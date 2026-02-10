package transactions_usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
)

const (
	checkUserActivity             = "CheckUser"
	checkCafeActivity             = "CheckCafe"
	checkFranchiseActivity        = "CheckFranchise"
	saveTransactionHeaderActivity = "SaveTransactionHeader"
	publishTransactionActivity    = "PublishTransaction"
	checkStatusActivity           = "CheckStatus"
)

type (
	TransactionWorkflow interface {
		CheckoutTransactionV3(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	}

	transactionWorkflow struct {
		workflow                      temporal_client.Workflow
		transactionHeaderRepository   repository_interfaces.TransactionHeaderRepository
		transactionDetailRepository   repository_interfaces.TransactionDetailRepository
		addressRepository             repository_interfaces.AddressRepository
		cartRepository                repository_interfaces.CartRepository
		userRepository                repository_interfaces.UserRepository
		cafeRepository                repository_interfaces.CafeRepository
		cafeFranchiseRepository       repository_interfaces.CafeFranchiseRepository
		productRepository             repository_interfaces.ProductRepository
		cafeProductRepository         repository_interfaces.CafeProductRepository
		transaction                   repository_interfaces.Transaction
		transactionTopic              topics_interfaces.TransactionTopic
		midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
		midtrans                      midtrans_client.Midtrans
		transactionCache              cache_interfaces.TransactionCache
		productCache                  cache_interfaces.ProductCache
		checkoutCache                 cache_interfaces.CheckoutCache
	}
)

func (t *transactionWorkflow) registerWorkflowAndActivities() {
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkUserActivity,
		Fn:   t.checkUser,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkCafeActivity,
		Fn:   t.checkCafe,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkFranchiseActivity,
		Fn:   t.checkFranchise,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: saveTransactionHeaderActivity,
		Fn:   t.saveTransactionHeader,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: publishTransactionActivity,
		Fn:   t.publishTransaction,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkStatusActivity,
		Fn:   t.checkStatus,
	})
	t.workflow.RegisterWorkflow(temporal_client.WorkflowDefinition{
		Name: "CreateTransaction",
		Fn:   t.transactionWorkflow,
	})
}

func (t *transactionWorkflow) CheckoutTransactionV3(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userID, ok := ctx.Value(enums.UserID).(uint64)
	if !ok {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user id from context`, fmt.Errorf("user id not found in context"))
	}

	workflowRun, err := t.workflow.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: "CreateTransaction",
	}, t.transactionWorkflow, userID, request)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start workflow`, err)
	}
	log.Println("Workflow started", workflowRun.GetID())

	// Use background context to avoid HTTP request timeout cancelling the wait
	err = t.workflow.GetWorkflowResult(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), &result)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get workflow result`, err)
	}

	if result == nil {
		return nil, apperror.NewCustomError(apperror.ErrNotFound, `result not found`, nil)
	}

	return result, nil
}

func NewTransactionWorkflow(
	transactionHeaderRepository repository_interfaces.TransactionHeaderRepository,
	transactionDetailRepository repository_interfaces.TransactionDetailRepository,
	addressRepository repository_interfaces.AddressRepository,
	cartRepository repository_interfaces.CartRepository,
	userRepository repository_interfaces.UserRepository,
	cafeRepository repository_interfaces.CafeRepository,
	cafeFranchiseRepository repository_interfaces.CafeFranchiseRepository,
	productRepository repository_interfaces.ProductRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	transaction repository_interfaces.Transaction,
	transactionTopic topics_interfaces.TransactionTopic,
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository,
	midtrans midtrans_client.Midtrans,
	transactionCache cache_interfaces.TransactionCache,
	productCache cache_interfaces.ProductCache,
	checkoutCache cache_interfaces.CheckoutCache,
	workflow temporal_client.Workflow,
) TransactionWorkflow {
	tw := &transactionWorkflow{
		transactionHeaderRepository:   transactionHeaderRepository,
		transactionDetailRepository:   transactionDetailRepository,
		addressRepository:             addressRepository,
		cartRepository:                cartRepository,
		userRepository:                userRepository,
		cafeRepository:                cafeRepository,
		cafeFranchiseRepository:       cafeFranchiseRepository,
		productRepository:             productRepository,
		cafeProductRepository:         cafeProductRepository,
		transaction:                   transaction,
		transactionTopic:              transactionTopic,
		midtransTransactionRepository: midtransTransactionRepository,
		midtrans:                      midtrans,
		transactionCache:              transactionCache,
		productCache:                  productCache,
		checkoutCache:                 checkoutCache,
		workflow:                      workflow,
	}
	tw.registerWorkflowAndActivities()
	return tw
}

package auto_transfer_usecases

import (
	"context"
	"fmt"

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
	autoTransferCheckUserActivity             = "AutoTransferCheckUser"
	autoTransferCheckCafeActivity             = "AutoTransferCheckCafe"
	autoTransferCheckFranchiseActivity        = "AutoTransferCheckFranchise"
	autoTransferSaveTransactionHeaderActivity = "AutoTransferSaveTransactionHeader"
	autoTransferPublishTransactionActivity    = "AutoTransferPublishTransaction"
	autoTransferStockDeductionActivity        = "AutoTransferStockDeduction"
	autoTransferSaveTransactionDetailActivity = "AutoTransferSaveTransactionDetail"
)

type (
	AutoTransferWorkflow interface {
		AutoTransfer(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	}

	autoTransferWorkflow struct {
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

func (a *autoTransferWorkflow) registerWorkflowAndActivities() {
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferCheckUserActivity,
		Fn:   a.checkUser,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferCheckCafeActivity,
		Fn:   a.checkCafe,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferCheckFranchiseActivity,
		Fn:   a.checkFranchise,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferSaveTransactionHeaderActivity,
		Fn:   a.saveTransactionHeader,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferPublishTransactionActivity,
		Fn:   a.publishTransaction,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferStockDeductionActivity,
		Fn:   a.stockDeduction,
	})
	a.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: autoTransferSaveTransactionDetailActivity,
		Fn:   a.saveTransactionDetail,
	})

	a.workflow.RegisterWorkflow(temporal_client.WorkflowDefinition{
		Name: "AutoTransfer",
		Fn:   a.autoTransferWorkflow,
	})
}

// AutoTransfer implements [AutoTransferWorkflow].
func (a *autoTransferWorkflow) AutoTransfer(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userID, ok := ctx.Value(enums.UserID).(uint64)
	if !ok {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user id from context`, fmt.Errorf("user id not found in context"))
	}

	workflowRun, err := a.workflow.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: "AutoTransfer",
	}, a.autoTransferWorkflow, userID, request)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start workflow`, err)
	}

	err = a.workflow.GetWorkflowResult(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), &result)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get workflow result`, err)
	}

	if result == nil {
		return nil, apperror.NewCustomError(apperror.ErrNotFound, `result not found`, nil)
	}

	return result, nil
}

func NewAutoTransferWorkflow(
	workflow temporal_client.Workflow,
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
) AutoTransferWorkflow {
	atw := &autoTransferWorkflow{
		workflow:                      workflow,
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
	}
	atw.registerWorkflowAndActivities()
	return atw
}

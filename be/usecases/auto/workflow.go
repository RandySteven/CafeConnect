package auto_transfer_usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	autoTransferCheckUserActivity             = "AutoTransferCheckUser"
	autoTransferCheckCafeActivity             = "AutoTransferCheckCafe"
	autoTransferCheckFranchiseActivity        = "AutoTransferCheckFranchise"
	autoTransferSaveTransactionHeaderActivity = "AutoTransferSaveTransactionHeader"
	autoTransferPublishTransactionActivity    = "AutoTransferPublishTransaction"
	autoTransferStockDeductionActivity        = "AutoTransferStockDeduction"
	autoTransferSaveTransactionDetailActivity = "AutoTransferSaveTransactionDetail"
	autoTransferRestoreStockActivity          = "AutoTransferRestoreStock"

	sgNoNeed   = ""
	sgMidtrans = "MidtransTransaction"
)

type (
	AutoTransferWorkflow interface {
		AutoTransfer(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	}

	autoTransferWorkflow struct {
		temporal                      temporal_client.Temporal
		workflow                      temporal_client.WorkflowExecution
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
	ao := &workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        5,
			InitialInterval:        1 * time.Minute,
			BackoffCoefficient:     2,
			MaximumInterval:        5 * time.Minute,
			NonRetryableErrorTypes: nonRetryableErrorTypes,
		},
	}
	a.workflow.AddTransitionActivityWithOptions(autoTransferCheckUserActivity, sgNoNeed, a.checkUser, nil)
	a.workflow.AddTransitionActivityWithOptions(autoTransferCheckCafeActivity, sgNoNeed, a.checkCafe, nil)
	a.workflow.AddTransitionActivityWithOptions(autoTransferCheckFranchiseActivity, sgNoNeed, a.checkFranchise, nil)
	a.workflow.AddTransitionActivityWithOptions(autoTransferSaveTransactionHeaderActivity, sgNoNeed, a.saveTransactionHeader, nil)
	a.workflow.AddTransitionActivityWithOptions(autoTransferStockDeductionActivity, sgNoNeed, a.stockDeduction, ao)
	a.workflow.AddTransitionActivityWithOptions(autoTransferPublishTransactionActivity, sgNoNeed, a.publishTransaction, nil)
	a.workflow.AddTransitionActivityWithOptions(autoTransferSaveTransactionDetailActivity, sgNoNeed, a.saveTransactionDetail, nil)

	a.workflow.AddBranchActivity(autoTransferRestoreStockActivity, a.restoreStock)

	a.workflow.RegisterWorkflow("AutoTransfer", a.autoTransferWorkflow)
}

// AutoTransfer implements [AutoTransferWorkflow].
func (a *autoTransferWorkflow) AutoTransfer(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userID, ok := ctx.Value(enums.UserID).(uint64)
	if !ok {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user id from context`, fmt.Errorf("user id not found in context"))
	}

	workflowRun, err := a.temporal.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: fmt.Sprintf("AutoTransfer-%s-%d", request.IdempotencyKey, time.Now().UnixNano()),
	}, a.autoTransferWorkflow, userID, request)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start workflow`, err)
	}

	err = a.temporal.GetWorkflowResult(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), &result)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get workflow result`, err)
	}

	if result == nil {
		return nil, apperror.NewCustomError(apperror.ErrNotFound, `result not found`, nil)
	}

	return result, nil
}

func NewAutoTransferWorkflow(
	temporal temporal_client.Temporal,
	workflow temporal_client.WorkflowExecution,
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
		temporal:                      temporal,
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

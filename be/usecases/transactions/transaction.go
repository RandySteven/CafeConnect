package transactions_usecases

import (
	"context"
	"errors"
	"log"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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
		CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, err error)
	}

	transactionWorkflow struct {
		worker                        worker.Worker
		client                        client.Client
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

func (t *transactionWorkflow) RegisterWorkflow(worker worker.Worker) {
	t.worker.RegisterActivityWithOptions(t.checkUser, activity.RegisterOptions{
		Name: checkUserActivity,
	})
	t.worker.RegisterActivityWithOptions(t.checkCafe, activity.RegisterOptions{
		Name: checkCafeActivity,
	})
	t.worker.RegisterActivityWithOptions(t.checkFranchise, activity.RegisterOptions{
		Name: checkFranchiseActivity,
	})
	t.worker.RegisterActivityWithOptions(t.saveTransactionHeader, activity.RegisterOptions{
		Name: saveTransactionHeaderActivity,
	})
	t.worker.RegisterActivityWithOptions(t.publishTransaction, activity.RegisterOptions{
		Name: publishTransactionActivity,
	})
	t.worker.RegisterActivityWithOptions(t.checkStatus, activity.RegisterOptions{
		Name: checkStatusActivity,
	})
	t.worker.RegisterWorkflow(t.transactionWorkflow)
}

func (t *transactionWorkflow) CheckoutTransactionV3(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, err error) {
	workflowRun, err := t.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: "CreateTransaction",
	}, "CreateTransaction", t.transactionWorkflow, request)
	if err != nil {
		return nil, err
	}
	log.Println("Workflow started", workflowRun.GetID())

	result = ctx.Value("result").(*responses.TransactionReceiptResponse)
	if result == nil {
		return nil, errors.New("result not found")
	}

	return result, nil
}

func (t *transactionWorkflow) transactionWorkflow(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, err error) {
	workflowCtx := workflow.WithChildOptions(ctx.(workflow.Context), workflow.ChildWorkflowOptions{
		WorkflowID: "CheckUser",
	})

	errCh := make(chan error)

	defer close(errCh)

	workflow.Go(workflowCtx, func(ctx workflow.Context) {
		checkUserActivity := workflow.ExecuteActivity(workflowCtx, t.checkUser, request)
		if err := checkUserActivity.Get(workflowCtx, nil); err != nil {
			errCh <- err
			return
		}

		checkCafeActivity := workflow.ExecuteActivity(workflowCtx, t.checkCafe, request)
		if err := checkCafeActivity.Get(workflowCtx, nil); err != nil {
			errCh <- err
			return
		}
	})

	select {
	case err = <-errCh:
		return nil, err
	default:
		checkFranchiseActivity := workflow.ExecuteActivity(workflowCtx, t.checkFranchise, request)
		if err := checkFranchiseActivity.Get(workflowCtx, nil); err != nil {
			return nil, err
		}

		saveTransactionHeaderActivity := workflow.ExecuteActivity(workflowCtx, t.saveTransactionHeader, request)
		if err := saveTransactionHeaderActivity.Get(workflowCtx, nil); err != nil {
			return nil, err
		}

		publishTransactionActivity := workflow.ExecuteActivity(workflowCtx, t.publishTransaction, request)
		if err := publishTransactionActivity.Get(workflowCtx, nil); err != nil {
			return nil, err
		}

		transactionHeader := ctx.Value("transactionHeader").(*models.TransactionHeader)

		result = &responses.TransactionReceiptResponse{
			ID:              transactionHeader.ID,
			TransactionCode: transactionHeader.TransactionCode,
			Status:          transactionHeader.Status,
			TransactionAt:   transactionHeader.TransactionAt.Local(),
		}
		ctx = context.WithValue(ctx, "result", result)

		return result, nil
	}
}

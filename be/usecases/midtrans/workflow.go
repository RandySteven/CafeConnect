package midtrans_usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
	"github.com/midtrans/midtrans-go"
	"go.temporal.io/sdk/workflow"
)

const (
	checkTransactionHeaderActivity    = "CheckTransactionHeader"
	checkoutListActivity              = "CheckoutList"
	createMidtransTransactionActivity = "CreateMidtransTransaction"

	sgNoNeed = ""
)

type (
	midtransCheckOut struct {
		Items       []midtrans.ItemDetails
		TotalAmount int64
	}

	MidtransWorkflow interface {
		CreateMidtransTransaction(ctx context.Context, message *messages.TransactionMidtransMessage) (result *midtrans_client.MidtransResponse, customErr *apperror.CustomError)
	}

	midtransWorkflow struct {
		temporal                      temporal_client.Temporal
		workflow                      temporal_client.WorkflowExecution
		transactionHeaderRepository   repository_interfaces.TransactionHeaderRepository
		midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
		transactionDetailRepository   repository_interfaces.TransactionDetailRepository
		cafeProductRepository         repository_interfaces.CafeProductRepository
		productRepository             repository_interfaces.ProductRepository
		midtrans                      midtrans_client.Midtrans
	}
)

func (m *midtransWorkflow) registerWorkflowAndActivities() {
	opts := &workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}

	m.workflow.AddTransitionActivityWithOptions(checkTransactionHeaderActivity, sgNoNeed, m.checkTransactionHeader, opts)
	m.workflow.AddTransitionActivityWithOptions(checkoutListActivity, sgNoNeed, m.checkoutList, opts)
	m.workflow.AddTransitionActivityWithOptions(createMidtransTransactionActivity, sgNoNeed, m.createMidtransTransaction, opts)
	m.workflow.RegisterWorkflow("MidtransTransaction", m.midtransTransaction)
}

// CreateMidtransTransaction starts the Midtrans workflow, signals it with
// the transaction data, and waits for the result.
func (m *midtransWorkflow) CreateMidtransTransaction(ctx context.Context, message *messages.TransactionMidtransMessage) (result *midtrans_client.MidtransResponse, customErr *apperror.CustomError) {
	workflowRun, err := m.workflow.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: fmt.Sprintf("MidtransTransaction-%s", message.TransactionCode),
	}, m.midtransTransaction)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start workflow`, err)
	}

	// Signal the workflow with the transaction data
	err = m.workflow.SignalWorkflow(ctx, workflowRun.GetID(), workflowRun.GetRunID(), "MidtransTransaction", message)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to signal midtrans workflow`, err)
	}

	// Wait for the workflow to complete and return the Midtrans response
	err = m.workflow.GetWorkflowResult(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), &result)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get workflow result`, err)
	}

	return result, nil
}

func NewMidtransWorkflow(temporal temporal_client.Temporal,
	workflow temporal_client.WorkflowExecution,
	transactionHeaderRepository repository_interfaces.TransactionHeaderRepository,
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository,
	transactionDetailRepository repository_interfaces.TransactionDetailRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	productRepository repository_interfaces.ProductRepository,
	midtrans midtrans_client.Midtrans) MidtransWorkflow {
	mw := &midtransWorkflow{
		temporal:                      temporal,
		workflow:                      workflow,
		transactionHeaderRepository:   transactionHeaderRepository,
		midtransTransactionRepository: midtransTransactionRepository,
		transactionDetailRepository:   transactionDetailRepository,
		cafeProductRepository:         cafeProductRepository,
		productRepository:             productRepository,
		midtrans:                      midtrans,
	}
	mw.registerWorkflowAndActivities()
	return mw
}

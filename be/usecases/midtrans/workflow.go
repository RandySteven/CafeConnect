package midtrans_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
	"github.com/midtrans/midtrans-go"
)

const (
	checkTransactionHeaderActivity    = "CheckTransactionHeader"
	checkoutListActivity              = "CheckoutList"
	createMidtransTransactionActivity = "CreateMidtransTransaction"
)

type (
	midtransCheckOut struct {
		Items       []midtrans.ItemDetails
		TotalAmount int64
	}

	MidtransWorkflow interface {
		CreateMidtransTransaction(ctx context.Context, message *messages.TransactionMidtransMessage) (err error)
	}

	midtransWorkflow struct {
		workflow                      temporal_client.Workflow
		transactionHeaderRepository   repository_interfaces.TransactionHeaderRepository
		midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
		transactionDetailRepository   repository_interfaces.TransactionDetailRepository
		cafeProductRepository         repository_interfaces.CafeProductRepository
		productRepository             repository_interfaces.ProductRepository
		midtrans                      midtrans_client.Midtrans
	}
)

func (m *midtransWorkflow) registerWorkflowAndActivities() {
	m.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkTransactionHeaderActivity,
		Fn:   m.checkTransactionHeader,
	})

	m.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkoutListActivity,
		Fn:   m.checkoutList,
	})

	m.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: createMidtransTransactionActivity,
		Fn:   m.createMidtransTransaction,
	})

	m.workflow.RegisterWorkflow(temporal_client.WorkflowDefinition{
		Name: "MidtransTransaction",
		Fn:   m.midtransTransaction,
	})
}

func (m *midtransWorkflow) CreateMidtransTransaction(ctx context.Context, message *messages.TransactionMidtransMessage) (err error) {
	return
}

func NewMidtransWorkflow(workflow temporal_client.Workflow) MidtransWorkflow {
	mw := &midtransWorkflow{
		workflow: workflow,
	}
	mw.registerWorkflowAndActivities()
	return mw
}

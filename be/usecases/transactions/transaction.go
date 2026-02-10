package transactions_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"go.temporal.io/sdk/workflow"
)

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

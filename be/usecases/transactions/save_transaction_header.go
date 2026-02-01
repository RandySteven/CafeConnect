package transactions_usecases

import (
	"context"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *transactionWorkflow) saveTransactionHeader(ctx context.Context, request *requests.CreateTransactionRequest) (err error) {
	user := ctx.Value("user").(*models.User)
	transactionHeader := &models.TransactionHeader{
		UserID:          user.ID,
		TransactionCode: utils.GenerateCode(24),
		CafeID:          request.CafeID,
		Status:          enums.TransactionPENDING.String(),
		TransactionAt:   time.Now(),
	}
	transactionHeader, err = t.transactionHeaderRepository.Save(ctx, transactionHeader)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, "transactionHeader", transactionHeader)

	return nil
}

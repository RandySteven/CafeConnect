package transactions_usecases

import (
	"context"
	"errors"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
)

func (t *transactionWorkflow) checkUser(ctx context.Context, request *requests.CreateTransactionRequest) (err error) {
	userId := ctx.Value(enums.UserID).(uint64)
	user, err := t.userRepository.FindByID(ctx, userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	ctx = context.WithValue(ctx, "user", user)
	return nil
}

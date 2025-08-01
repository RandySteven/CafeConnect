package consumer_interfaces

import "context"

type TransactionConsumer interface {
	MidtransTransactionRecord(ctx context.Context) error
	MidtransPaymentConfirmation(ctx context.Context) error
}

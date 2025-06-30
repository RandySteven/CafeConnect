package job_interfaces

import "context"

type TransactionJob interface {
	CheckMidtransStatus(ctx context.Context) (err error)
	HardFailedPendingTrx(ctx context.Context) (err error)
}

package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type TransactionUsecase interface {
	CreateTransactionV1(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	CreateTransactionV2(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	GetUserTransactions(ctx context.Context) (result []*responses.TransactionListResponse, customErr *apperror.CustomError)
	GetTransactionByCode(ctx context.Context, transactionCode string) (result *responses.TransactionDetailResponse, customErr *apperror.CustomError)
}

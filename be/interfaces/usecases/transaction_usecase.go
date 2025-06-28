package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type TransactionUsecase interface {
	CheckoutTransactionV1(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	CheckoutTransactionV2(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	CreateTransactionV1(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	CreateTransactionV2(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	GetUserTransactions(ctx context.Context) (result []*responses.TransactionListResponse, customErr *apperror.CustomError)
	GetTransactionByCode(ctx context.Context, transactionCode string) (result *responses.TransactionDetailResponse, customErr *apperror.CustomError)
	CheckReceipt(ctx context.Context, transactionCode string) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
	PaymentConfirmation(ctx context.Context, request *requests.PaymentConfirmationRequest) (result []*responses.PaymentConfirmationResponse, message string, customErr *apperror.CustomError)
}

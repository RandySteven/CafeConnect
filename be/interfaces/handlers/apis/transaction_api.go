package api_interfaces

import "net/http"

type TransactionApi interface {
	CheckoutTransactionV1(w http.ResponseWriter, r *http.Request)
	CheckoutTransactionV2(w http.ResponseWriter, r *http.Request)
	CheckoutTransactionV3(w http.ResponseWriter, r *http.Request)
	GetUserTransactions(w http.ResponseWriter, r *http.Request)
	GetTransactionByTransactionCode(w http.ResponseWriter, r *http.Request)
	CheckReceipt(w http.ResponseWriter, r *http.Request)
	PaymentConfirmation(w http.ResponseWriter, r *http.Request)
}

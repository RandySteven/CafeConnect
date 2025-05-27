package api_interfaces

import "net/http"

type TransactionApi interface {
	CheckoutTransaction(w http.ResponseWriter, r *http.Request)
	GetUserTransactions(w http.ResponseWriter, r *http.Request)
	GetTransactionByTransactionCode(w http.ResponseWriter, r *http.Request)
}

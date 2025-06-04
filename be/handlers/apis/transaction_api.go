package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type TransactionApi struct {
	usecase usecase_interfaces.TransactionUsecase
}

func (t *TransactionApi) CheckoutTransactionV1(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `result`
	)
	result, customErr := t.usecase.CreateTransactionV1(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success create transaction`, &dataKey, result, nil)
}

func (t *TransactionApi) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `transactions`
	)
	result, customErr := t.usecase.GetUserTransactions(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success get transactions`, &dataKey, result, nil)
}

func (t *TransactionApi) GetTransactionByTransactionCode(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `transaction`
		vars    = mux.Vars(r)
	)
	transactionCode := vars[`transactionCode`]
	result, customErr := t.usecase.GetTransactionByCode(ctx, transactionCode)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success get transaction`, &dataKey, result, nil)
}

func (t *TransactionApi) CheckoutTransactionV2(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.CreateTransactionRequest{}
		dataKey = `result`
	)

	if err := utils.BindJSON(r, &request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}

	result, customErr := t.usecase.CheckoutTransactionV2(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success create transaction`, &dataKey, result, nil)
}

var _ api_interfaces.TransactionApi = &TransactionApi{}

func newTransactionApi(usecase usecase_interfaces.TransactionUsecase) *TransactionApi {
	return &TransactionApi{
		usecase: usecase,
	}
}

package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"net/http"
)

type AddressApi struct {
	usecase usecase_interfaces.AddressUsecase
}

func (a *AddressApi) AddUserAddress(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.AddAddressRequest{}
		dataKey = `result`
	)

	if err := utils.BindRequest(r, &request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `bad request`, nil, nil, err)
		return
	}

	result, customErr := a.usecase.AddUserAddress(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusCreated, `success add address user`, &dataKey, result, nil)
}

func (a *AddressApi) GetUserAddresses(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `addresses`
	)

	result, customErr := a.usecase.GetUserAddresses(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get address user`, &dataKey, result, nil)
}

var _ api_interfaces.AddressApi = &AddressApi{}

func newAddressApi(usecase usecase_interfaces.AddressUsecase) *AddressApi {
	return &AddressApi{
		usecase: usecase,
	}
}

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

type CartApi struct {
	usecase usecase_interfaces.CartUsecase
}

func (c *CartApi) AddCart(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.AddToCartRequest{}
		dataKey = `result`
	)

	if err := utils.BindRequest(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}
	result, customErr := c.usecase.AddToCart(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusCreated, `success add cart user`, &dataKey, result, nil)
}

func (c *CartApi) GetCart(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `cart`
	)
	result, customErr := c.usecase.GetUserCart(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get user cart`, &dataKey, result, nil)
}

var _ api_interfaces.CartApi = &CartApi{}

func newCartApi(usecase usecase_interfaces.CartUsecase) *CartApi {
	return &CartApi{
		usecase: usecase,
	}
}

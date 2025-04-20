package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"net/http"
)

type CafeApi struct {
	usecase usecase_interfaces.CafeUsecase
}

func (c *CafeApi) GetListCafeFranchise(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `franchises`
	)
	result, customErr := c.usecase.GetListCafeFranchises(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success login user`, &dataKey, result, nil)
}

func (c *CafeApi) RegisterCafe(w http.ResponseWriter, r *http.Request) {
}

func (c *CafeApi) GetCafeDetail(w http.ResponseWriter, r *http.Request) {
}

func (c *CafeApi) GetCafeProducts(w http.ResponseWriter, r *http.Request) {
}

func (c *CafeApi) GetListOfCafeBasedOnRadius(w http.ResponseWriter, r *http.Request) {
}

var _ api_interfaces.CafeApi = &CafeApi{}

func newCafeApi(usecase usecase_interfaces.CafeUsecase) *CafeApi {
	return &CafeApi{
		usecase: usecase,
	}
}

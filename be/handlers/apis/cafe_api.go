package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"net/http"
)

type CafeApi struct {
	usecase usecase_interfaces.CafeUsecase
}

func (c *CafeApi) RegisterCafe(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *CafeApi) GetCafeDetail(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *CafeApi) GetCafeProducts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *CafeApi) GetListOfCafeBasedOnRadius(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

var _ api_interfaces.CafeApi = &CafeApi{}

func newCafeApi(usecase usecase_interfaces.CafeUsecase) *CafeApi {
	return &CafeApi{
		usecase: usecase,
	}
}

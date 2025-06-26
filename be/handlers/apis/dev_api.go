package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	"github.com/RandySteven/CafeConnect/be/utils"
	"net/http"
)

type DevApi struct {
}

func (d *DevApi) CheckHealth(w http.ResponseWriter, r *http.Request) {
	utils.ResponseHandler(w, http.StatusOK, `success check health`, nil, nil, nil)
}

func (d *DevApi) SendEmail(w http.ResponseWriter, r *http.Request) {}

func (d *DevApi) RouterList(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

var _ api_interfaces.DevApi = &DevApi{}

func newDevApi() *DevApi {
	return &DevApi{}
}

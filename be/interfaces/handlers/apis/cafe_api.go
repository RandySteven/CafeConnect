package api_interfaces

import "net/http"

type CafeApi interface {
	RegisterCafeAndFranchise(w http.ResponseWriter, r *http.Request)
	GetCafeDetail(w http.ResponseWriter, r *http.Request)
	GetCafeProducts(w http.ResponseWriter, r *http.Request)
	GetListOfCafeBasedOnRadius(w http.ResponseWriter, r *http.Request)
	GetListCafeFranchise(w http.ResponseWriter, r *http.Request)
}

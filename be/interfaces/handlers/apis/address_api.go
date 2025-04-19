package api_interfaces

import "net/http"

type AddressApi interface {
	AddUserAddress(w http.ResponseWriter, r *http.Request)
}

package api_interfaces

import "net/http"

type CartApi interface {
	AddCart(w http.ResponseWriter, r *http.Request)
	GetCart(w http.ResponseWriter, r *http.Request)
}

package api_interfaces

import "net/http"

type DevApi interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
	RouterList(w http.ResponseWriter, r *http.Request)
}

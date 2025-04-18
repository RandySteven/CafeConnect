package api_interfaces

import "net/http"

type OnboardingApi interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	GoogleLogin(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	GetOnboardUser(w http.ResponseWriter, r *http.Request)
}

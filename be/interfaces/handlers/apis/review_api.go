package api_interfaces

import "net/http"

type ReviewApi interface {
	GetCafeReviews(w http.ResponseWriter, r *http.Request)
	AddCafeReview(w http.ResponseWriter, r *http.Request)
}

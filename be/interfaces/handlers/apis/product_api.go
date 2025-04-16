package api_interfaces

import "net/http"

type ProductApi interface {
	AddProduct(w http.ResponseWriter, r *http.Request)
	GetListOfProducts(w http.ResponseWriter, r *http.Request)
	GetProductDetail(w http.ResponseWriter, r *http.Request)
}

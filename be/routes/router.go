package routes

import (
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/RandySteven/CafeConnect/be/handlers/apis"
	"github.com/RandySteven/CafeConnect/be/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request)

	Router struct {
		path        string
		handler     HandlerFunc
		method      string
		middlewares []enums.Middleware
	}

	RouterPrefix map[enums.RouterPrefix][]*Router
)

func NewEndpointRouters(api *apis.APIs) RouterPrefix {
	endpoint := make(RouterPrefix)

	endpoint[enums.DevPrefix] = []*Router{
		Get(`/check-health`, api.DevApi.CheckHealth),
	}

	endpoint[enums.OnboardingPrefix] = []*Router{
		Post(`/register`, api.OnboardingApi.RegisterUser),
		Post(`/login`, api.OnboardingApi.LoginUser),
		Get(``, api.OnboardingApi.GetOnboardUser, enums.AuthenticationMiddleware),
	}

	endpoint[enums.RolePrefix] = []*Router{
		Get(``, api.RoleApi.GetRoleList),
	}

	endpoint[enums.AddressPrefix] = []*Router{
		Post(``, api.AddressApi.AddUserAddress, enums.AuthenticationMiddleware),
		Get(``, api.AddressApi.GetUserAddresses, enums.AuthenticationMiddleware),
	}

	endpoint[enums.CafePrefix] = []*Router{
		Get(`/franchises`, api.CafeApi.GetListCafeFranchise),
		Post(`/franchises`, api.CafeApi.RegisterCafeAndFranchise),
		Get(`/{id}`, api.CafeApi.GetCafeDetail),
		Post(``, api.CafeApi.GetListOfCafeBasedOnRadius),
		Post(`/add-outlet`, api.CafeApi.AddCafeOutlet),
	}

	endpoint[enums.ReviewPrefix] = []*Router{
		Post(`/cafe-review`, api.ReviewApi.GetCafeReviews),
		Post(``, api.ReviewApi.AddCafeReview, enums.AuthenticationMiddleware),
	}

	endpoint[enums.ProductPrefix] = []*Router{
		Post(``, api.ProductApi.AddProduct),
		Get(`/{id}`, api.ProductApi.GetProductDetail),
		Post(`/cafe-list`, api.ProductApi.GetListOfProducts),
	}

	endpoint[enums.CartPrefix] = []*Router{
		Post(``, api.CartApi.AddCart, enums.AuthenticationMiddleware),
		Get(``, api.CartApi.GetCart, enums.AuthenticationMiddleware),
	}

	endpoint[enums.TransactionPrefix] = []*Router{
		Get(`/v1/check-out`, api.TransactionApi.CheckoutTransactionV1, enums.AuthenticationMiddleware),
		Post(`/v2/check-out`, api.TransactionApi.CheckoutTransactionV2, enums.AuthenticationMiddleware),
		Get(`/{transactionCode}`, api.TransactionApi.GetTransactionByTransactionCode, enums.AuthenticationMiddleware),
		Get(``, api.TransactionApi.GetUserTransactions, enums.AuthenticationMiddleware),
	}

	return endpoint
}

func InitRouter(routers RouterPrefix, r *mux.Router) {
	middleware := middlewares.NewMiddlewares()
	clientMiddleware := middlewares.RegisterClientMiddleware(middleware)
	serverMiddleware := middlewares.RegisterServerMiddleware(middleware)

	r.Use(
		serverMiddleware.LoggingMiddleware,
		serverMiddleware.CorsMiddleware,
		serverMiddleware.TimeoutMiddleware,
		serverMiddleware.CheckHealthMiddleware,
		clientMiddleware.AuthenticationMiddleware,
		clientMiddleware.RateLimiterMiddleware,
	)

	devRouter := r.PathPrefix(enums.DevPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.DevPrefix] {
		devRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.DevPrefix.ToString())
	}

	onboardingRouter := r.PathPrefix(enums.OnboardingPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.OnboardingPrefix] {
		middleware.RegisterMiddleware(enums.OnboardingPrefix, router.method, router.path, router.middlewares)
		onboardingRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.OnboardingPrefix.ToString())
	}

	cafeRouter := r.PathPrefix(enums.CafePrefix.ToString()).Subrouter()
	for _, router := range routers[enums.CafePrefix] {
		middleware.RegisterMiddleware(enums.CafePrefix, router.method, router.path, router.middlewares)
		cafeRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.CafePrefix.ToString())
	}

	reviewRouter := r.PathPrefix(enums.ReviewPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.ReviewPrefix] {
		middleware.RegisterMiddleware(enums.ReviewPrefix, router.method, router.path, router.middlewares)
		reviewRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.ReviewPrefix.ToString())
	}

	productRouter := r.PathPrefix(enums.ProductPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.ProductPrefix] {
		middleware.RegisterMiddleware(enums.ProductPrefix, router.method, router.path, router.middlewares)
		productRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.ProductPrefix.ToString())
	}

	cartRouter := r.PathPrefix(enums.CartPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.CartPrefix] {
		middleware.RegisterMiddleware(enums.CartPrefix, router.method, router.path, router.middlewares)
		cartRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.CartPrefix.ToString())
	}

	transactionRouter := r.PathPrefix(enums.TransactionPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.TransactionPrefix] {
		middleware.RegisterMiddleware(enums.TransactionPrefix, router.method, router.path, router.middlewares)
		transactionRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.TransactionPrefix.ToString())
	}

	roleRouter := r.PathPrefix(enums.RolePrefix.ToString()).Subrouter()
	for _, router := range routers[enums.RolePrefix] {
		middleware.RegisterMiddleware(enums.RolePrefix, router.method, router.path, router.middlewares)
		roleRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.RolePrefix.ToString())
	}

	addressRouter := r.PathPrefix(enums.AddressPrefix.ToString()).Subrouter()
	for _, router := range routers[enums.AddressPrefix] {
		middleware.RegisterMiddleware(enums.AddressPrefix, router.method, router.path, router.middlewares)
		addressRouter.HandleFunc(router.path, router.handler).Methods(router.method)
		router.RouterLog(enums.AddressPrefix.ToString())
	}

}

func (router *Router) RouterLog(prefix string) {
	log.Printf("%12s | %4s/ \n", router.method, prefix+router.path)
}

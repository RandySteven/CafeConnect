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

	endpoint[enums.CafePrefix] = []*Router{
		Get(`/franchises`, api.CafeApi.GetListCafeFranchise),
		Post(`/franchises`, api.CafeApi.RegisterCafeAndFranchise),
		Get(`/{id}`, api.CafeApi.GetCafeDetail),
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
}

func (router *Router) RouterLog(prefix string) {
	log.Printf("%12s | %4s/ \n", router.method, prefix+router.path)
}

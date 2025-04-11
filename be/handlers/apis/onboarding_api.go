package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"net/http"
)

type OnboardingApi struct {
	usecase usecase_interfaces.OnboardingUsecase
}

func (o *OnboardingApi) RegisterUser(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GoogleLogin(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GoogleCallback(w http.ResponseWriter, r *http.Request) {
}

var _ api_interfaces.OnboardingApi = &OnboardingApi{}

func NewOnboardingApi(usecase usecase_interfaces.OnboardingUsecase) *OnboardingApi {
	return &OnboardingApi{
		usecase: usecase,
	}
}

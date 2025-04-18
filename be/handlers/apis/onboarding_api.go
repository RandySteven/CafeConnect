package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"net/http"
)

type OnboardingApi struct {
	usecase usecase_interfaces.OnboardingUsecase
}

func (o *OnboardingApi) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.RegisterUserRequest{}
		dataKey = `data`
	)

	imageFile, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	request.ProfilePicture = imageFile

	if err := utils.BindMultipartForm(r, &request); err != nil {
		return
	}
	ctx2 := context.WithValue(ctx, enums.FileHeader, fileHeader)
	ctx = ctx2
	result, customErr := o.usecase.RegisterUser(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusCreated, `success register user`, &dataKey, result, nil)
}

func (o *OnboardingApi) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GoogleLogin(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GoogleCallback(w http.ResponseWriter, r *http.Request) {
}

var _ api_interfaces.OnboardingApi = &OnboardingApi{}

func newOnboardingApi(usecase usecase_interfaces.OnboardingUsecase) *OnboardingApi {
	return &OnboardingApi{
		usecase: usecase,
	}
}

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
	"strconv"
)

type OnboardingApi struct {
	usecase usecase_interfaces.OnboardingUsecase
}

func (o *OnboardingApi) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		rID          = uuid.NewString()
		ctx          = context.WithValue(r.Context(), enums.RequestID, rID)
		longitude, _ = strconv.ParseFloat(r.FormValue("longitude"), 64)
		latitude, _  = strconv.ParseFloat(r.FormValue("latitude"), 64)
		request      = &requests.RegisterUserRequest{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Username:     r.FormValue("username"),
			Email:        r.FormValue("email"),
			Password:     r.FormValue("password"),
			PhoneNumber:  r.FormValue("phone_number"),
			DoB:          r.FormValue("dob"),
			ReferralCode: r.FormValue("referral_code"),
			Address:      r.FormValue("address"),
			Longitude:    longitude,
			Latitude:     latitude,
		}
		dataKey = `data`
	)

	imageFile, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	request.ProfilePicture = imageFile

	if err := utils.BindRequest(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
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
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.LoginUserRequest{}
		dataKey = `token`
	)

	if err := utils.BindJSON(r, &request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}

	result, customErr := o.usecase.LoginUser(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success login user`, &dataKey, result, nil)
}

func (o *OnboardingApi) GoogleLogin(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GoogleCallback(w http.ResponseWriter, r *http.Request) {
}

func (o *OnboardingApi) GetOnboardUser(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `result`
	)
	result, customErr := o.usecase.GetOnboardUser(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get user`, &dataKey, result, nil)
}

var _ api_interfaces.OnboardingApi = &OnboardingApi{}

func newOnboardingApi(usecase usecase_interfaces.OnboardingUsecase) *OnboardingApi {
	return &OnboardingApi{
		usecase: usecase,
	}
}

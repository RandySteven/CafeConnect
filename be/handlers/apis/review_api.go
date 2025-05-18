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

type ReviewApi struct {
	usecase usecase_interfaces.ReviewUsecase
}

func (r2 *ReviewApi) GetCafeReviews(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.GetCafeReviewRequest{}
		dataKey = `review`
	)

	if err := utils.BindRequest(r, &request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}

	result, customErr := r2.usecase.GetCafeReviews(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success get review list`, &dataKey, result, nil)
}

func (r2 *ReviewApi) AddCafeReview(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.AddReviewRequest{}
		dataKey = `result`
	)

	if err := utils.BindRequest(r, &request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}

	result, customErr := r2.usecase.AddCafeReview(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}

	utils.ResponseHandler(w, http.StatusOK, `success add review list`, &dataKey, result, nil)
}

var _ api_interfaces.ReviewApi = &ReviewApi{}

func newReviewApi(usecase usecase_interfaces.ReviewUsecase) *ReviewApi {
	return &ReviewApi{
		usecase: usecase,
	}
}

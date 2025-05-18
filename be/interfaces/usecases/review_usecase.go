package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type ReviewUsecase interface {
	GetCafeReviews(ctx context.Context, request *requests.GetCafeReviewRequest) (result *responses.GetReviewsResponse, customErr *apperror.CustomError)
	AddCafeReview(ctx context.Context, request *requests.AddReviewRequest) (result *responses.AddReviewResponse, customErr *apperror.CustomError)
}

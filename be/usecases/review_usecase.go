package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/google/uuid"
	"sync"
	"time"
)

type reviewUsecase struct {
	reviewRepo repository_interfaces.ReviewRepository
	cafeRepo   repository_interfaces.CafeRepository
	userRepo   repository_interfaces.UserRepository
}

func (r *reviewUsecase) GetCafeReviews(ctx context.Context, request *requests.GetCafeReviewRequest) (result *responses.GetReviewsResponse, customErr *apperror.CustomError) {
	var (
		wg            sync.WaitGroup
		numbOfWorkers = 2
		customErrCh   = make(chan *apperror.CustomError)
	)
	result = &responses.GetReviewsResponse{
		CafeID: request.CafeID,
	}
	wg.Add(numbOfWorkers)

	go func() {
		defer wg.Done()
		reviews, err := r.reviewRepo.FindByCafeID(ctx, request.CafeID)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get reviews`, err)
			return
		}
		if len(reviews) == 0 {
			return
		}
		resultReviews := make([]*responses.ReviewsResponse, len(reviews))
		for index, review := range reviews {
			user, err := r.userRepo.FindByID(ctx, review.UserID)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user`, err)
				return
			}
			resultReviews[index] = &responses.ReviewsResponse{
				User: struct {
					ID             uint64 `json:"id"`
					Name           string `json:"name"`
					ProfilePicture string `json:"profile_picture"`
				}{ID: user.ID, Name: user.Name, ProfilePicture: user.ProfilePicture},
				Score:     review.Score,
				Comment:   review.Comment,
				CreatedAt: review.CreatedAt,
			}
		}
		result.Reviews = resultReviews
	}()

	go func() {
		defer wg.Done()
		avg, err := r.reviewRepo.AvgCafeRating(ctx, request.CafeID)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get avg`, err)
		}

		result.AvgScore = avg
	}()

	go func() {
		wg.Wait()
		close(customErrCh)
	}()

	select {
	case customErr = <-customErrCh:
		return nil, customErr
	default:
		return result, nil
	}
}

func (r *reviewUsecase) AddCafeReview(ctx context.Context, request *requests.AddReviewRequest) (result *responses.AddReviewResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)
	review := &models.Review{
		UserID:  userId,
		CafeID:  request.CafeID,
		Score:   request.Score,
		Comment: request.Comment,
	}
	review, err := r.reviewRepo.Save(ctx, review)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to add review`, err)
	}
	return &responses.AddReviewResponse{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}, nil
}

var _ usecase_interfaces.ReviewUsecase = &reviewUsecase{}

func newReviewUsecase(
	reviewRepo repository_interfaces.ReviewRepository,
	cafeRepo repository_interfaces.CafeRepository,
	userRepo repository_interfaces.UserRepository,
) *reviewUsecase {
	return &reviewUsecase{
		reviewRepo: reviewRepo,
		cafeRepo:   cafeRepo,
		userRepo:   userRepo,
	}
}

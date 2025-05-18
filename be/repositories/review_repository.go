package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type reviewRepository struct {
	db repository_interfaces.DBX
}

func (r *reviewRepository) Save(ctx context.Context, entity *models.Review) (result *models.Review, err error) {
	id, err := mysql_client.Save[models.Review](ctx, r.db(ctx), queries.InsertIntoReview,
		&entity.UserID, &entity.CafeID, &entity.Score, &entity.Comment)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (r *reviewRepository) FindByID(ctx context.Context, id uint64) (result *models.Review, err error) {
	result = &models.Review{}
	err = mysql_client.FindByID[models.Review](ctx, r.db(ctx), queries.SelectReviewByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *reviewRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Review, err error) {
	return mysql_client.FindAll[models.Review](ctx, r.db(ctx), queries.SelectsReview)
}

func (r *reviewRepository) FindByCafeID(ctx context.Context, cafeID uint64) (result []*models.Review, err error) {
	rows, err := r.db(ctx).QueryContext(ctx, queries.SelectReviewsByCafeID.String(), cafeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		review := &models.Review{}
		err = rows.Scan(
			&review.ID,
			&review.UserID,
			&review.CafeID,
			&review.Score,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, review)
	}
	return result, nil
}

func (r *reviewRepository) AvgCafeRating(ctx context.Context, cafeID uint64) (result float64, err error) {
	err = r.db(ctx).QueryRowContext(ctx, queries.SelectReviewAvgByCafeID.String(), cafeID).Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

var _ repository_interfaces.ReviewRepository = &reviewRepository{}

func newReviewRepository(db repository_interfaces.DBX) *reviewRepository {
	return &reviewRepository{
		db: db,
	}
}

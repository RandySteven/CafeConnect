package queries

const (
	InsertIntoReview GoQuery = `
		INSERT INTO reviews (user_id, cafe_id, score, comment)
		VALUES 
		    (?, ?, ?, ?)
	`

	SelectsReview GoQuery = `
		SELECT id, user_id, cafe_id, score, comment, created_at, updated_at, deleted_at
		FROM
		    reviews
	`

	SelectReviewByID GoQuery = `
		SELECT id, user_id, cafe_id, score, comment, created_at, updated_at, deleted_at
		FROM
		    reviews
		WHERE id = ?
	`

	SelectReviewsByCafeID GoQuery = `
		SELECT id, user_id, cafe_id, score, comment, created_at, updated_at, deleted_at
		FROM
		    reviews
		WHERE cafe_id = ?
	`

	SelectReviewAvgByCafeID GoQuery = `
		SELECT AVG(score)
		FROM
		    reviews
		WHERE cafe_id = ?
	`
)

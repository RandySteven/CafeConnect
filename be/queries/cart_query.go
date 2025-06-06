package queries

const (
	InsertIntoCart GoQuery = `
		INSERT INTO carts (user_id, cafe_product_id, qty)
		VALUES 
		    (?, ?, ?)
	`

	SelectCartsByUserID GoQuery = `
		SELECT id, user_id, cafe_product_id, qty, created_at, updated_at, deleted_at
		FROM
		    carts
		WHERE user_id = ?
	`

	SelectCartByUserIDAndCafeProductID GoQuery = `
		SELECT id, user_id, cafe_product_id, qty, created_at, updated_at, deleted_at
		FROM
		    carts
		WHERE user_id = ? AND cafe_product_id = ?
	`

	UpdateCartByID GoQuery = `
		UPDATE carts
			SET 
			    user_id = ?,
			    cafe_product_id = ?,
			    qty = ?,
			    created_at = ?,
			    updated_at = ?,
			    deleted_at = ?
		WHERE
		    id = ?
	`

	DeleteCartByUserID GoQuery = `
		DELETE FROM carts WHERE user_id = ?
	`

	DeleteCartByUserIDAndCafeProductID GoQuery = `
		DELETE FROM carts WHERE 
		                      user_id = ? AND cafe_product_id = ?
	`
)

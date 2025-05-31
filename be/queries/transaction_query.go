package queries

const (
	InsertTransactionHeader GoQuery = `
		INSERT INTO transaction_headers (user_id, cafe_id, transaction_code, status, transaction_at)
		VALUES 
		    (?, ?, ?, ?, ?)
	`

	InsertTransactionDetail GoQuery = `
		INSERT INTO transaction_details (transaction_id, cafe_product_id, qty)
		VALUES 
		    (?, ?, ?)
	`

	SelectTransactionHeaderByTransactionCode GoQuery = `
		SELECT id, user_id, cafe_id, transaction_code, status, transaction_at, created_at, updated_at, deleted_at
		FROM
		    transaction_headers
		WHERE 
		    transaction_code = ?
	`

	SelectTransactionHeaderByUserID GoQuery = `
		SELECT id, user_id, cafe_id, transaction_code, status, transaction_at, created_at, updated_at, deleted_at
		FROM
		    transaction_headers
		WHERE 
			user_id = ?
	`

	SelectTransactionDetailsByTransactionId GoQuery = `
		SELECT id, transaction_id, cafe_product_id, qty, created_at, updated_at, deleted_at
		FROM
		    transaction_details
		WHERE 
		    transaction_id = ?
	`

	UpdateTransactionHeader GoQuery = `
		UPDATE transaction_headers SET
			user_id = ?,
			cafe_id = ?,
			transaction_code = ?,
			status = ?,
			transaction_at = ?,
		    created_at = ?,
			updated_at = ?,
			deleted_at = ?
		WHERE id = ?
	`
)

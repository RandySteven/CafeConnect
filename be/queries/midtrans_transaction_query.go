package queries

const (
	InsertMidtransTransaction GoQuery = `
		INSERT INTO midtrans_transactions (transaction_code, total_amt, token, redirect_url)
		VALUES 
		    (?, ?, ?, ?)
	`

	SelectMidtransByID GoQuery = `

	`

	SelectMidtransByTransactionCode GoQuery = `
		SELECT id, transaction_code, total_amt, token, redirect_url, created_at, updated_at, deleted_at
		FROM
		    midtrans_transactions
		WHERE
		    transaction_code = ?
	`
)

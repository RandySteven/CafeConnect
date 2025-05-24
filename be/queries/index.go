package queries

const (
	TransactionCodeIndex IndexQuery = `
		CREATE INDEX transaction_code
		ON transaction_headers (transaction_code)
	`

	DropIndex IndexQuery = `ALTER TABLE %s DROP INDEX %s`
)

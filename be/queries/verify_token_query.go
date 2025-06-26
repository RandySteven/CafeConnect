package queries

const (
	InsertVerifyToken GoQuery = `
		INSERT INTO verify_tokens (token, user_id, is_clicked, expired_time)
		VALUES 
		    (?, ?, ?, ?)
	`

	SelectVerifyTokens GoQuery = `
		SELECT id, token, user_id, is_clicked, expired_time, created_at, updated_at, deleted_at
		FROM
		    verify_tokens
	`

	SelectVerifyTokenByID GoQuery = `
		SELECT id, token, user_id, is_clicked, expired_time, created_at, updated_at, deleted_at
		FROM
		    verify_tokens
		WHERE id = ?
	`

	SelectVerifyTokenByToken GoQuery = `
		SELECT id, token, user_id, is_clicked, expired_time, created_at, updated_at, deleted_at
		FROM
		    verify_tokens
		WHERE token = ?
	`

	UpdateVerifyToken GoQuery = `
		UPDATE verify_tokens
			SET
			    token = ?,
			    user_id = ?,
			    is_clicked = ?,
			    expired_time = ?,
			    created_at = ?,
			    updated_at = ?,
			    deleted_at = ?
		WHERE id = ?
	`
)

package queries

const (
	InsertReferral GoQuery = `
		INSERT INTO referrals (code, user_id, expired_time, status)
		VALUES 
		    (?, ?, ?, ?)
	`

	SelectReferrals GoQuery = `
		SELECT id, code, user_id, expired_time, status, numb_of_usage, created_at, updated_at, deleted_at
		FROM
		    referrals
	`

	SelectReferralByID GoQuery = `
		SELECT id, code, user_id, expired_time, status, numb_of_usage, created_at, updated_at, deleted_at
		FROM
		    referrals
		WHERE id = ?
	`

	SelectReferralByUserID GoQuery = `
		SELECT id, code, user_id, expired_time, status, numb_of_usage, created_at, updated_at, deleted_at
		FROM
		    referrals
		WHERE user_id = ?
	`

	SelectReferralByCode GoQuery = `
		SELECT id, code, user_id, expired_time, status, numb_of_usage, created_at, updated_at, deleted_at
		FROM
		    referrals
		WHERE code = ?
	`

	UpdateReferralByID GoQuery = `
		UPDATE referrals
		SET
		    code = ?,
		    user_id = ?,
		    expired_time = ?,
		    status = ?,
		    numb_of_usage = ?,
		    created_at = ?,
		    updated_at = ?,
		    deleted_at = ?
		WHERE
		    id = ?
	`
)

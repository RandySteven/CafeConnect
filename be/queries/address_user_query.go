package queries

const (
	InsertAddressUser GoQuery = `
		INSERT INTO address_users (address_id, user_id)
		VALUES 
		    (?, ?)
	`

	SelectAddressUsers GoQuery = `
		SELECT id, address_id, user_id, is_default, created_at, updated_at, deleted_at
		FROM
		    address_users
	`

	SelectAddressUserByID GoQuery = `
		SELECT id, address_id, user_id, is_default, created_at, updated_at, deleted_at
		FROM
			address_users
		WHERE id = ?
	`

	SelectAddressUserByUserID GoQuery = `
		SELECT id, address_id, user_id, is_default, created_at, updated_at, deleted_at
		FROM
			address_users
		WHERE user_id = ?
	`

	SelectAddressUserByAddressAndUserID GoQuery = `
		SELECT id, address_id, user_id, is_default, created_at, updated_at, deleted_at
		FROM
			address_users
		WHERE address_id = ? AND user_id = ?	
	`
)

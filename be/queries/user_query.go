package queries

const (
	InsertUser GoQuery = `
		INSERT INTO users (name, username, email, password, profile_picture, phone_number, dob)
		VALUES 
		    (?, ?, ?, ?, ?, ?)
	`

	SelectUsers GoQuery = `
		SELECT id, name, username, email, password, phone_number, profile_picture, dob, created_at, updated_at, deleted_at
		FROM
		    users
	`

	SelectUserByID GoQuery = `
		SELECT id, name, username, email, password, phone_number, profile_picture, dob, created_at, updated_at, deleted_at
		FROM
		    users
		WHERE
		    id = ?
	`

	SelectPhoneNumber GoQuery = `
		SELECT id, name, username, email, password, phone_number, profile_picture, dob, created_at, updated_at, deleted_at
		FROM
		    users
		WHERE
		    phone_number = ?
	`

	SelectEmail GoQuery = `
		SELECT id, name, username, email, password, phone_number, profile_picture, dob, created_at, updated_at, deleted_at
		FROM
		    users
		WHERE
		    email = ?
	`

	SelectUsername GoQuery = `
		SELECT id, name, username, email, password, phone_number, profile_picture, dob, created_at, updated_at, deleted_at
		FROM
		    users
		WHERE
		    username = ?
	`
)

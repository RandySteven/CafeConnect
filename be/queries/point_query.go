package queries

const (
	InsertPoint GoQuery = `
		INSERT INTO points (point, user_id)
		VALUES 
		    (?, ?)
	`

	SelectPoints GoQuery = `
		SELECT id, point, user_id, created_at, updated_at, deleted_at
		FROM
		    points
	`

	SelectPointByID GoQuery = `
		SELECT id, point, user_id, created_at, updated_at, deleted_at
		FROM
		    points
		WHERE id = ?
	`

	UpdatePointByID GoQuery = `
		UPDATE points SET
			point = ?,
			user_id = ?,
			created_at = ?,
			updated_at = ?,
			deleted_at = ?
		WHERE id = ?
	`
)

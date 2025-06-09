package queries

const (
	InsertRoleQuery GoQuery = `
		INSERT INTO roles (role)
		VALUES 
		    (?)
	`

	SelectRolesQuery GoQuery = `
		SELECT id, role, created_at, updated_At, deleted_At
		FROM
		    roles
	`

	SelectRoleByIDQuery GoQuery = `
		SELECT id, role, created_at, updated_At, deleted_At
		FROM
		    roles
		WHERE id = ?
	`
)

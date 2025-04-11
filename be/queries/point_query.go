package queries

const (
	InsertPoint GoQuery = `
		INSERT INTO points (point, user_id)
		VALUES 
		    (?, ?)
	`
)

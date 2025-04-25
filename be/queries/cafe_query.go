package queries

const (
	InsertIntoCafe GoQuery = `
		INSERT INTO cafes (address_id, cafe_franchise_id, cafe_type, photo_urls, open_hour, close_hour)
		VALUES 
		    (?, ?, ?, ?, ?, ?)
	`
)

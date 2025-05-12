package queries

const (
	InsertIntoCafe GoQuery = `
		INSERT INTO cafes (address_id, cafe_franchise_id, cafe_type, photo_urls, open_hour, close_hour)
		VALUES 
		    (?, ?, ?, ?, ?, ?)
	`

	SelectCafeByID GoQuery = `
		SELECT id, address_id, cafe_franchise_id, cafe_type, photo_urls, open_hour, close_hour, created_at, updated_at, deleted_at
		FROM	
		    cafes
		WHERE
		    id = ?
	`

	SelectCafesByCafeFranchiseID GoQuery = `
		SELECT id, address_id, cafe_franchise_id, cafe_type, photo_urls, open_hour, close_hour, created_at, updated_at, deleted_at
		FROM	
		    cafes
		WHERE
			cafe_franchise_id = ?
	`

	SelectCafeByAddressID GoQuery = `
		SELECT id, address_id, cafe_franchise_id, cafe_type, photo_urls, open_hour, close_hour, created_at, updated_at, deleted_at
		FROM	
		    cafes
		WHERE
			address_id = ?
	`
)

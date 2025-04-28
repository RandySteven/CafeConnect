package queries

const (
	InsertIntoProduct GoQuery = `
		INSERT INTO products (name, photo_url, product_category_id)
		VALUES 
		    (?, ?, ?)
	`

	SelectProducts GoQuery = `
		SELECT id, name, photo_url, product_category_id, created_at, updated_at, deleted_at
		FROM
		    products
	`

	SelectProductByID GoQuery = `
		SELECT id, name, photo_url, product_category_id, created_at, updated_at, deleted_at
		FROM
		    products
		WHERE id = ?
	`
)

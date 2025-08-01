package queries

const (
	InsertCafeProduct GoQuery = `
		INSERT INTO cafe_products (cafe_id, product_id, price)
		VALUES 
		    (?, ?, ?)
	`

	SelectCafeProducts GoQuery = `
		SELECT id, cafe_id, product_id, price, stock, status, created_at, updated_at, deleted_at
		FROM
		    cafe_products
	`

	SelectCafeProductByID GoQuery = `
		SELECT 
		    id, cafe_id, product_id, price, stock, status, created_at, updated_at, deleted_at
		FROM
		    cafe_products 
		WHERE id = ?
	`

	SelectCafeProductsByCafeID GoQuery = `
		SELECT id, cafe_id, product_id, price, stock, status, created_at, updated_at, deleted_at
		FROM
		    cafe_products
		WHERE
		    cafe_id = ?
	`

	SelectCafeProductsInCafeIDs GoQuery = `
		SELECT id, cafe_id, product_id, price, stock, status, created_at, updated_at, deleted_at
		FROM
		    cafe_products
		WHERE
		    cafe_id IN 
	`

	SelectCafeIdByCafeProductIDs GoQuery = `
		SELECT cafe_id FROM cafe_products WHERE id IN %s GROUP BY cafe_id;
	`

	UpdateCafeProductByID GoQuery = `
		UPDATE cafe_products
		SET
		    cafe_id = ?,
		    product_id = ?,
		    price = ?,
		    stock = ?,
		    status = ?,
		    created_at = ?,
		    updated_at = ?,
		    deleted_at = ?
		WHERE
		    id = ?
	`
)

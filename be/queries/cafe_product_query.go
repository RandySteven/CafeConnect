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
		    cp.id, cp.cafe_id, cp.product_id, cp.price, cp.stock, cp.status, cp.created_at, cp.updated_at, cp.deleted_at,
		    p.id, p.name, p.photo_url, p.product_category_id, p.created_at, p.updated_at, p.deleted_at,
		    c.id, c.address_id, c.cafe_franchise_id, c.cafe_type,
		FROM
		    cafe_products cp
		INNER JOIN
		    cafes c
		ON cp.cafe_id = c.id
		INNER JOIN
		    products p 
		ON cp.product_id, p.id
		WHERE id = ?
	`
)

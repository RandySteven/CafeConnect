package queries

const (
	SelectProductCategories GoQuery = `
		SELECT id, category, created_at, updated_at, deleted_at
		FROM
		    product_categories
	`

	SelectProductCategoriesByID GoQuery = `
		SELECT id, category, created_at, updated_at, deleted_at
		FROM
		    product_categories
		WHERE id = ?
	`
)

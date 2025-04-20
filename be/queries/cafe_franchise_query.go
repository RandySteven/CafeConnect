package queries

const (
	SelectCafeFranchises GoQuery = `
		SELECT id, name, logo_url, created_at, updated_at, deleted_at
		FROM
		    cafe_franchises
	`
)

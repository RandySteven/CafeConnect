package queries

const (
	InsertAddress GoQuery = `
		INSERT INTO addresses (address, coordinate)
		VALUES 
		    (?, POINT(?, ?))
	`

	SelectAddressByID GoQuery = `
		SELECT 
		  id,
		  address,
		  ST_X(coordinate) AS longitude,
		  ST_Y(coordinate) AS latitude,
		  created_at, updated_at, deleted_at
		FROM addresses
		WHERE id = ?
	`

	SelectAddressByRadiusNKm GoQuery = `
		SELECT 
			id,
			address,
			ST_X(coordinate) AS longitude,
			ST_Y(coordinate) AS latitude,
			created_at,
			updated_at,
			deleted_at
		FROM addresses
		WHERE ST_Distance_Sphere(coordinate, POINT(?, ?)) <= ?
	` // longitude, latitude
)

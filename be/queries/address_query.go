package queries

const (
	InsertAddress = `
		INSERT INTO addresses (address, coordinate)
		VALUES 
		    (?, POINT(?, ?))
	`
)

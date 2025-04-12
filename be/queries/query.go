package queries

type (
	GoQuery        string
	MigrationQuery string
)

func (q GoQuery) String() string {
	return string(q)
}

func (m MigrationQuery) String() string {
	return string(m)
}

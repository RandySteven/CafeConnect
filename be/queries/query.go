package queries

type (
	GoQuery        string
	MigrationQuery string
	IndexQuery     string
)

func (q GoQuery) String() string {
	return string(q)
}

func (m MigrationQuery) String() string {
	return string(m)
}

func (i IndexQuery) String() string {
	return string(i)
}

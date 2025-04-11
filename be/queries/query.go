package queries

type (
	GoQuery string
)

func (q GoQuery) String() string {
	return string(q)
}

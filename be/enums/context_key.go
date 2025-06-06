package enums

type ContextKey string

const (
	UserID     ContextKey = `user_id`
	RoleID                = `role_id`
	RequestID             = `request_id`
	Env                   = `env`
	ClientIP              = `client_ip`
	FileHeader            = `file_header`
	FileObject            = `file_object`
	QtyCart               = `qty_cart`
	QtyTrx                = `qty_trx`
)

func (c ContextKey) ToString() string {
	return string(c)
}

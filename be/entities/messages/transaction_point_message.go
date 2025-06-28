package messages

type TransactionPointMessage struct {
	UserID uint64 `json:"user_id"`
	Point  uint64 `json:"point_id"`
}

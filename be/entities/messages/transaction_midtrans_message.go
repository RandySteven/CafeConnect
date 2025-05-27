package messages

type TransactionMidtransMessage struct {
	UserID         uint64   `json:"user_id"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	PhoneNumber    string   `json:"phone_number"`
	CafeProductIDs []uint64 `json:"cafe_product_ids"`
}

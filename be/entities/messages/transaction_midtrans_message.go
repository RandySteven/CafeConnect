package messages

import "github.com/RandySteven/CafeConnect/be/entities/payloads/requests"

type TransactionMidtransMessage struct {
	UserID            uint64                   `json:"user_id"`
	FName             string                   `json:"f_name"`
	Email             string                   `json:"email"`
	Phone             string                   `json:"phone"`
	LName             string                   `json:"l_name"`
	CafeID            uint64                   `json:"cafe_id"`
	TransactionCode   string                   `json:"transaction_code"`
	CafeFranchiseName string                   `json:"cafe_franchise_name"`
	CheckoutList      []*requests.CheckoutList `json:"checkout_list"`
}

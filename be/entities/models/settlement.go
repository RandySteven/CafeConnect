package models

import "time"

type Settlement struct {
	ID                    uint64
	MidtransTransactionID string
	TransactionCode       string
	FranchiseID           uint64
	CafeID                uint64
	GrossAmount           string
	PaymentType           string
	SettlementStatus      string
	Metadata              string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

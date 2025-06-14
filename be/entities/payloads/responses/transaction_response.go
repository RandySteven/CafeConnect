package responses

import (
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"time"
)

type (
	TransactionDetailItem struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		Price    uint64 `json:"price"`
		ImageURL string `json:"image_url"`
		Qty      uint64 `json:"qty"`
	}

	TransactionReceiptResponse struct {
		ID               uint64                            `json:"id"`
		TransactionCode  string                            `json:"transaction_code"`
		Status           string                            `json:"status"`
		TransactionAt    time.Time                         `json:"transaction_at"`
		MidtransResponse *midtrans_client.MidtransResponse `json:"midtrans_response"`
	}

	TransactionDetailResponse struct {
		ID              uint64                   `json:"id"`
		TransactionCode string                   `json:"transaction_code"`
		TransactionTime time.Time                `json:"transaction_at"`
		Status          string                   `json:"status"`
		CreatedAt       time.Time                `json:"created_at"`
		UpdatedAt       time.Time                `json:"updated_at"`
		Items           []*TransactionDetailItem `json:"items"`
	}

	CafeResponse struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		ImageURL string `json:"image_url"`
	}

	TransactionListResponse struct {
		ID              uint64        `json:"id"`
		Cafe            *CafeResponse `json:"cafe"`
		TransactionCode string        `json:"transaction_code"`
		Status          string        `json:"status"`
		TransactionAt   time.Time     `json:"transaction_at"`
		CreatedAt       time.Time     `json:"created_at"`
		UpdatedAt       time.Time     `json:"updated_at"`
		DeletedAt       *time.Time    `json:"deleted_at"`
	}
)

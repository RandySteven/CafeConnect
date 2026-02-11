package requests

type (
	CheckoutList struct {
		CafeProductID uint64 `json:"cafe_product_id"`
		Qty           uint64 `json:"qty"`
	}

	CreateTransactionRequest struct {
		IdempotencyKey string          `json:"idempotency_key"`
		CafeID         uint64          `json:"cafe_id"`
		Checkouts      []*CheckoutList `json:"checkouts"`
	}

	ReceiptRequest struct {
		TransactionCode string `json:"transaction_code"`
	}

	MidtransNotificationRequest struct {
		TransactionStatus string `json:"transaction_status"`
		OrderID           string `json:"order_id"`
		FraudStatus       string `json:"fraud_status"`
		PaymentType       string `json:"payment_type"`
		TransactionTime   string `json:"transaction_time"`
		TransactionID     string `json:"transaction_id"`
		GrossAmount       string `json:"gross_amount"`
		SignatureKey      string `json:"signature_key"`
		StatusCode        string `json:"status_code"`
		MerchantID        string `json:"merchant_id"`
		SettlementTime    string `json:"settlement_time,omitempty"`
		Bank              string `json:"bank,omitempty"`
		VAAccountNumber   string `json:"va_numbers,omitempty"`
		PaymentCode       string `json:"payment_code,omitempty"`
		BillerCode        string `json:"biller_code,omitempty"`
		Store             string `json:"store,omitempty"`
		ExpiryTime        string `json:"expiry_time,omitempty"`
	}

	PaymentConfirmationRequest struct {
		TransactionCode string `json:"transaction_code"`
	}
)

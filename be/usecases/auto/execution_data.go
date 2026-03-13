package auto_transfer_usecases

import (
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
)

// DeductedProduct tracks a stock deduction that was applied, so it can
// be rolled back by the restoreStock compensation activity.
type (
	DeductedProduct struct {
		CafeProductID uint64 `json:"cafe_product_id"`
		Qty           uint64 `json:"qty"`
	}

	ExecutionData struct {
		// Request inputs data
		UserID  uint64                             `json:"user_id"`
		Request *requests.CreateTransactionRequest `json:"request"`

		// Intermediate results — populated by activities
		User              *models.User                         `json:"user,omitempty"`
		Cafe              *models.Cafe                         `json:"cafe,omitempty"`
		Franchise         *models.CafeFranchise                `json:"franchise,omitempty"`
		TransactionHeader *models.TransactionHeader            `json:"transaction_header,omitempty"`
		MidtransMessage   *messages.TransactionMidtransMessage `json:"midtrans_message,omitempty"`
		MidtransResponse  *midtrans_client.MidtransResponse    `json:"midtrans_response,omitempty"`

		// Branching — controls which activity runs next.
		// If empty, Execute follows the default sequential order.
		// Set to an activity name to branch (e.g., for compensation).
		CurrentActivity string `json:"current_activity,omitempty"`
		NextActivity    string `json:"next_activity,omitempty"`

		// Compensation tracking
		DeductedProducts     []*DeductedProduct `json:"deducted_products,omitempty"`
		StockDeductionFailed bool               `json:"stock_deduction_failed,omitempty"`
	}
)

func (t *ExecutionData) GetActivity() string     { return t.CurrentActivity }
func (t *ExecutionData) SetActivity(name string) { t.CurrentActivity = name }
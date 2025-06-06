// type (
//     TransactionDetailItem struct {
//     ID       uint64 `json:"id"`
//     Name     string `json:"name"`
//     Price    uint64 `json:"price"`
//     ImageURL string `json:"image_url"`
//     Qty      uint64 `json:"qty"`
// }
//
// TransactionReceiptResponse struct {
//     ID              uint64    `json:"id"`
//     TransactionCode string    `json:"transaction_code"`
//     Status          string    `json:"status"`
//     TransactionAt   time.Time `json:"transaction_at"`
// }
//
// TransactionDetailResponse struct {
//     ID              uint64                   `json:"id"`
//     TransactionCode string                   `json:"transaction_code"`
//     TransactionTime time.Time                `json:"transaction_at"`
//     Status          string                   `json:"status"`
//     CreatedAt       time.Time                `json:"created_at"`
//     UpdatedAt       time.Time                `json:"updated_at"`
//     Items           []*TransactionDetailItem `json:"items"`
// }
//
// TransactionListResponse struct {
//     ID              uint64     `json:"id"`
//     TransactionCode string     `json:"transaction_code"`
//     Status          string     `json:"status"`
//     TransactionAt   time.Time  `json:"transaction_at"`
//     CreatedAt       time.Time  `json:"created_at"`
//     UpdatedAt       time.Time  `json:"updated_at"`
//     DeletedAt       *time.Time `json:"deleted_at"`
// }
// )

interface TransactionDetailItem {
    id: number
    name: string
    price: number
    image_url: string
    qty: number
}

export interface TransactionReceiptResponse {

}

export interface TransactionDetailResponse {
    id: number
    transaction_code: string
    transaction_time: string
    status: string
    created_at: string
    updated_at: string
    items: TransactionDetailItem[]
}

interface CafeResponse {
    id: number
    name: string
    address: string
    image_url: string
}

export interface TransactionListResponse {
    id: number
    transaction_code: string
    cafe: CafeResponse
    transaction_time: string
    status: string
    transaction_at: string
    created_at: string
    updated_At: string
    deleted_at: string
}
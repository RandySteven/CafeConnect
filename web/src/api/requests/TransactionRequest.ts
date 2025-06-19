interface CheckoutItem {
    cafe_product_id: number
    qty: number
}

export interface CheckoutV2Request {
    cafe_id: number
    checkouts: CheckoutItem[]
}

export interface ReceiptRequest {
    transaction_code: string
}
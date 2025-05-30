interface CheckoutList {
    cafe_id : number
    cafe_name : string
    items : CafeCartItems[]
}

interface CafeCartItems {
    product_id : number
    product_name : string
    product_price : number
    product_image : string
    qty : number
    created_at : string
    updated_at : string
    deleted_at : string
}

export interface ListCartResponse {
    user_id : number
    checkout_list: CheckoutList[]
}
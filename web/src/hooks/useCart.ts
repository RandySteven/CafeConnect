import {ListCartResponse} from "@/api/responses/CartResponse";
import {useEffect, useState} from "react";
import {GET} from "@/api/api";
import {GET_CARTS} from "@/api/endpoint";

export const useCart = () : ListCartResponse => {
    const [cartResponse, setCartResponse] = useState<ListCartResponse>({
        user_id: 0,
        checkout_list: []
    })

    useEffect(() => {
        const fetch = async () => {
            const result = await GET(GET_CARTS, true)
            setCartResponse(result.data.cart)
        }

        fetch()
    }, []);

    return cartResponse
}
import {useEffect, useState} from "react";
import {ListProductResponse} from "@/api/responses/ProductResponse";
import {POST} from "@/api/api";

export const useProductByCafe = (cafeId: number) : ListProductResponse[] => {
    const [listProductResponse, setListProductResponse] = useState<ListProductResponse[]>([])

    useEffect(() => {
        const fetchProducts = async () => {
            const result = await POST("products/cafe-list", false, {
                cafe_id: Number(cafeId)
            })
            if(result.data.menus != null) {
                setListProductResponse(result.data.menus)
            }
        }

        fetchProducts()
    }, []);

    return listProductResponse
}
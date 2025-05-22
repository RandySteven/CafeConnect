import {useEffect, useState} from "react";
import {ListProductResponse} from "@/api/responses/ProductResponse";
import {POST} from "@/api/api";

export const useProductByCafe = (cafeIds: number[]) : ListProductResponse[] => {
    const [listProductResponse, setListProductResponse] = useState<ListProductResponse[]>([])
    console.log(`request : `, cafeIds)
    useEffect(() => {
        const fetchProducts = async () => {
            const result = await POST("products/cafe-list", false, {
                cafe_ids: cafeIds
            })
            if(result.data.menus != null) {
                setListProductResponse(result.data.menus)
            }
        }

        fetchProducts()
    }, []);

    return listProductResponse
}
"use client";

import {ListCafeResponse} from "@/api/responses/CafeResponse";
import {useEffect, useState} from "react";
import {ListRadiusCafeRequest} from "@/api/requests/CafeRequest";
import {POST} from "@/api/api";

export const useListCafeWithRadius = (longitude, latitude, radius : number) : ListCafeResponse[] => {
    const [listCafeResponse, setListCafeResponse] = useState<ListCafeResponse[]>([])

    const listCafeRequest: ListRadiusCafeRequest = {
        point: {
            longitude: longitude,
            latitude: latitude
        },
        radius: radius
    }

    useEffect( () => {

        console.log(listCafeRequest)
        const fetchCafes = async () => {
            try {
                const result = await POST("cafes", false, listCafeRequest)
                console.log(`resultnya `, result.data.cafes)
                setListCafeResponse(result.data.cafes)
            }catch (err) {
                console.log(err)
            }
        };

        fetchCafes()
    }, []);

    return listCafeResponse
}
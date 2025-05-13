"use client";

import {ListCafeResponse} from "@/api/responses/CafeResponse";
import {useEffect, useState} from "react";
import {ListRadiusCafeRequest} from "@/api/requests/CafeRequest";
import {POST} from "@/api/api";

export const useListCafeWithRadius = (address_id, radius : number) : ListCafeResponse  => {
    const [listCafeResponse, setListCafeResponse] = useState<ListCafeResponse[]>([])

    useEffect(() => {
        const listCafeRequest : ListRadiusCafeRequest = {
            address_id: address_id,
            radius: radius
        }

        POST(`/cafes`, false, listCafeRequest)
            .then((data) => {
                return setListCafeResponse(data.data.cafes)
            })
            .catch((error) => {
                return error
            })

    }, []);

    return listCafeResponse
}
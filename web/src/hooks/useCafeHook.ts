"use client";

import {ListCafeResponse} from "@/api/responses/CafeResponse";
import {useEffect, useState} from "react";

export const useListCafeWithRadius = (radius : number) : ListCafeResponse  => {
    const [listCafeResponse, setListCafeResponse] = useState<ListCafeResponse[]>([])

    useEffect(() => {
        const listCafeRequest = ListCa

    }, []);

    return
}
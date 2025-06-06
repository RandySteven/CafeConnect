"use client";

import {DetailCafeResponse, ListCafeResponse} from "@/api/responses/CafeResponse";
import {useEffect, useState} from "react";
import {ListRadiusCafeRequest} from "@/api/requests/CafeRequest";
import {GET, POST} from "@/api/api";
import {GET_CAFE_RADIUS} from "@/api/endpoint";

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

        const fetchCafes = async () => {
            try {
                const result = await POST(GET_CAFE_RADIUS, false, listCafeRequest)
                if(result.data.cafes != null) {
                    setListCafeResponse(result.data.cafes)
                }
            }catch (err) {
                console.log(err)
            }
        };

        fetchCafes()
    }, []);

    return listCafeResponse
}

export const useCafeDetail = (id : number) : DetailCafeResponse => {
    const [detailCafeResponse, setDetailCafeResponse] = useState<DetailCafeResponse>({})

    useEffect(() => {
        const fetchResult = async () => {
            try　{
                const result = await GET(`cafes/${id}`, false)
                setDetailCafeResponse(result.data.cafe)
            }catch (err) {
                console.log(err)
            }
        }

        fetchResult()
    }, []);

    return detailCafeResponse
}
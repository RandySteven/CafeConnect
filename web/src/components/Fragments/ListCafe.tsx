"use client";

import {Fragment, useEffect, useState} from "react";
import {ListCard} from "@/components/Elements/Card";
import {Box} from "@mui/material";
import {setItem} from "@/utils/common";
import {ListCafeResponse} from "@/api/responses/CafeResponse";
import {POST} from "@/api/api";
import {GET_CAFE_RADIUS} from "@/api/endpoint";

export const ListCafe = () => {
    const [long, setLong] = useState<number | null>(null);
    const [lat, setLat] = useState<number | null>(null);
    const [listCafes, setListCafes] = useState<ListCafeResponse[]>([])

    useEffect(() => {
        if (typeof window !== "undefined") {
            const longValue = parseFloat(localStorage.getItem("long") ?? "0");
            const latValue = parseFloat(localStorage.getItem("lat") ?? "0");

            const fetch = async () => {
                const result = await POST(GET_CAFE_RADIUS, false, {
                        point: {
                            longitude: longValue?? 0,
                            latitude: latValue?? 0
                        },
                        radius : 2500
                    }
                );

                setListCafes(result.data.cafes)
            }

            fetch()
        }
    }, []);

    useEffect(() => {
        if (listCafes.length > 0 && long !== null && lat !== null) {
            const ids = listCafes.map((cafe) => cafe.id);
            setItem("cafeIds", ids);
        }
    }, [listCafes, long, lat]);

    return <Fragment>
        <Box
            sx={{
                py: 2
            }}
        >
            {
                listCafes.map((cafe, index) => (
                    <ListCard
                        link={`/cafes/${cafe.id}`}
                        key={index}
                        type="cafe"
                        img={cafe.logo_url}
                        name={cafe.name}
                        logoURL={cafe.logo_url}
                        status={cafe.status}
                        openHour={cafe.open_hour}
                        closeHour={cafe.close_hour}
                        address={cafe.address}
                    />
                ))
            }
        </Box>
    </Fragment>
}
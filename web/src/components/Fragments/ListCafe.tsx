"use client";

import {Fragment, useEffect, useState} from "react";
import {ListCard} from "@/components/Elements/Card";
import {Box} from "@mui/material";
import {setItem} from "@/utils/common";
import {ListCafeResponse} from "@/api/responses/CafeResponse";
import {POST} from "@/api/api";
import {GET_CAFE_RADIUS} from "@/api/endpoint";
import {Coordinate, useAWSLocationServiceMultiPinPoint} from "@/hooks/useGeoLoc";
import {AWSMap} from "@/components/Fragments/AWSMap";

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

    let coordinates : Coordinate[] = []

    return <Fragment>
        <Box
            sx={{
                py: 2
            }}
        >
            <h1>
                List Cafe Terdekat Rumah Anda
            </h1>
            {
                listCafes.map((cafe, index) => {
                    coordinates.push({
                        lat: cafe.address.latitude,
                        long: cafe.address.longitude
                    })
                    return (
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
                            address={cafe.address.address}
                        />
                    )
                })
            }
        </Box>
        <CafeListMaps coordinates={coordinates} />
    </Fragment>
}

export const CafeListMaps = (props : {
    coordinates: Coordinate[]
}) => {
    const href = useAWSLocationServiceMultiPinPoint(props.coordinates)
    return (
        <div ref={href} style={{width: "100%", height: "500px"}}/>
    )
}
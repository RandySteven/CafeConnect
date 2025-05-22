"use client";

import {Fragment} from "react";
import {useListCafeWithRadius} from "@/hooks/useCafeHook";
import {ListCard} from "@/components/Elements/Card";
import {Box} from "@mui/material";
import {setItem} from "@/utils/common";

export const ListCafe = () => {
    const listCafes = useListCafeWithRadius(109.34695696792511, -0.03646908129222186, 2500)
    let cafeIds : number[] = []
    for(let i = 0 ; i < listCafes.length ; i++) {
        cafeIds[i] = listCafes[i].id
    }

    setItem("cafeIds", cafeIds)

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
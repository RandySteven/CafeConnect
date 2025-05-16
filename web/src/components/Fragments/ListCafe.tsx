"use client";

import {Fragment} from "react";
import {useListCafeWithRadius} from "@/hooks/useCafeHook";
import {ListCard} from "@/components/Elements/Card";
import {Box} from "@mui/material";

export const ListCafe = () => {
    const listCafes = useListCafeWithRadius(109.325903544546240, 0, 25000)
    return <Fragment>
        <Box
            sx={{
                py: 2
            }}
        >
            {
                listCafes.map((cafe, index) => (
                    <ListCard
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
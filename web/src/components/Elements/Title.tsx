import { Fragment } from "react";
import { Box, Typography } from "@mui/material";
import {useDummy} from "@/hooks/useDummy";
import {CafeTitleProp} from "@/interfaces/props/TitleProp";
import {CafeListMaps} from "@/components/Fragments/ListCafe";
import {Coordinate, useAWSLocationService, useAWSLocationServiceSinglePinPoint} from "@/hooks/useGeoLoc";

export const AppTitle = () => {
    return (
        <Fragment>
            <Box display="flex" alignItems="center" >
                <img
                    src="/assets/img/cafeConnect-logo.png"
                    alt="CafeConnect Logo"
                    style={{ height: 40, marginRight: 8 }}
                />
            </Box>
        </Fragment>
    );
};

export const DummyTitle = () => {
    useDummy()

    return (
        <Fragment>
            <Box textAlign="center">
                <h1>INI TITLE</h1>
            </Box>
        </Fragment>
    );
};

export const CafeTitle = (prop : CafeTitleProp) => {
    console.log(`prop : `, prop)
    // const coordinates : Coordinate[] = [{
    //     long: prop.coordinate.longitude,
    //     lat: prop.coordinate.latitude,
    // }]
    return (
        <Fragment>
            <Box
                display="flex"
                flexDirection="column"
                alignItems="center"
                justifyContent="center"
                textAlign="center"
            >
                <img
                    src={prop.img}
                    alt={prop.name}
                    style={{
                        height: 240,
                        width: 240,
                        objectFit: "cover",
                        marginBottom: 8,
                    }}
                />
                <Typography variant="h6">{prop.name}</Typography>
                <Typography variant="body2">
                    {prop.address}
                </Typography>
                <CafeMap longitude={prop.longitude} latitude={prop.latitude} />
            </Box>
        </Fragment>
    )
}

export const CafeMap = (prop : {
    longitude: number
    latitude: number
}) => {
    const href = useAWSLocationServiceSinglePinPoint({
        long: prop.longitude,
        lat: prop.latitude,
    })
    return (
        <div ref={href} style={{width: "100%", height: "500px"}}/>
    )
}
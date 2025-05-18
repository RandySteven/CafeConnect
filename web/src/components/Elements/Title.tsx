import { Fragment } from "react";
import { Box, Typography } from "@mui/material";
import {useDummy} from "@/hooks/useDummy";
import {CafeTitleProp} from "@/interfaces/props/TitleProp";

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
                        height: 80,
                        width: 80,
                        borderRadius: "50%", // full circle
                        objectFit: "cover",
                        marginBottom: 8,
                    }}
                />
                <Typography variant="h6">{prop.name}</Typography>
                <Typography variant="body2">
                    {prop.address}
                </Typography>
            </Box>
        </Fragment>
    )
}
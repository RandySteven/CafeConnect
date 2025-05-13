import { Fragment } from "react";
import { Box, Typography } from "@mui/material";

export const AppTitle = () => {
    return (
        <Fragment>
            <Box display="flex" alignItems="center" >
                <img
                    src="./assets/img/cafeConnect-logo.png"
                    alt="CafeConnect Logo"
                    style={{ height: 40, marginRight: 8 }}
                />
            </Box>
        </Fragment>
    );
};

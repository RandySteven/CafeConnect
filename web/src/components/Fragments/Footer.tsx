import {Fragment} from "react";
import {BottomNavigation} from "@mui/material";

export const Footer = () => {
    return <Fragment>
        <BottomNavigation
            sx={{
                bgcolor: "#C38844",
                bottom: 0,
                left: 0,
                right: 0,
                zIndex: 1200,
                position: "fixed"
            }}
        >

        </BottomNavigation>
    </Fragment>
}
import React, { Fragment, ReactNode } from "react";
import { Box } from "@mui/material";

type ContainerProps = {
    children: ReactNode;
};

export const Container = ({ children }: ContainerProps) => {
    return (
        <Fragment>
            <Box
                sx={{
                    bgcolor: "#fff",
                    px: { xs: 2, sm: 3, md: 4 }, // responsive horizontal padding
                    maxWidth: "1200px",
                    mx: "auto", // center horizontally
                    width: "100%",
                    maxHeight: "inherit"
                }}
            >
                {children}
            </Box>
        </Fragment>
    );
};

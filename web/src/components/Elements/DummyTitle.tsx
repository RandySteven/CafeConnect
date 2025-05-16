"use client";

import { Fragment } from "react";
import { Box } from "@mui/material";
import {useDummy} from "@/hooks/useDummy";

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

"use client";

import {Fragment} from "react";
import {AppBar, Box, Container, Toolbar} from "@mui/material";
import {NavbarElementContent} from "@/components/Elements/NavbarElement";
import {useNavbarContent} from "@/hooks/useNavbarContentHook";
import {AppTitle} from "@/components/Elements/Title";

export const Navbar = () => {
    const navbar = useNavbarContent()
    return <Fragment>
        <AppBar
            position="fixed"
            sx={{
                bgcolor: "#C38844",
                py: 1,
            }}
        >
            <Toolbar>
                <Box
                    sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                        width: "100%",
                        px: 5
                    }}
                >
                    <AppTitle />
                    <NavbarElementContent navbarContents={navbar} />
                </Box>
            </Toolbar>
        </AppBar>
    </Fragment>
}
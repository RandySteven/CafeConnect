"use client";

import {Fragment} from "react";
import {Box, Card, CardContent, List} from "@mui/material";
import {useMenu} from "@/hooks/useMenu";
import {MenuCard} from "@/components/Elements/Card";
import {Grid} from "@mui/system";

export const GridMenu = () => {
    const menus = useMenu()
    return (
        <Fragment>
            <Box alignItems="center">
                <Grid container spacing={2}>
                    {menus.menus.map((menu, idx) => (
                        <Grid item xs={12} sm={6} md={3} key={idx}>
                            <MenuCard name={menu.menu} icon={menu.icon} link={menu.link} />
                        </Grid>
                    ))}
                </Grid>
            </Box>
        </Fragment>
    );
}
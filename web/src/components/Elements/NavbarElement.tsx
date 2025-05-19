import {Fragment} from "react";
import {NavbarProp} from "@/interfaces/props/NavbarProp";
import {Box, Button} from "@mui/material";

export const NavbarElementContent = (props : NavbarProp) => {
    return <Fragment>
        <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
            {props.navbarContents.contents.map((item) => (
                <Button key={item.href} sx={{ color: '#fff' }}>
                    {item.title}
                </Button>
            ))}
            <Button sx={{
                color: '#fff'
            }}>
                Login
            </Button>
        </Box>
    </Fragment>
}
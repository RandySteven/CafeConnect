import React, {Fragment, use} from "react";
import {NavbarProp} from "@/interfaces/props/NavbarProp";
import {Avatar, Box, Button, IconButton, Link, Menu, MenuItem} from "@mui/material";
import {getToken} from "@/utils/common";
import {useOnboarding} from "@/hooks/useOnboardingHook";
import {redirect} from "next/navigation";
import {red} from "@mui/material/colors";

export const NavbarElementContent = (props : NavbarProp) => {
    return <Fragment>
        <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
            {props.navbarContents.contents.map((item) => (
                    <Button key={item.href}>
                        <Link href={item.href}  sx={{ color: '#fff' }}>
                            {item.title}
                        </Link>
                    </Button>
            ))}
           <UserAccountMenu />
        </Box>
    </Fragment>
}

export const getOnboardUser = () => {
}

export const UserAccountMenu = () => {
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

    const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleProfile = () => {
        redirect(`/profile`)
    }

    const handleHistory = () => {
        redirect(`/histories`)
    }

    let user = useOnboarding()
    return <Fragment>
        <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            color="inherit"
        >
            <Avatar src={user.profile_picture} alt={user.name} />
        </IconButton>
        <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
                vertical: 'top',
                horizontal: 'right',
            }}
            keepMounted
            transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
            }}
            open={Boolean(anchorEl)}
            onClose={handleClose}
        >
            <MenuItem onClick={handleProfile}>Profile</MenuItem>
            <MenuItem onClick={handleHistory}>History</MenuItem>
            <MenuItem onClick={handleMenu}>Logout</MenuItem>
        </Menu>
    </Fragment>
}
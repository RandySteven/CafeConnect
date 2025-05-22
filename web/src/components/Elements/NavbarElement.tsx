import {Fragment} from "react";
import {NavbarProp} from "@/interfaces/props/NavbarProp";
import {Box, Button, Link} from "@mui/material";
import {getToken} from "@/utils/common";
import {useOnboarding} from "@/hooks/useOnboardingHook";

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
    let user = useOnboarding()
    console.log(user.name)
    return <Fragment>
        <Button sx={{
            color: '#fff'
        }}>
            {
                user.name !== "" ? (
                    <>
                        {user.username}
                    </>
                ) : (
                    <>
                        Login
                    </>
                )
            }
        </Button>
    </Fragment>
}
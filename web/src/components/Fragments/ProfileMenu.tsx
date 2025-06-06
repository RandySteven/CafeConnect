import {Fragment} from "react";
import {useProfileMenu} from "@/hooks/useOnboardingHook";
import {Button} from "@mui/material";

export const ProfileMenu = () => {
    const profileMenus = useProfileMenu()

    return <Fragment>
        {profileMenus.map((menu, index) => (
            <Button>
                {menu.menu}
            </Button>
        ))}
    </Fragment>
}
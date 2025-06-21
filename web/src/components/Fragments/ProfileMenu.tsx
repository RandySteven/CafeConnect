import {Fragment} from "react";
import {useProfileMenu} from "@/hooks/useOnboardingHook";
import {Box, Button, Link} from "@mui/material";

export const ProfileMenu = () => {
    const profileMenus = useProfileMenu()

    return (
        <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            gap={2} // spacing between buttons
        >
            {profileMenus.map((menu, index) => (
                <Button
                    key={index}
                    variant="contained"
                    href={`/${menu.menu}`}
                    fullWidth
                    sx={{ maxWidth: 300, backgroundColor: `#C38844` }} // Optional: limit button width
                >
                    {menu.menu}
                </Button>
            ))}
        </Box>
    );
}
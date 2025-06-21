import { Box } from "@mui/material";
import { ProfileImage, ProfileName } from "@/components/Elements/Profile";
import { useOnboarding } from "@/hooks/useOnboardingHook";

export const OnboardingProfile = () => {
    const onboardUser = useOnboarding();

    return (
        <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            gap={2}
        >
            <ProfileImage
                imageURL={onboardUser.profile_picture}
                name={onboardUser.name}
                sx={{
                    borderRadius: "100%",
                }}
            />
            <ProfileName
                name={onboardUser.name}
                username={onboardUser.username}
            />
        </Box>
    );
};

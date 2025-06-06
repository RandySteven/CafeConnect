import {Fragment} from "react";
import {ProfileImage, ProfileName} from "@/components/Elements/Profile";
import {useOnboarding} from "@/hooks/useOnboardingHook";

export const OnboardingProfile = () => {
    const onboardUser = useOnboarding()

    return <Fragment>
        <ProfileImage
            imageURL={onboardUser.profile_picture}
            name={onboardUser.name}
        />
        <ProfileName
            name={onboardUser.name}
            username={onboardUser.username}
        />
    </Fragment>
}
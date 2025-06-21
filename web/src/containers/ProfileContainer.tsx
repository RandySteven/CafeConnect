"use client";

import {Fragment} from "react";
import {OnboardingProfile} from "@/components/Fragments/OnboardingProfile";
import {ProfileMenu} from "@/components/Fragments/ProfileMenu";
import {AWSMap} from "@/components/Fragments/AWSMap";
import {Container} from "@/components/Fragments/Container";

export const ProfileContainer = () => {
    return <Fragment>
        <OnboardingProfile />
        <ProfileMenu />
        <AWSMap />
    </Fragment>
}
import {OnboardUserResponse} from "@/api/responses/OnboardingResponse";
import {useEffect, useState} from "react";
import {GET} from "@/api/api";
import {ONBOARD_ENDPOINT} from "@/api/endpoint";

export const useOnboarding = () : OnboardUserResponse => {
    const [onboardingUserResponse, setOnboardingUserResponse] = useState<OnboardUserResponse>({
        id: 0,
        name: "",
        username: "",
        email: "",
        profile_picture: "",
        point: 0,
        addresses: [],
        created_at: "",
        updated_at: "",
        deleted_at: ""
    })

    useEffect(() => {
        let result = GET(ONBOARD_ENDPOINT, true)
            .then((data) => setOnboardingUserResponse(data))
            .catch(error => {
                return error
            })
    }, []);

    return onboardingUserResponse
}
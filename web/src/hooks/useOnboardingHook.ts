import {LoginResponse, OnboardUserResponse} from "@/api/responses/OnboardingResponse";
import {useEffect, useState} from "react";
import {GET, POST} from "@/api/api";
import {LOGIN_ENDPOINT, ONBOARD_ENDPOINT} from "@/api/endpoint";
import {LoginRequest} from "@/api/requests/OnboardingRequest";
import {setToken} from "@/utils/common";

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
        let fetchResult = async () => {
            let result = await GET(ONBOARD_ENDPOINT, true)
            setOnboardingUserResponse(result.data.result)
        }

        fetchResult()


    }, []);

    return onboardingUserResponse
}


export const useLogin = async (request: LoginRequest): Promise<LoginResponse> => {
    const result = await POST(LOGIN_ENDPOINT, false, request);
    setToken(result.data.token.access_token);
    return result.data.token;
};
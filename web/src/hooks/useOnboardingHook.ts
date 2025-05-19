import {LoginResponse, OnboardUserResponse} from "@/api/responses/OnboardingResponse";
import {useEffect, useState} from "react";
import {GET, POST} from "@/api/api";
import {LOGIN_ENDPOINT, ONBOARD_ENDPOINT} from "@/api/endpoint";
import {LoginRequest} from "@/api/requests/OnboardingRequest";

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

export const useLogin = (request : LoginRequest) : LoginResponse => {
    const [loginResponse, setLoginResponse] = useState<LoginResponse>({
        access_token: "",
        refresh_token: "",
        login_time: ""
    })

    useEffect(() => {
        const fetchLogin = async () => {
            const result = await POST(LOGIN_ENDPOINT, false, request)
            setLoginResponse(result.data.token)
        }

        fetchLogin()
    }, [])

    return loginResponse
}
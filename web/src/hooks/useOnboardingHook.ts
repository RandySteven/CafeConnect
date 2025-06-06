import {LoginResponse, OnboardAddressUser, OnboardUserResponse} from "@/api/responses/OnboardingResponse";
import {useEffect, useState} from "react";
import {GET, POST} from "@/api/api";
import {LOGIN_ENDPOINT, ONBOARD_ENDPOINT} from "@/api/endpoint";
import {LoginRequest} from "@/api/requests/OnboardingRequest";
import {setToken} from "@/utils/common";
import {ProfileMenuContent} from "@/interfaces/contents/ProfileMenuContent";

const setAddress = (addresses : OnboardAddressUser[]) => {
    addresses.forEach((address) => {
        if(address.is_default == true) {
            localStorage.setItem(`lat`, address.latitude)
            localStorage.setItem(`long`, address.longitude)
        }
    })
}

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
            setAddress(result.data.result.addresses)
        }
        fetchResult()
    }, []);

    return onboardingUserResponse
}


export const useLogin = async (request: LoginRequest): Promise<LoginResponse> => {
    const result = await POST(LOGIN_ENDPOINT, false, request);
    setToken(result.data.token.access_token);
    let onboardUser = await GET(ONBOARD_ENDPOINT, true)
    setAddress(onboardUser.data.result.addresses)
    return result.data.token;
};

export const useProfileMenu = () : ProfileMenuContent[] => {
    const [profileMenuContent, setProfileMenuContent] = useState<ProfileMenuContent[]>([])

    useEffect(() => {
        fetch(`/contents/json/profileMenu.json`)
            .then((response) => {return response.json()})
            .then((data) => setProfileMenuContent(data))
            .catch((error) => {
                console.log(error)
                return error
            })
    }, []);

    return profileMenuContent
}
// import {cookies} from "next/headers";

export const onSubmit = () => {

}

export const setToken = (token: string) => {
    localStorage.setItem("access_token", token);
};

export const getToken = (): string | null => {
    return localStorage.getItem("access_token");
};

export const clearToken = () => {
    localStorage.removeItem("access_token");
};

// export const setTokenCookie = async (token: string) => {
//     const cookieStore = await cookies()
//     cookieStore.set({
//         name: 'token',
//         value: token,
//         httpOnly: true
//     })
// }
//
// export const getTokenCookie = async () => {
//     const cookieStore = await cookies()
//     return cookieStore.get('token')
// }
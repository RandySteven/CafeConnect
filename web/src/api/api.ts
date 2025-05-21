"use client";
import { redirect } from "next/navigation";
import {getToken} from "@/utils/common";

export const GET = async (url:string, requiredJwt: boolean) => {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
    };

    if (requiredJwt) {
        // const token = document.cookie.split('; ').find(row => row.startsWith('token='))?.split('=')[1];
        // const token = await getTokenCookie()
        const token = getToken()
        if (token) {
            headers['Authorization'] = `Bearer ${token}`
        }
    }

    try {
        const start = performance.now()
        const response = await fetch((`http://localhost:8081/${url}`), {
            method: `GET`,
            headers: headers
        })

        const end = performance.now()
        const responseTime = end - start;

        if(response.status === 401) {
            redirect('/login')
        }
        if(response.status === 500) {
            console.error("Failed response with status:", response.status);
        }
        console.log(`Response time: ${responseTime} milliseconds`);
        return response.json()
    }catch(error) {
        throw error
    }
}

export const POST = async (url: string, requiredJwt: boolean, body: any) => {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
    };

    if (requiredJwt) {
        // const token = document.cookie.split('; ').find(row => row.startsWith('token='))?.split('=')[1];
        // const token = await getTokenCookie()
        // if (token) headers['Authorization'] = `Bearer ${token}`;
    }

    try {

        console.log(`masuk post harusnya`)
        const start = performance.now()
        const response = await fetch(`http://localhost:8081/${url}`, {
            method: `POST`,
            headers,
            body: JSON.stringify(body),
        });
        const end = performance.now()
        const responseTime = end - start

        console.log(response)

        if(response.status === 401) {
            redirect('/login')
        }
        if(response.status === 500) {
            console.error("Failed response with status:", response.status);
        }
        console.log(`${url} Response time: ${responseTime} milliseconds`);
        return response.json();
    } catch (error) {
        console.error("Error in POST:", error);
        throw error;
    }
}
"use client";

import {Navbar} from "@/interfaces/contents/NavbarContent";
import {useEffect, useState} from "react";
import {error} from "next/dist/build/output/log";

export const useNavbarContent = () : Navbar => {
    const [navbar, setNavbar] = useState<Navbar>({
        contents: []
    })

    useEffect(() => {
        fetch(`/contents/json/navbarContent.json`)
            .then((response) => {return response.json()})
            .then((data) => setNavbar(data))
            .catch((error) => {
                console.log(error)
                return error
            })
    }, []);

    return navbar
}
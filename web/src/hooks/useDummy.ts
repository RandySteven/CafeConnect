import {useEffect} from "react";
import {GET} from "@/api/api";

export const useDummy = () => {
    useEffect(() => {
        const fetchData = GET("dev/check-health", false)
            .then((data) => {
                console.log(data.message)
            })
            .catch((err) => console.log(`error karna `, err))

    }, []);
}
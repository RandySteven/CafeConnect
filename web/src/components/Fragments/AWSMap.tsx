import {Fragment, useEffect, useState} from "react";
import {useAWSLocationService, useAWSLocationServiceSinglePinPoint} from "@/hooks/useGeoLoc";

export const AWSMap = () => {
    const [longValue, setLongValue] = useState<number>(0)
    const [latValue, setLatValue] = useState<number>(0)

    useEffect(() => {
        if (typeof window !== "undefined") {
            const longValue = parseFloat(localStorage.getItem("long") ?? "0");
            const latValue = parseFloat(localStorage.getItem("lat") ?? "0");

            setLongValue(longValue)
            setLatValue(latValue)
        }
    }, []);

    const href = useAWSLocationServiceSinglePinPoint(
        {
            lat: latValue,
            long: longValue
        }
    )
    return <Fragment>
        <div ref={href} style={{ width: "100%", height: "500px" }}  />
    </Fragment>
}
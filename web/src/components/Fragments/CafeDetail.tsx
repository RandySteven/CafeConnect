"use client";

import {Fragment} from "react";
import {IDParam} from "@/interfaces/props/ParamProp";
import {useCafeDetail} from "@/hooks/useCafeHook";
import {CafeTitle} from "@/components/Elements/Title";

export const CafeDetail = (prop : IDParam) => {
    const cafe = useCafeDetail(prop.id)

    return <Fragment>
        <CafeTitle
            name={cafe.name}
            img={cafe.logo_url}
            address={cafe.address?.address}
        />
    </Fragment>
}
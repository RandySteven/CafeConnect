"use client";

import {Fragment, useEffect, useState} from "react";
import {IDParam} from "@/interfaces/props/ParamProp";
import {useCafeDetail} from "@/hooks/useCafeHook";
import {CafeMap, CafeTitle} from "@/components/Elements/Title";
import {Coordinate} from "@/hooks/useGeoLoc";
import {CafeListMaps} from "@/components/Fragments/ListCafe";
import {sleep} from "maplibre-gl/src/util/test/util";

export const CafeDetail = (prop : IDParam) => {
    const cafe = useCafeDetail(prop.id)
    console.log(cafe)
    return <Fragment>
        <CafeTitle
            name={cafe.name}
            img={cafe.logo_url}
            address={cafe.address?.address}
            longitude={Number(cafe.address?.longitude)}
            latitude={Number(cafe.address?.latitude)}
        />
        {/*<CafeMap longitude={cafe.address?.longitude} latitude={cafe.address?.latitude} />*/}
    </Fragment>
}
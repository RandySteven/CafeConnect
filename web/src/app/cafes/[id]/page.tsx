"use client";

import {Fragment, use} from "react";
import {HomeContainer} from "@/containers/HomeContainer";
import {ParamProp} from "@/interfaces/props/ParamProp";
import {CafeDetailContainer} from "@/containers/CafeDetailContainer";

export default function Home({prop} : ParamProp) {
    const resolvedParams = use(prop)
    return (
        <Fragment>
            <CafeDetailContainer id={resolvedParams.id} />
        </Fragment>
    );
}

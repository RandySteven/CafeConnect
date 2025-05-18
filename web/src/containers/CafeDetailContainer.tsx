import {IDParam} from "@/interfaces/props/ParamProp";
import {Fragment} from "react";
import {CafeDetail} from "@/components/Fragments/CafeDetail";

export const CafeDetailContainer = (id : IDParam) => {
    return <Fragment>
        <CafeDetail id={id.id} />
    </Fragment>
}
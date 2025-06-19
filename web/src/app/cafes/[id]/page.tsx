import {Fragment, use} from "react";
import {ParamProp} from "@/interfaces/props/ParamProp";
import {CafeDetailContainer} from "@/containers/CafeDetailContainer";

export default function Home({params} : ParamProp) {
    const resolvedParams = use(params)
    return (
        <Fragment>
            <CafeDetailContainer id={parseInt(resolvedParams.id)} />
        </Fragment>
    );
}

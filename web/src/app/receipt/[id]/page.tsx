import {Fragment, use} from "react";
import {ReceiptContainer} from "@/containers/ReceiptContainer";
import {ParamProp} from "@/interfaces/props/ParamProp";

export default function Home({params} : ParamProp) {
    const resolvedParams = use(params)

    return <Fragment>
        <ReceiptContainer code={resolvedParams.id}/>
    </Fragment>
}
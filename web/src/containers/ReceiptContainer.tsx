import {Fragment} from "react";
import {Receipt} from "@/components/Fragments/Receipt";
import {CodeParam} from "@/interfaces/props/ParamProp";

export const ReceiptContainer = (code : CodeParam) => {
    return <Fragment>
        <Receipt code={code.code}/>
    </Fragment>
}
import {Fragment} from "react";
import {CodeParam} from "@/interfaces/props/ParamProp";
import {useHistoryDetail} from "@/hooks/useTransaction";

export const HistoryDetailContainer = (code : CodeParam) => {
    const transactionDetail = useHistoryDetail(code.code)
    return <Fragment>

    </Fragment>
}
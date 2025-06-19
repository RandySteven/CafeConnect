import {Fragment} from "react";
import {getTotalAmounts} from "@/utils/common";
import {Typography} from "@mui/material";

export const TotalAmount = (props : {
    amounts : number[]
}) => {
    const totalAmount = getTotalAmounts(props.amounts)
    return <Fragment>
        <Typography>
            {totalAmount}
        </Typography>
    </Fragment>
}
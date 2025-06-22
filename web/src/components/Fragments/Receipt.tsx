"use client";

import {Fragment, useState} from "react";
import {useReceipt} from "@/hooks/useTransaction";
import {CodeParam} from "@/interfaces/props/ParamProp";
import {PaymentButton} from "@/components/Elements/Button";
import {TransactionReceiptResponse} from "@/api/responses/TransactionResponse";
import {ReceiptInfo} from "@/components/Elements/ReceiptInfo";

export const Receipt = (code : CodeParam) => {
    // let receipt : TransactionReceiptResponse
    // const [loading, setLoading] = useState<boolean>(true)
    console.log(code)
    // setInterval(() => {
        const receipt = useReceipt({
            transaction_code: code.code
        })

    //     if(receipt.status !== "PENDING") {
    //         setLoading(false)
    //     }
    //
    // }, 60000)
    let paymentButton = <div>Loading</div>
    if(receipt?.midtrans_response != null) {
        paymentButton = <PaymentButton midtransUrl={receipt?.midtrans_response.redirect_url} snapToken={receipt?.midtrans_response?.token} transactionCode={receipt.transaction_code}/>
    }


    return <Fragment>
        <ReceiptInfo transactionCode={receipt?.transaction_code} transactionAt={receipt?.transaction_at} status={receipt?.status} />
        <PaymentButton midtransUrl={receipt?.midtrans_response?.redirect_url} snapToken={receipt?.midtrans_response?.token} transactionCode={receipt?.transaction_code}/>
    </Fragment>
}
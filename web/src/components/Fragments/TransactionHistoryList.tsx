"use client";

import {Fragment} from "react";
import {useTransactionList} from "@/hooks/useTransaction";
import {ListCard} from "@/components/Elements/Card";

export const TransactionHistoryList = () => {
    const transactionList = useTransactionList()
    return <Fragment>
        {transactionList.map((transaction, index) => (
            <div>
                {transaction.transaction_code}
            </div>
        ))}
    </Fragment>
}
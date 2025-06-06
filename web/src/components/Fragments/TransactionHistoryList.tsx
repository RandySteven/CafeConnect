"use client";

import {Fragment} from "react";
import {useTransactionList} from "@/hooks/useTransaction";
import {ListCard} from "@/components/Elements/Card";

export const TransactionHistoryList = () => {
    const transactionList = useTransactionList()
    return <Fragment>
        {transactionList.map((transaction, index) => (
            <ListCard
                type={`transaction`}
                img={transaction.cafe.image_url}
                name={transaction.cafe.name}
                address={transaction.cafe.address}
                status={transaction.status}
            >
                {transaction.cafe.name} + {transaction.cafe.address}
            </ListCard>
        ))}
    </Fragment>
}
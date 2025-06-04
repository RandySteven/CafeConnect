import {TransactionListResponse} from "@/api/responses/TransactionResponse";
import {useEffect, useState} from "react";
import {GET} from "@/api/api";
import {TRANSACTION_LIST} from "@/api/endpoint";

export const useTransactionList = () : TransactionListResponse[] => {
    const [transactionListResponse, setTransactionListResponse] = useState<TransactionListResponse[]>([])

    useEffect(() => {
        const fetch = async () => {
            const result = await GET(TRANSACTION_LIST, true)

            setTransactionListResponse(result.data.transactions)
        }

        fetch()
    }, []);

    return transactionListResponse
}
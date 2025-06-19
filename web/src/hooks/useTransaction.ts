import {TransactionListResponse, TransactionReceiptResponse} from "@/api/responses/TransactionResponse";
import {useEffect, useState} from "react";
import {GET, POST} from "@/api/api";
import {TRANSACTION_LIST, TRANSACTION_RECEIPT} from "@/api/endpoint";
import {ReceiptRequest} from "@/api/requests/TransactionRequest";

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

export const useReceipt = (request : ReceiptRequest): TransactionReceiptResponse=> {
    const [transactionReceipt, setTransactionReceipt] = useState<TransactionReceiptResponse>()

    useEffect(() => {
        const fetch = async () => {
            const result = await POST(TRANSACTION_RECEIPT, true, request)

            setTransactionReceipt(result.data.result)
        }

        fetch()
    }, []);

    return transactionReceipt
}
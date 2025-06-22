import {Fragment, useEffect, useState} from "react";
import {Button} from "@mui/material";
import {PaymentButtonProp} from "@/interfaces/props/ButtonProp";
import {useMidtransSnap} from "@/hooks/useTransaction";
import {POST} from "@/api/api";
import {PAYMENT_CONFIRMATION} from "@/api/endpoint";


export const PaymentButton = (props : PaymentButtonProp) => {
    const [disable, setDisable] = useState<boolean>(true);

    useMidtransSnap()

    useEffect(() => {
        if (props) {
            setDisable(false);
        }
    }, [props.snapToken]);

    const handleSnapPay = () => {
        if (!window.snap || !props.snapToken) return;

        window.snap.pay(props.snapToken, {
            onSuccess: async (result) => {
                console.log("✅ Payment Success:", result);
                try {
                    const response = await POST(PAYMENT_CONFIRMATION, true, {
                        transaction_code: props.transactionCode,
                        status: "SUCCESS",
                    });

                    console.log("✅ Backend confirmation sent:", response);
                } catch (error) {
                    console.error("❌ Failed to confirm transaction to backend:", error);
                }
            },

            onPending: (result) => {
                console.log("⏳ Payment Pending:", result);
            },
            onError: (result) => {
                console.error("❌ Payment Error:", result);
            },
            onClose: () => {
                console.log("Modal closed by user.");
            },
        });
    };


    return <Fragment>
        <Button onClick={handleSnapPay}
                style={{
                    backgroundColor: 'blue',
                    color: 'white',
                    fontStyle: 'bold'
                }}
                disableElevation={disable}
        >
            Pay
        </Button>
    </Fragment>
}
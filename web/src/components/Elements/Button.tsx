import {Fragment, useEffect, useState} from "react";
import {Button} from "@mui/material";
import {PaymentButtonProp} from "@/interfaces/props/ButtonProp";


export const PaymentButton = (props : PaymentButtonProp) => {
    const [disable, setDisable] = useState<boolean>(true)

    useEffect(() => {
        if(props.midtransUrl != null) {
            setDisable(true)
        }
    }, []);

    return <Fragment>
        <Button onClick={() => {
            location.href = props?.midtransUrl
        }}
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
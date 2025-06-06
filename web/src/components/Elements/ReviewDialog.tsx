import React, { Fragment, useState } from "react";
import { TransitionProps } from '@mui/material/transitions';
import Dialog from '@mui/material/Dialog';
import Slide from '@mui/material/Slide';
import {
    Button,
    DialogActions,
    DialogContent,
    DialogTitle,
    TextField,
    Box,
    Stack,
} from "@mui/material";
import {POST} from "@/api/api";
import {GET_CARTS} from "@/api/endpoint";
import {redirect} from "next/navigation";

const Transition = React.forwardRef(function Transition(
    props: TransitionProps & {
        children: React.ReactElement<any, any>;
    },
    ref: React.Ref<unknown>,
) {
    return <Slide direction="up" ref={ref} {...props} />;
});

export const ReviewDialog = (
    props: {
        open: boolean;
        handleClose: () => void;
        product: {
            id: number
            name: string;
            imageURL: string;
        };
    }
) => {
    console.log(`item id `, props.product.id)
    const [qty, setQty] = useState<number>(0);

    const handleIncrease = () => {
        setQty(qty + 1);
    };

    const handleDecrease = () => {
        setQty(prev => (prev <= 0 ? 0 : prev - 1));
    };

    const handleConfirm = async () => {
        try {
            let request = {
                cafe_product_id: props.product.id,
                qty: qty
            }
            const response = await POST(GET_CARTS, true, request)

            props.handleClose();
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <Fragment>
            <Dialog
                open={props.open}
                TransitionComponent={Transition}
                keepMounted
                onClose={props.handleClose}
                aria-describedby="alert-dialog-slide-description"
            >
                <DialogTitle>{props.product.name}</DialogTitle>
                <DialogContent>
                    <Box display="flex" flexDirection="column" alignItems="center">
                        <img
                            src={props.product.imageURL}
                            alt={props.product.name}
                            style={{
                                height: 240,
                                width: 240,
                                objectFit: "cover",
                                marginBottom: 16,
                                borderRadius: 8,
                            }}
                        />
                        <Stack direction="row" spacing={2} alignItems="center">
                            <Button variant="outlined" onClick={handleDecrease}>-</Button>
                            <TextField
                                type="number"
                                value={qty}
                                inputProps={{ readOnly: true, style: { textAlign: 'center' } }}
                                size="small"
                                style={{ width: 60 }}
                            />
                            <Button variant="outlined" onClick={handleIncrease}>+</Button>
                        </Stack>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={props.handleClose} color="primary">Close</Button>
                    <Button onClick={handleConfirm} color="primary" variant="contained">Confirm</Button>
                </DialogActions>
            </Dialog>
        </Fragment>
    );
};

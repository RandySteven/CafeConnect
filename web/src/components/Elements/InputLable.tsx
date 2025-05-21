import React, { Fragment } from "react";
import { TextField } from "@mui/material";

interface InputLabelProp {
    id: string;
    label: string;
    inputType: string;
    placeholder?: string;
    value: any;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const InputLabel = (prop: InputLabelProp) => {
    return (
        <Fragment>
            <TextField
                id={prop.id}
                label={prop.label}
                value={prop.value}
                placeholder={prop.placeholder}
                onChange={prop.onChange}
                type={prop.inputType}
                fullWidth
                margin="normal"
            />
        </Fragment>
    );
};

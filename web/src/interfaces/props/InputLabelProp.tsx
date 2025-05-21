import React from "react";

export interface InputLabelProp {
    id: string
    label: string
    inputType: string
    placeholder: string
    value: any
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => {}
}
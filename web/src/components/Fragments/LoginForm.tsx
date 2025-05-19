import {Fragment} from "react";
import {FormContainer} from "@/components/Elements/FormContainer";
import {FormHelperText, Input, InputLabel} from "@mui/material";

export const LoginForm = () => {
    return <Fragment>
        <FormContainer title={`Login`}>
            <InputLabel htmlFor="my-input">Email address</InputLabel>
            <Input id="my-input" aria-describedby="my-helper-text" />
            <FormHelperText id="my-helper-text">We'll never share your email.</FormHelperText>

            <InputLabel htmlFor="my-input-2">Password</InputLabel>
            <Input id="my-input-2" aria-describedby="my-helper-text-2" />
            <FormHelperText id="my-helper-text-2">We'll never share your password.</FormHelperText>

        </FormContainer>
    </Fragment>
}
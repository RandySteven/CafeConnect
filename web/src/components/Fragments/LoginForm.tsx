"use client";

import { useRouter } from "next/navigation";
import {FormContainer} from "@/components/Elements/FormContainer";
import {InputLabel} from "@/components/Elements/InputLable";
import React, {Fragment, useState} from "react";
import {LoginRequest} from "@/api/requests/OnboardingRequest";
import {Box, Button} from "@mui/material";
import {useLogin} from "@/hooks/useOnboardingHook";

export const LoginForm = () => {
    const router = useRouter()
    const [loginRequest, setLoginRequest] = useState<LoginRequest>({
        email: "",
        password: ""
    });

    const handleChange = (field: keyof LoginRequest) =>
        (event: React.ChangeEvent<HTMLInputElement>) => {
            setLoginRequest({
                ...loginRequest,
                [field]: event.target.value
            });
        };

    const handleSubmit = async () => {
        try {
            const result = await useLogin(loginRequest)
            router.push(`/`)
        } catch (error) {
            console.error("Login failed", error);
        }
    };

    return (
        <Fragment>
            <FormContainer title="Login">
                <Box display="flex" flexDirection="column" gap={2} mt={2}>
                    <InputLabel
                        id="email"
                        label="Email"
                        inputType="email"
                        placeholder="Enter your email"
                        value={loginRequest.email}
                        onChange={handleChange("email")}
                    />
                    <InputLabel
                        id="password"
                        label="Password"
                        inputType="password"
                        placeholder="Enter your password"
                        value={loginRequest.password}
                        onChange={handleChange("password")}
                    />
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={handleSubmit}
                    >
                        Login
                    </Button>
                </Box>
            </FormContainer>
        </Fragment>
    );
};

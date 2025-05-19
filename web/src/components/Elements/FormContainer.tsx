import { ContainerProps, FormContainerProps } from "@/interfaces/props/ContainerProp";
import { Fragment } from "react";
import { Box, FormControl, Typography } from "@mui/material";

export const FormContainer = ({ title, children }: FormContainerProps) => {
    return (
        <Fragment>
            <FormControl fullWidth>
                <Box
                    display="flex"
                    flexDirection="column"
                    justifyContent="center"
                    alignItems="center"
                    textAlign="center"
                    sx={{ minHeight: '100vh', gap: 2 }} // Full height and spacing
                >
                    <Box display="flex" alignItems="center" gap={1}>
                        <img
                            src="/assets/img/cafeConnect-icon.png"
                            alt="CafeConnect icon"
                            style={{ height: 50 }}
                        />
                        <Typography variant="h4" fontWeight="bold">
                            {title}
                        </Typography>
                    </Box>
                    <Box width="100%" maxWidth={400}>
                        {children}
                    </Box>
                </Box>
            </FormControl>
        </Fragment>
    );
};

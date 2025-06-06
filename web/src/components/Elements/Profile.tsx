import React, {Fragment} from "react";
import {ProfileImageProp, ProfileNameProp} from "@/interfaces/props/ProfileProp";
import {Avatar, Box, CardContent} from "@mui/material";

export const ProfileImage = (prop : ProfileImageProp) => {
    return <Fragment>
        <CardContent sx={{ display: 'flex', gap: 2 }}>
            <Avatar src={prop.imageURL} alt={prop.name} />
        </CardContent>
    </Fragment>
}

export const ProfileName = (prop : ProfileNameProp) => {
    return <Fragment>
        <Box>
            <h1>{prop.name}</h1>
        </Box>
    </Fragment>
}
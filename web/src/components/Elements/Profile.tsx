import React, {Fragment} from "react";
import {ProfileImageProp, ProfileNameProp} from "@/interfaces/props/ProfileProp";
import {Avatar, Box, CardContent} from "@mui/material";

export const ProfileImage = (prop : ProfileImageProp) => {
    return <Fragment>
        <CardContent sx={prop.sx}>
            <Avatar src={prop.imageURL} alt={prop.name} sx={{
                width: 150, height: 150
            }}/>
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
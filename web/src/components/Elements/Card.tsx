import {Fragment} from "react";
import {Box, Card, CardMedia, CardProps} from "@mui/material";
import {CardProp} from "@/interfaces/props/CardProp";

export const ListCard = (prop : CardProp) => {
    return <Fragment>
        <Card
            sx={{
                display: 'flex'
            }}
        >
            <CardMedia
                component={"img"}
                image={prop.img}
                alt={""}
            />
            <Box
                sx={{

                }}
            >

            </Box>
        </Card>
    </Fragment>
}

export const ListCardContent = () => {
    return <Fragment>

    </Fragment>
}

export const GridCard = () => {
    return <Fragment>

    </Fragment>
}
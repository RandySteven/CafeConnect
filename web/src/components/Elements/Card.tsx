import {Fragment} from "react";
import {
    Box,
    Card, CardActionArea, CardContent,
    CardMedia,
    CardProps,
    Link,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Typography
} from "@mui/material";
import {CardProp} from "@/interfaces/props/CardProp";
import {MenuProp} from "@/interfaces/props/MenuProp";
import Image from "next/image";

interface CafeDataProp {
    name : string,
    status : string,
    openHour : string
    closeHour : string
    address : string
}

export const ListCard = (prop : CardProp) => {
    return (
        <Fragment>
            <Card sx={{ display: "flex", p: 2, my: 2}}>
                <CardMedia
                    component="img"
                    image={prop.img}
                    alt="card image"
                    sx={{ width: 120, height: 120, objectFit: "cover", borderRadius: 2 }}
                />

                <Box ml={2} display="flex" flexDirection="column" justifyContent="center">
                    {prop.type === "product" && (
                        <>
                            <Typography variant="h6">{prop.title}</Typography>
                            <Typography variant="body2" color="text.secondary">
                                {prop.description}
                            </Typography>
                        </>
                    )}

                    {prop.type === "cafe" && (
                        <CafeCard
                            name={prop.name}
                            status={prop.status}
                            openHour={prop.openHour}
                            closeHour={prop.closeHour}
                            address={prop.address}
                        />
                    )}
                </Box>
            </Card>
        </Fragment>
    );
}

const CafeCard = (prop : CafeDataProp) => {
    return <Fragment>
        <>
            <Typography variant="h6">{prop.name}</Typography>
            <Typography variant="body2">Status: {prop.status}</Typography>
            <Typography variant="body2">
                Hours: {prop.openHour} - {prop.closeHour}
            </Typography>
            <Typography variant="body2">
                Address: {prop.address}
            </Typography>

        </>
    </Fragment>
}

export const GridCard = () => {
    return <Fragment>

    </Fragment>
}

export const MenuCard = (prop : MenuProp) => {
    return <Fragment>
        <Link href={prop.link} passHref legacyBehavior>
            <Card sx={{ height: "100%" }}>
                <CardActionArea component="a">
                    <Box
                        sx={{
                            display: "flex",
                            justifyContent: "center",
                            alignItems: "center",
                            mb: 2,
                        }}>
                        <Box
                            component="img"
                            src={prop.icon}
                            alt={`${name} icon`}
                            sx={{
                                width: 120,
                                height: 120,
                                borderRadius: "50%",
                                objectFit: "cover",
                            }}
                        />
                    </Box>

                    {/*<CardContent>*/}
                    {/*    <Typography variant="h6" component="div" align="center">*/}
                    {/*        {prop.name}*/}
                    {/*    </Typography>*/}
                    {/*</CardContent>*/}
                </CardActionArea>
            </Card>
        </Link>
    </Fragment>
}
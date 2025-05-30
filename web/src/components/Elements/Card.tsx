import {Fragment, useState} from "react";
import {
    Avatar,
    Box, Button,
    Card, CardActionArea, CardActions, CardContent,
    CardMedia,
    CardProps,
    Link,
    ListItemButton,
    ListItemIcon,
    ListItemText, Rating,
    Typography
} from "@mui/material";
import {CardProp, CommentProp} from "@/interfaces/props/CardProp";
import {MenuProp} from "@/interfaces/props/MenuProp";
import Image from "next/image";
import {ReviewDialog} from "@/components/Elements/ReviewDialog";

interface CafeDataProp {
    name : string,
    status : string,
    openHour : string
    closeHour : string
    address : string
}

interface ProductDataProp {
    name: string
    stock: number
}

export const ListCard = (prop : CardProp) => {
    return (
        <Fragment>
            <Link href={prop.link}>
                <Card sx={{ display: "flex", p: 2, my: 2}}>
                    <CardMedia
                        component="img"
                        image={prop.img}
                        alt="card image"
                        sx={{ width: 120, height: 120, objectFit: "cover", borderRadius: 2 }}
                    />

                    <Box ml={2} display="flex" flexDirection="column" justifyContent="center">
                        {prop.type === "product" && (
                            <ProductCard
                                name={prop.name}
                                stock={prop.stock} />
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
            </Link>

        </Fragment>
    );
}

export const CommentCard = (prop : CommentProp) => {
    return <Fragment>
        <Card variant="outlined" sx={{ mb: 2 }}>
            <CardContent sx={{ display: 'flex', gap: 2 }}>
                <Avatar src={prop.avatar} alt={prop.name} />
                <Box>
                    <Box sx={{ backgroundColor: '#f0f2f5', p: 1.5, borderRadius: 2 }}>
                        <Typography variant="subtitle2" fontWeight="bold">
                            {prop.name}
                        </Typography>
                        <Rating
                            name="read-only"
                            value={prop.score} // must be a number between 0-5
                            readOnly
                            precision={0.5}
                            size="small"
                            sx={{ mb: 0.5 }}
                        />
                        <Typography variant="body2">
                            {prop.comment}
                        </Typography>
                    </Box>
                    {prop.timestamp && (
                        <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5 }}>
                            {prop.timestamp}
                        </Typography>
                    )}
                </Box>
            </CardContent>
        </Card>
    </Fragment>
}

const CafeCard = (prop : CafeDataProp) => {
    return <Fragment>
        <>
            <Typography variant="h6">{prop.name}</Typography>
            <Status status={prop.status} />
            <Hour openHour={prop.openHour} closeHour={prop.closeHour} />
            <Typography variant="body2">
                Address: {prop.address}
            </Typography>

        </>
    </Fragment>
}

const ProductCard = (prop : ProductDataProp) => {
    return <>
        <>
            <Typography variant="h6">{prop.name}</Typography>
            <Typography variant="body2" color="text.secondary">
                {prop.stock}
            </Typography>
        </>
    </>
}

const Status = (prop : {
    status: string
}) => {
    let color = ``
    switch (prop.status) {
        case "OPEN":
            color = `#4caf50`
            break
        case "CLOSED":
            color = `#d50000`
            break
    }

    return(
        <Typography color={color}>
            {prop.status}
        </Typography>
    )
}

const Hour = (prop : {
    openHour: string
    closeHour: string
}) => {
    let result = ``
    if (prop.openHour === "00:00:00" && prop.closeHour === "00:00:00") {
        result = "24 Hours"
    }else {
        result = `${prop.openHour} - ${prop.closeHour}`
    }

    return <Typography>
        <b>{result}</b>
    </Typography>
}

export const GridCard = (item : {
    id: string
    image: string
    name: string
    price: number
    stock: number
}) => {
    console.log(`item id `, item.id)
    const [open, setOpen] = useState<boolean>(false)

    const handleClickOpen = () => {
        setOpen(true);
    };
    const handleClose = () => {
        setOpen(false);
    };

    return <Fragment>
        <Card sx={{ height: '100%' }}>
            <CardMedia
                component="img"
                height="160"
                image={item.image}
                alt={item.name}
                sx={{ height: 180,
                    width: '100%',
                    objectFit: 'cover',
                    borderTopLeftRadius: 4,
                    borderTopRightRadius: 4, }}
            />
            <CardContent>
                <Typography gutterBottom variant="b" component="div">
                    {item.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Price: Rp {item.price.toLocaleString()}
                </Typography>
                <Typography variant="body2" color={item.stock > 0 ? "text.primary" : "error"}>
                    Stock: {item.stock > 0 ? item.stock : "Out of stock"}
                </Typography>
            </CardContent>
            <CardActions>
                <Button
                    size="small"
                    variant="contained"
                    color="primary"
                    disabled={item.stock === 0}
                    onClick={handleClickOpen}
                >
                    See Detail
                </Button>
                <ReviewDialog open={open} handleClose={handleClose}  product={
                    {
                        id: Number(item.id),
                        imageURL: item.image,
                        name: item.name
                    }
                }/>
            </CardActions>
        </Card>
    </Fragment>
}

export const MenuCard = (prop : MenuProp) => {
    return <Fragment>
        <Link href={prop.link}>
            <Card sx={{ height: "100%" }}>
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
            </Card>
        </Link>
    </Fragment>
}
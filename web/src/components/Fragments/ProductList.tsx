"use client";

import {Fragment} from "react";
import {useProductByCafe} from "@/hooks/useProduct";
import {Grid} from "@mui/system";
import {GridCard, ListCard} from "@/components/Elements/Card";
import {useListCafeWithRadius} from "@/hooks/useCafeHook";
import {Box} from "@mui/material";
import {wait} from "next/dist/lib/wait";
import {getItem} from "@/utils/common";

export const ProductList = (prop : {
    id : number
}) => {
    const products = useProductByCafe([Number(prop.id)])
    return <Fragment>
        <Grid container spacing={3}>
            {products.map((product, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                    <GridCard
                            id={product.id}
                            image={product.photo}
                              name={product.name} price={product.price} stock={product.stock} />
                </Grid>
            ))}
        </Grid>
    </Fragment>
}

export const MenuList = () => {
    let cafeIds : number[] = []
    let cafeIdStorages = getItem("cafeIds")?.split(",")
    for(let i = 0 ; i < cafeIdStorages?.length; i++) {
        cafeIds[i] = Number(cafeIdStorages[i])
    }
    const products = useProductByCafe(cafeIds)
    return <Fragment>
        <Box
            sx={{
                py: 2
            }}
        >
            {
                products.map((product, index) => (
                    <ListCard
                        link={`/products/${product.id}`}
                        key={index}
                        type="product"
                        img={product.photo}
                        name={product.name}
                    />
                ))
            }
        </Box>
    </Fragment>
}
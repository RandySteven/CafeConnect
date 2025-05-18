"use client";

import {Fragment} from "react";
import {useProductByCafe} from "@/hooks/useProduct";
import {Grid} from "@mui/system";
import {GridCard} from "@/components/Elements/Card";

export const ProductList = (prop : {
    id : number
}) => {
    const products = useProductByCafe(prop.id)

    return <Fragment>
        <Grid container spacing={3}>
            {products.map((product, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                    <GridCard image={product.photo}
                              name={product.name} price={product.price} stock={product.stock} />
                </Grid>
            ))}
        </Grid>
    </Fragment>
}
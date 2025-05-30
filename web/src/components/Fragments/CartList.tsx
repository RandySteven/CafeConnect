"use client";

import {Fragment} from "react";
import {useCart} from "@/hooks/useCart";
import {ListCard} from "@/components/Elements/Card";
import {Box, Typography} from "@mui/material";

export const CartList = () => {
    const carts = useCart()

    return <Fragment>
        <Box sx={{ py: 2 }}>
            {carts.checkout_list.map((checkout, index) => (
                <Box key={index} sx={{ mb: 3 }}>
                    <Typography variant="h6" sx={{ mb: 1 }}>
                        {checkout.cafe_name}
                    </Typography>
                    {checkout.items.map((product, pIndex) => (
                        <ListCard
                            key={pIndex}
                            link={`/products/${product.product_id}`}
                            type="product"
                            name={product.product_name}
                            img={product.product_image}
                            stock={product.qty}
                        />
                    ))}
                </Box>
            ))}
        </Box>
    </Fragment>
}
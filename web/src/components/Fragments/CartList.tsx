import {
    Box,
    Typography,
    Checkbox,
    FormControlLabel,
    Button,
} from "@mui/material";
import { Fragment, useState } from "react";
import { useCart } from "@/hooks/useCart";
import { ListCard } from "@/components/Elements/Card";
import {POST} from "@/api/api";
import {TRANSACTION_V2_CHECKOUT} from "@/api/endpoint";
import {redirect} from "next/navigation";

export const CartList = () => {
    const carts = useCart();
    const [selectedProducts, setSelectedProducts] = useState<Set<number>>(new Set());

    const toggleProduct = (productId: number) => {
        setSelectedProducts((prev) => {
            const newSet = new Set(prev);
            newSet.has(productId) ? newSet.delete(productId) : newSet.add(productId);
            return newSet;
        });
    };

    const toggleCafe = (cafeItemIds: number[]) => {
        const allSelected = cafeItemIds.every((id) => selectedProducts.has(id));
        setSelectedProducts((prev) => {
            const newSet = new Set(prev);
            cafeItemIds.forEach((id) =>
                allSelected ? newSet.delete(id) : newSet.add(id)
            );
            return newSet;
        });
    };

    const handleCheckout = async () => {
        const groupedByCafe: Record<number, { cafe_product_id: number; qty: number }[]> = {};

        carts.checkout_list.forEach((checkout) => {
            const selectedItems = checkout.items.filter((item) =>
                selectedProducts.has(item.product_id)
            );

            if (selectedItems.length > 0) {
                groupedByCafe[checkout.cafe_id] = selectedItems.map((item) => ({
                    cafe_product_id: item.product_id,
                    qty: item.qty, // using existing qty in cart
                }));
            }
        });

        for (const [cafe_id_str, checkouts] of Object.entries(groupedByCafe)) {
            const cafe_id = parseInt(cafe_id_str, 10);

            const body = {
                cafe_id,
                checkouts,
            };

            try {
                const res = await POST(TRANSACTION_V2_CHECKOUT, true, body)
                redirect(``)
            } catch (err) {
                console.error("Checkout error:", err);
            }
        }
    };

    return (
        <Fragment>
            <Box sx={{ py: 2 }}>
                {carts.checkout_list.map((checkout, index) => {
                    const cafeProductIds = checkout.items.map((item) => item.product_id);
                    const allProductsSelected = cafeProductIds.every((id) => selectedProducts.has(id));

                    return (
                        <Box key={index} sx={{ mb: 3 }}>
                            <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={allProductsSelected}
                                            onChange={() => toggleCafe(cafeProductIds)}
                                        />
                                    }
                                    label={<Typography variant="h6">{checkout.cafe_name}</Typography>}
                                />
                            </Box>

                            {checkout.items.map((product, pIndex) => (
                                <Box key={pIndex} sx={{ display: "flex", alignItems: "center" }}>
                                    <Checkbox
                                        checked={selectedProducts.has(product.product_id)}
                                        onChange={() => toggleProduct(product.product_id)}
                                    />
                                    <ListCard
                                        link={`/products/${product.product_id}`}
                                        type="cart"
                                        name={product.product_name}
                                        img={product.product_image}
                                        stock={product.qty}
                                    />
                                </Box>
                            ))}
                        </Box>
                    );
                })}

                <Box sx={{ mt: 4 }}>
                    <Button
                        variant="contained"
                        color="primary"
                        disabled={selectedProducts.size === 0}
                        onClick={handleCheckout}
                    >
                        Checkout Selected
                    </Button>
                </Box>
            </Box>
        </Fragment>
    );
};

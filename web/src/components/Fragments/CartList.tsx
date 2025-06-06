import {
    Box,
    Typography,
    Checkbox,
    FormControlLabel,
    Button,
} from "@mui/material";
import {Fragment, useEffect, useMemo, useState} from "react";
import { useCart } from "@/hooks/useCart";
import { ListCard } from "@/components/Elements/Card";
import {POST} from "@/api/api";
import {TRANSACTION_V2_CHECKOUT} from "@/api/endpoint";
import {redirect} from "next/navigation";

export const CartList = () => {
    const carts = useCart();
    const [selectedProducts, setSelectedProducts] = useState<Set<number>>(new Set());
    const [productQuantities, setProductQuantities] = useState<Record<number, number>>({});

    const toggleProduct = (productId: number) => {
        setSelectedProducts((prev) => {
            const newSet = new Set(prev);
            newSet.has(productId) ? newSet.delete(productId) : newSet.add(productId);
            return newSet;
        });
    };

    const initialQuantities = useMemo(() => {
        const quantities: Record<number, number> = {};
        carts.checkout_list.forEach((checkout) => {
            checkout.items.forEach((item) => {
                quantities[item.product_id] = item.qty;
            });
        });
        return quantities;
    }, [carts.checkout_list]);

    useEffect(() => {
        setProductQuantities(initialQuantities);
    }, [initialQuantities]);


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

    const increaseQty = (productId: number) => {
        setProductQuantities((prev) => ({
            ...prev,
            [productId]: (prev[productId] || 0) + 1,
        }));
    };

    const decreaseQty = (productId: number) => {
        setProductQuantities((prev) => ({
            ...prev,
            [productId]: Math.max(0, (prev[productId] || 0) - 1),
        }));
    };

    useEffect(() => {
        const initialQuantities: Record<number, number> = {};
        carts.checkout_list.forEach((checkout) => {
            checkout.items.forEach((item) => {
                initialQuantities[item.product_id] = item.qty;
            });
        });
        setProductQuantities(initialQuantities);
    }, [carts.checkout_list]);


    const handleCheckout = async () => {
        const groupedByCafe: Record<number, { cafe_product_id: number; qty: number }[]> = {};

        carts.checkout_list.forEach((checkout) => {
            const selectedItems = checkout.items.filter((item) =>
                selectedProducts.has(item.product_id)
            );

            if (selectedItems.length > 0) {
                groupedByCafe[checkout.cafe_id] = selectedItems.map((item) => ({
                    cafe_product_id: item.product_id,
                    qty: productQuantities[item.product_id] ?? item.qty,
                }));
            }
        });

        for (const [cafe_id_str, checkouts] of Object.entries(groupedByCafe)) {
            const cafe_id = parseInt(cafe_id_str, 10);
            const body = { cafe_id, checkouts };

            try {
                const res = await POST(TRANSACTION_V2_CHECKOUT, true, body);
                redirect(`histories`);
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
                                        stock={productQuantities[product.product_id] || product.qty}
                                        onIncrease={() => increaseQty(product.product_id)}
                                        onDecrease={() => decreaseQty(product.product_id)}
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

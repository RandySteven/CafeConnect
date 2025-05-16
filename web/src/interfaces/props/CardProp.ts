export type BaseCardProp = {
    img: string;
};

export type ProductCardProp = BaseCardProp & {
    type: 'product';
    title: string;
    description: string;
};

export type CafeCardProp = BaseCardProp & {
    type: 'cafe';
    name: string;
    status: string;
    openHour: string;
    closeHour: string;
    address: string;
};

export type CardProp = ProductCardProp | CafeCardProp;
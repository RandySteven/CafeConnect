export type BaseCardProp = {
    link: string;
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

export interface CommentProp {
    avatar: string
    name: string
    score: number
    comment: string
    timestamp: string
}
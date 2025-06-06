export type BaseCardProp = {
    img: string;
};

export type CartCardProp = BaseCardProp & {
    type: 'cart';
    name: string;
    stock: number;
    onIncrease: () => void;
    onDecrease: () => void;
};

export type CafeCardProp = BaseCardProp & {
    link: string;
    type: 'cafe';
    name: string;
    status: string;
    openHour: string;
    closeHour: string;
    address: string;
};

export type TransactionCardProp = BaseCardProp & {
    link: string
    type: 'transaction'
    name: string
    status: string
    address: string
}

export type CardProp = CartCardProp | CafeCardProp | TransactionCardProp;

export interface CommentProp {
    avatar: string
    name: string
    score: number
    comment: string
    timestamp: string
}
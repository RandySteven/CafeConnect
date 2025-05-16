import {Fragment} from "react";
import {GridMenu} from "@/components/Fragments/GridMenu";
import {ListCafe} from "@/components/Fragments/ListCafe";

export const HomeContainer = () => {
    return <Fragment>
        <GridMenu />
        <ListCafe />
    </Fragment>
}
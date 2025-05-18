import {IDParam} from "@/interfaces/props/ParamProp";
import {Fragment} from "react";
import {CafeDetail} from "@/components/Fragments/CafeDetail";
import {CommentSections} from "@/components/Fragments/CommentSections";
import {ProductList} from "@/components/Fragments/ProductList";
import {Box} from "@mui/material";

export const CafeDetailContainer = (id : IDParam) => {
    return <Fragment>
        <Box>
            <CafeDetail id={id.id} />
            <CommentSections id={id.id} />
            <ProductList id={id.id} />
        </Box>
    </Fragment>
}
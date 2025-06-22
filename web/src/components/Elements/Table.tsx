import {Fragment} from "react";
import {TableCell, TableHead, TableRow} from "@mui/material";
import {TableItemDetailProp} from "@/interfaces/props/TableProp";

export const TableTransactionItemHeader = () => {
    return <Fragment>
        <TableHead>
            <TableRow>
                <TableCell align="center" colSpan={3}>
                    Details
                </TableCell>
                <TableCell align="right">Price</TableCell>
            </TableRow>
            <TableRow>
                <TableCell>Desc</TableCell>
                <TableCell align="right">Qty.</TableCell>
                <TableCell align="right">Unit</TableCell>
                <TableCell align="right">Sum</TableCell>
            </TableRow>
        </TableHead>
    </Fragment>
}

export const TableTransactionItemRow = (prop : TableItemDetailProp) => {
    return <Fragment>
        {
            prop.rows.map((row, index) => (
                <TableRow key={index}>
                    <TableCell>{row.name}</TableCell>
                    <TableCell align="right">{row.qty}</TableCell>
                    <TableCell align="right">{row.price}</TableCell>
                    <TableCell align="right">{(row.price * row.qty)}</TableCell>
                </TableRow>
            ))
        }
    </Fragment>
}
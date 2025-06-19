import {Fragment} from "react";
import {ReceiptProp} from "@/interfaces/props/ReceiptProp";
import {Box} from "@mui/material";
import {Grid} from "@mui/system";

export const ReceiptInfo = (prop : ReceiptProp) => {
    return <Fragment>
        <Box>
            {/*<Grid container spacing={2}>*/}
            {/*    <Grid>*/}
            {/*        Transaction Code*/}
            {/*    </Grid>*/}
            {/*    <Grid>*/}
            {/*        {prop.transactionCode}*/}
            {/*    </Grid>*/}
            {/*    <Grid>*/}
            {/*        Transaction Status*/}
            {/*    </Grid>*/}
            {/*    <Grid>*/}
            {/*        {prop.status}*/}
            {/*    </Grid>*/}
            {/*    <Grid>*/}
            {/*        Transaction At*/}
            {/*    </Grid>*/}
            {/*    <Grid>*/}
            {/*        {prop.transactionAt}*/}
            {/*    </Grid>*/}
            {/*</Grid>*/}
            <table>
                <tbody>
                    <tr>
                        <td>
                            Transaction Code
                        </td>
                        <td>
                            {prop.transactionCode}
                        </td>
                    </tr>
                    <tr>
                        <td>
                            Transaction Status
                        </td>
                        <td>
                            {prop.status}
                        </td>
                    </tr>
                    <tr>
                        <td>
                            Transaction At
                        </td>
                        <td>
                            {prop.transactionAt}
                        </td>
                    </tr>
                </tbody>
            </table>
        </Box>
    </Fragment>
}
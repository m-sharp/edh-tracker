import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { CommanderColumn, StatColumns } from "../stats";
import { Deck } from "../types";

export default function Decks(): ReactElement {
    const decks = useLoaderData() as Array<Deck>;

    const formatColumn: GridColDef = {
        field: "format_name",
        headerName: "Format",
        minWidth: 130,
    };

    const columns = [
        CommanderColumn,
        formatColumn,
        ...StatColumns,
    ];

    return (
        <Box id="decks" style={{height: 500}}>
            <DataGrid
                rows={decks}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "points", sort: "desc"}],
                    },
                }}
            />
        </Box>
    );
}

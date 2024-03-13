import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { StatColumns } from "../stats";
import { Player } from "../types";

export default function View(): ReactElement {
    const players = useLoaderData() as Array<Player>;

    const columns: Array<GridColDef> = [
        {
            field: "name",
            headerName: "Player Name",
            renderCell: (params) => (
                <Link to={`/player/${params.row.id}`}>{params.row.name}</Link>
            ),
            hideable: false,
            flex: 1,
            minWidth: 100,
        },
        ...StatColumns,
    ];

    return (
        <Box id="players" sx={{height: 500}}>
            <DataGrid
                rows={players}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "record", sort: "desc" }],
                    },
                }}
            />
        </Box>
    );
}

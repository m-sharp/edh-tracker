import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { Player } from "./player";
import { CreatedAtColumn, StatColumns } from "../stats";

export async function getPlayers(): Promise<Array<Player>> {
    const res = await fetch(`http://localhost:8080/api/players`);
    return await res.json();
}

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
        CreatedAtColumn,
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

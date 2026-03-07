import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { StatColumns } from "../stats";
import { Player } from "../types";

export default function View(): ReactElement {
    // TODO: Should be getting players for a given pod. Will be supplanted by the new pods view described in TODOs in @app/src/routes/decks.tsx

    // TODO: Need to introduce the context of a logged in user. The logged in user should be able to, for instance:
    //      - view and managed their pods
    //      - view and manage their decks
    //      - view players in their pod
    //      - add and invite new players to their pod (will need roles for "PodManager" and "PodMember" - manager can add members)
    //      - view a pod player's decks
    //      - view a pod player's games & record within the pod
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

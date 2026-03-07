import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { CommanderColumn, StatColumns } from "../stats";
import { Deck } from "../types";

export default function Decks(): ReactElement {
    // TODO: There should be a pod landing page that supplants the Decks page.
    // TODO: The pod landing page should have various tabs that show:
    //      - A DataGrid of the decks in the pod like the one shown here - need pagination implemented for loading speed
    //      - A DataGrid or nice list of players in the pod - "PodManagers" should be able to remove members from the pod
    //      - A DataGrid of recent pod games - need pagination implemented for loading speed
    //      - A settings tab for "PodManagers" that lets them change the pod name and close (softDelete) the pod

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

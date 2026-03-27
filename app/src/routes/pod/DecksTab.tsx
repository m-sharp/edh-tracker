import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { StatColumns } from "../../components/stats";
import { Deck } from "../../types";

interface PodDecksTabProps {
    decks: Deck[];
    podId: number;
}

export default function PodDecksTab({ decks }: PodDecksTabProps): ReactElement {
    const columns: GridColDef[] = [
        {
            field: "name",
            headerName: "Deck",
            flex: 1,
            minWidth: 200,
            renderCell: (params) => (
                <Link to={`/player/${params.row.player_id}/deck/${params.row.id}`}>
                    {params.row.name}
                    {params.row.commanders && ` (${params.row.commanders.commander_name})`}
                </Link>
            ),
        },
        { field: "format_name", headerName: "Format", width: 120 },
        ...StatColumns,
    ];

    return (
        <Box sx={{ height: { xs: 400, sm: 600 }, width: "100%" }}>
            <DataGrid
                rows={decks}
                columns={columns}
                slots={{ toolbar: GridToolbar }}
                initialState={{
                    sorting: {
                        sortModel: [{ field: "record", sort: "desc" }],
                    },
                }}
            />
        </Box>
    );
}

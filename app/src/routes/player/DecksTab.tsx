import { ReactElement } from "react";
import { Box, Skeleton, Typography } from "@mui/material";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { AsyncComponentHelper } from "../../components/common";
import { GetDecksForPlayer } from "../../http";
import { CommanderColumn, StatColumns } from "../../components/stats";

interface PlayerDecksTabProps {
    playerId: number;
}

export default function PlayerDecksTab({ playerId }: PlayerDecksTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetDecksForPlayer(playerId));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={750} />;
    }
    if (error) {
        return (
            <Box sx={{ p: 2 }}>
                <Typography variant="body2" color="error">Could not load decks. Refresh to try again.</Typography>
            </Box>
        );
    }

    if (data && data.length === 0) {
        return (
            <Box sx={{ p: 2 }}>
                <Typography variant="body2" color="text.secondary">No decks yet.</Typography>
            </Box>
        );
    }

    const columns = [
        CommanderColumn,
        ...StatColumns,
        { field: "retired", headerName: "Is Retired", type: "boolean", width: 100 },
    ];

    return (
        <Box style={{ height: 750, width: "100%" }}>
            <DataGrid
                rows={data}
                columns={columns}
                slots={{ toolbar: GridToolbar }}
                initialState={{ sorting: { sortModel: [{ field: "name", sort: "asc" }] } }}
            />
        </Box>
    );
}

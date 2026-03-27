import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { Box, Button, Skeleton, Typography } from "@mui/material";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { useAuth } from "../../auth";
import { AsyncComponentHelper } from "../../components/common";
import { GetDecksForPlayer } from "../../http";
import { Deck } from "../../types";
import { CommanderColumn, StatColumns } from "../../components/stats";

interface PlayerDecksTabProps {
    playerId: number;
}

export default function PlayerDecksTab({ playerId }: PlayerDecksTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetDecksForPlayer(playerId));
    const { user } = useAuth();
    const isOwner = user?.player_id === playerId;

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

    const visibleRows = (data ?? []).filter((d: Deck) => !d.retired);

    if (data && visibleRows.length === 0) {
        return (
            <Box sx={{ p: 2, display: "flex", flexDirection: "column", gap: 2 }}>
                <Typography variant="body2" color="text.secondary">
                    {isOwner ? "No decks yet. Add a deck to get started." : "No decks yet."}
                </Typography>
                {isOwner && (
                    <Button variant="outlined" component={Link} to="/deck/new">
                        Add Deck
                    </Button>
                )}
            </Box>
        );
    }

    const columns = [
        CommanderColumn,
        ...StatColumns,
    ];

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {isOwner && (
                <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
                    <Button variant="outlined" component={Link} to="/deck/new">
                        Add Deck
                    </Button>
                </Box>
            )}
            <Box style={{ height: 750, width: "100%" }}>
                <DataGrid
                    rows={visibleRows}
                    columns={columns}
                    slots={{ toolbar: GridToolbar }}
                    initialState={{ sorting: { sortModel: [{ field: "name", sort: "asc" }] } }}
                />
            </Box>
        </Box>
    );
}

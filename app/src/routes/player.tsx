import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box, Skeleton } from "@mui/material";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { AsyncComponentHelper } from "../common";
import { GetDecksForPlayer } from "../http";
import { CommanderColumn, Record, StatColumns } from "../stats";
import { Player } from "../types";

export default function View(): ReactElement {
    const player = useLoaderData() as Player;

    return (
        <Box id="player" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <Box sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
                <h1>{player.name}</h1>
                <Record record={player.stats.record} />
            </Box>
            <Box sx={{width: "100%", display: "flex", flexDirection: "row", justifyContent: "space-evenly", py: 3}}>
                <span><strong>Games Played:</strong> {player.stats.games}</span>
                <span><strong>Total Kills:</strong> {player.stats.kills}</span>
                <span><strong>Total Points:</strong> {player.stats.points}</span>
            </Box>
            <DeckDisplay player={player} />
            <Box sx={{width: "100%", display: "flex", justifyContent: "flex-end", pt: 1}}>
                <em>Player created at: {new Date(player.created_at).toLocaleString()}</em>
            </Box>
        </Box>
    );
}

interface DeckDisplayProps {
    player: Player;
}

function DeckDisplay({ player }: DeckDisplayProps): ReactElement {
    const {data, loading, error} = AsyncComponentHelper(GetDecksForPlayer(player.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={750} />;
    }
    if (error) {
        return <span>Error Loading Player's Decks: {error.message}</span>;
    }

    const columns = [
        CommanderColumn,
        ...StatColumns,
        {
            field: "retired",
            headerName: "Is Retired",
            type: "boolean",
            width: 100,
        },
    ];

    // ToDo: Style DataGrid - https://mui.com/x/react-data-grid/style
    return (
        <Box style={{ height: 750, width: "100%" }}>
            <DataGrid
                rows={data}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "name", sort: "asc"}]
                    }
                }}
            />
        </Box>
    );
}

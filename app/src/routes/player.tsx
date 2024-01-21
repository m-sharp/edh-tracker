import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import Skeleton from "@mui/material/Skeleton";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { Deck } from "./deck";
import { AsyncComponentHelper } from "../common";
import { CommanderColumn, Record, RecordDict, StatColumns } from "../stats";

export interface Player {
    id: number;
    name: string;
    ctime: string;
    record: RecordDict;
    games: number;
    kills: number;
    points: number;
}

export async function getPlayer({ params }: LoaderFunctionArgs): Promise<Player> {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${params.playerId}`);
    return res.json();
}

export default function View(): ReactElement {
    const player = useLoaderData() as Player;

    return (
        <Box id="player" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <Box sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
                <h1>{player.name}</h1>
                <Record record={player.record} />
            </Box>
            <Box sx={{width: "100%", display: "flex", flexDirection: "row", justifyContent: "space-evenly", py: 3}}>
                <span><strong>Games Played:</strong> {player.games}</span>
                <span><strong>Total Kills:</strong> {player.kills}</span>
                <span><strong>Total Points:</strong> {player.points}</span>
            </Box>
            <DeckDisplay player={player} />
            <Box sx={{width: "100%", display: "flex", justifyContent: "flex-end", pt: 1}}>
                <em>Player created at: {new Date(player.ctime).toLocaleString()}</em>
            </Box>
        </Box>
    );
}

interface DeckDisplayProps {
    player: Player;
}

function DeckDisplay({ player }: DeckDisplayProps): ReactElement {
    const {data, loading, error} = AsyncComponentHelper(getDecksForPlayer(player.id));

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
                        sortModel: [{field: "commander", sort: "asc"}]
                    }
                }}
            />
        </Box>
    );
}

async function getDecksForPlayer(id: number): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    return await res.json();
}

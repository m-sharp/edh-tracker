import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import Skeleton from "@mui/material/Skeleton";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { Deck } from "./deck";
import { AsyncComponentHelper } from "../common";
import { CommanderColumn, CreatedAtColumn, Record, RecordDict, StatColumns } from "../stats";

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
        <div id="player">
            <h1>{player.name}&apos;s Page!</h1>
            <p>Created At: {new Date(player.ctime).toLocaleString()}</p>
            <p>Games Played: {player.games}</p>
            <p>Record: <Record record={player.record}/></p>
            <p>Total Kills: {player.kills}</p>
            <p>Total Points: {player.points}</p>
            <p>Decks:</p>
            <DeckDisplay player={player}/>
        </div>
    );
}

interface DeckDisplayProps {
    player: Player;
}

function DeckDisplay({ player }: DeckDisplayProps): ReactElement {
    const {data, loading, error} = AsyncComponentHelper(getDecksForPlayer(player.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} width={"75%"} />;
    }
    if (error) {
        return <span>Error Loading Player's Decks: {error.message}</span>;
    }

    const columns = [
        CommanderColumn,
        ...StatColumns,
        CreatedAtColumn,
        {
            field: "retired",
            headerName: "Is Retired",
            type: "boolean",
            width: 100,
        },
    ];

    // ToDo: Style DataGrid - https://mui.com/x/react-data-grid/style
    return (
        <div style={{ height: 500, width: "75%" }}>
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
        </div>
    );
}

async function getDecksForPlayer(id: number): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    return await res.json();
}
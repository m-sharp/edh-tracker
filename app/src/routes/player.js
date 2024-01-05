import { useEffect, useState } from "react";
import { useLoaderData } from "react-router-dom";
import Skeleton from "@mui/material/Skeleton";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { CommanderColumn, CreatedAtColumn, Record, StatColumns } from "../common";

export async function getPlayer({ params }) {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${params.playerId}`);
    return res.json();
}

export default function Player() {
    const player = useLoaderData();

    return (
        <div id="player">
            <h1>{player.name}&apos;s Page!</h1>
            <p>Created At: {player.ctime}</p>
            <p>Games Played: {player.games}</p>
            <p>Record: <Record record={player.record}/></p>
            <p>Total Kills: {player.kills}</p>
            <p>Total Points: {player.points}</p>
            <p>Decks:</p>
            <DeckDisplay player={player}/>
        </div>
    );
}

function DeckDisplay({ player }) {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const playerDecks = await getDecksForPlayer(player.id);
                setData(playerDecks);
                setLoading(false);
            } catch (error) {
                setError(error);
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} width={"75%"} />
    }
    if (error) {
        return <span>Error: {error.message}</span>;
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

async function getDecksForPlayer(id) {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    return await res.json();
}

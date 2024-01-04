import {useEffect, useState} from "react";
import {Link, useLoaderData} from "react-router-dom";
import Skeleton from '@mui/material/Skeleton';
import { DataGrid, GridToolbar } from '@mui/x-data-grid';

import {Record} from "../common";

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
        }

        fetchData();
    }, []);

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} width={"75%"} />
    }
    if (error) {
        return <span>Error: {error.message}</span>;
    }

    const columns = [
        {
            field: "commander",
            headerName: "Commander",
            renderCell: (params) => (
                <Link to={`/deck/${params.row.id}`}>{params.row.commander}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        {
            field: "ctime",
            headerName: "Created At",
            type: "dateTime",
            valueGetter: ({ value }) => value && new Date(value),
            minWidth: 250,
        },
        {
            field: "retired",
            headerName: "Is Retired",
            type: "boolean",
            width: 100,
        },
    ];

    // ToDo: Necessary?
    const rows = data.map(deck => ({
        "id": deck.id,
        "commander": deck.commander,
        "retired": deck.retired,
        "ctime": deck.ctime,
    }));

    // ToDo: Style DataGrid - https://mui.com/x/react-data-grid/style
    return (
        <div style={{ height: 500, width: "75%" }}>
            <DataGrid rows={rows} columns={columns} slots={{toolbar: GridToolbar}} />
        </div>
    );
}

async function getDecksForPlayer(id) {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    const decks = await res.json();

    return decks.map((deck) => ({
        id: deck.id.toString(),
        player_id: deck.player_id,
        commander: deck.commander,
        retired: deck.retired,
        ctime: deck.ctime,
    }));
}

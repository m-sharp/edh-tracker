import { Link, useLoaderData } from "react-router-dom";
import {Record} from "../common";
import {DataGrid, GridToolbar} from "@mui/x-data-grid";

export async function getPlayers() {
    const res = await fetch(`http://localhost:8080/api/players`);
    const players = await res.json();

    // ToDo: Is this really necessary?
    return players.map((player) => ({
        id: player.id.toString(),
        name: player.name,
        ctime: player.ctime,
        record: player.record,
        games: player.games,
        kills: player.kills,
        points: player.points,
    }));
}

export default function Players() {
    const players = useLoaderData();

    // ToDo: Default sorting be record or points?
    const columns = [
        {
            field: "name",
            headerName: "Player Name",
            renderCell: (params) => (
                <Link to={`/player/${params.row.id}`}>{params.row.name}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        {
            field: "record",
            headerName: "Record",
            renderCell: (params) => (
                <Record record={params.row.record}/>
            ),
            // ToDo: Custom sorting
            sortable: false,
            minWidth: 150,
        },
        {
            field: "kills",
            headerName: "Total Kills",
            type: "number",
            minWidth: 100,
        },
        {
            field: "points",
            headerName: "Total Points",
            type: "number",
            minWidth: 100,
        },
        {
            field: "games",
            headerName: "Games Played",
            type: "number",
            minWidth: 100,
        },
        {
            field: "ctime",
            headerName: "Created At",
            type: "dateTime",
            valueGetter: ({ value }) => value && new Date(value),
            minWidth: 250,
        },
    ]

    return (
        <div id="players" style={{height: 500, width: "75%"}}>
            <DataGrid rows={players} columns={columns} slots={{toolbar: GridToolbar}} />
            {/*<ul>*/}
            {/*    {players.map(player => (*/}
            {/*        <li key={player.id}>*/}
            {/*            <Link to={`/player/${player.id}`}>{player.name}</Link>*/}
            {/*        </li>*/}
            {/*    ))}*/}
            {/*</ul>*/}
        </div>
    );
}

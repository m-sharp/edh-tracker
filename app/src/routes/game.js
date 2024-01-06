import { Link, useLoaderData } from "react-router-dom";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

export async function getGame({ params }) {
    const res = await fetch(`http://localhost:8080/api/game?game_id=${params.gameId}`);
    return res.json();
}

export default function Game() {
    const game = useLoaderData();

    const columns = [
        {
            field: "place",
            headerName: "Place",
            type: "number",
            minWidth: 100,
        },
        {
            field: "commander",
            headerName: "Commander",
            renderCell: (params) => (
                <Link to={`/deck/${params.row.deck_id}`}>{params.row.commander}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        {
            field: "kill_count",
            headerName: "Kills",
            type: "number",
            minWidth: 100,
        },
        {
            field: "points",
            headerName: "Points",
            type: "number",
            minWidth: 100,
        },
    ];

    return (
        <div id="game">
            <h1>Game #{game.id}</h1>
            <p>Game Played On: {new Date(game.ctime).toLocaleString()}</p>
            <p>Description: {game.description}</p>
            <div style={{height: 355, width: "75%"}}>
                <DataGrid
                    rows={game.results}
                    columns={columns}
                    slots={{toolbar: GridToolbar}}
                    initialState={{
                        sorting: {
                            sortModel: [{field: "place", sort: "asc"}],
                        },
                    }}
                />
            </div>
        </div>
    );
}

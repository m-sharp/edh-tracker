import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";
import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { Game } from "../matches";

export async function getGame({ params }: LoaderFunctionArgs): Promise<Game> {
    const res = await fetch(`http://localhost:8080/api/game?game_id=${params.gameId}`);
    return res.json();
}

export default function View(): ReactElement {
    const game = useLoaderData() as Game;

    const columns: Array<GridColDef> = [
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
        <Box id="game" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <h1>Game #{game.id}</h1>
            <em>{new Date(game.ctime).toLocaleString()}</em>
            <p>Description: {game.description}</p>
            <Box sx={{height: 355, width: "100%"}}>
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
            </Box>
        </Box>
    );
}

import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { Game } from "../types";

export default function View(): ReactElement {
    // TODO: A better route would probably be /pod/<PodID>/game/<GameID>
    // TODO: Should display: Game description, decks that played DataGrid, CreatedAt
    // TODO: Logged in user that is a "PodManager" should be able to do the following:
    //      - edit the Game description
    //      - edit individual Game Result via a modal. Edit kills, place, points, deck etc
    //      - Add and remove Game Results from the game
    //      - delete the Game

    const game = useLoaderData() as Game;

    const columns: Array<GridColDef> = [
        {
            field: "place",
            headerName: "Place",
            type: "number",
            minWidth: 100,
        },
        {
            field: "deck_name",
            headerName: "Deck",
            renderCell: (params) => (
                <Link to={`/deck/${params.row.deck_id}`}>{params.row.deck_name}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        {
            field: "commander_name",
            headerName: "Commander",
            flex: 1,
            renderCell: (params) => params.row.commander_name ?? "—",
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
            {/* TODO: Would probably be better to say "{PodName} Game #{game.id}" */}
            <h1>Game #{game.id}</h1>
            <em>{new Date(game.created_at).toLocaleString()}</em>
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

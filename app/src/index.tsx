import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { CssBaseline } from "@mui/material";

import ErrorPage from "./routes/error"
import DeckView, { getDeck } from "./routes/deck";
import DecksView, { getDecks } from "./routes/decks";
import GameView, { getGame } from "./routes/game";
import GamesView, { getGames } from "./routes/games";
import PlayerView, { getPlayer } from "./routes/player";
import PlayersView, { getPlayers } from "./routes/players";
import Root from "./routes/root";

import "./styles.css";

// ToDo: Change to jsx declarations?
const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            {
                path: "decks",
                element: <DecksView />,
                loader: getDecks,
            },
            {
                path: "deck/:deckId",
                element: <DeckView />,
                loader: getDeck,
            },
            {
                path: "games",
                element: <GamesView />,
                loader: getGames,
            },
            {
                path: "game/:gameId",
                element: <GameView />,
                loader: getGame,
            },
            {
                path: "players",
                element: <PlayersView />,
                loader: getPlayers,
            },
            {
                path: "player/:playerId",
                element: <PlayerView />,
                loader: getPlayer,
            },
        ],
    },
]);

createRoot(document.getElementById("root") as HTMLElement).render(
    <StrictMode>
        <CssBaseline enableColorScheme />
        <RouterProvider router={router} />
    </StrictMode>
);

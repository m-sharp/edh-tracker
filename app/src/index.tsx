import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { CssBaseline } from "@mui/material";

import {
    GetDeck, GetDecks,
    GetGame, GetGames,
    GetPlayer, GetPlayers,
} from "./http";
import ErrorPage from "./routes/error"
import DeckView from "./routes/deck";
import DecksView from "./routes/decks";
import GameView from "./routes/game";
import GamesView from "./routes/games";
import NewGameView, { createGame } from "./routes/new";
import PlayerView from "./routes/player";
import PlayersView from "./routes/players";
import Root from "./routes/root";

import "./styles.css";

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            {
                path: "decks",
                element: <DecksView />,
                loader: GetDecks,
            },
            {
                path: "deck/:deckId",
                element: <DeckView />,
                loader: GetDeck,
            },
            {
                path: "games",
                element: <GamesView />,
                loader: GetGames,
            },
            {
                path: "game/:gameId",
                element: <GameView />,
                loader: GetGame,
            },
            {
                path: "players",
                element: <PlayersView />,
                loader: GetPlayers,
            },
            {
                path: "player/:playerId",
                element: <PlayerView />,
                loader: GetPlayer,
            },
            {
                path: "new-game",
                element: <NewGameView />,
                loader: GetDecks,
                action: createGame,
            }
        ],
    },
]);

createRoot(document.getElementById("root") as HTMLElement).render(
    <StrictMode>
        <CssBaseline enableColorScheme />
        <RouterProvider router={router} />
    </StrictMode>
);

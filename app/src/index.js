import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import ErrorPage from "./routes/error"
import Deck, { getDeck } from "./routes/deck";
import Decks, { getDecks } from "./routes/decks";
import Player, { getPlayer } from "./routes/player";
import Players, { getPlayers } from "./routes/players";
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
                element: <Decks />,
                loader: getDecks,
            },
            {
                path: "deck/:deckId",
                element: <Deck />,
                loader: getDeck,
            },
            {
                path: "players",
                element: <Players />,
                loader: getPlayers,
            },
            {
                path: "player/:playerId",
                element: <Player />,
                loader: getPlayer,
            },
        ],
    },
]);

createRoot(document.getElementById("root")).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>
);

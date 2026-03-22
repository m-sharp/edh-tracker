import { ReactElement, StrictMode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, Navigate, RouterProvider, useNavigate } from "react-router-dom";
import { CssBaseline, Typography } from "@mui/material";

import { AuthProvider, useAuth } from "./auth";
import {
    GetDeck,
    GetGame,
    GetNewDeckInfo,
    GetPlayer,
    GetPodsForPlayer,
} from "./http";
import ErrorPage from "./routes/error"
import DeckView from "./routes/deck";
import GameView from "./routes/game";
import LoginPage from "./routes/login";
import NewGameView, { createGame } from "./routes/new";
import PlayerView from "./routes/player";
import RequireAuth from "./routes/RequireAuth";
import Root from "./routes/root";

import "./styles.css";

// TODO: Move these views into their own files?
function HomeView(): ReactElement {
    const { user } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (!user) {
            return;
        }

        GetPodsForPlayer(user.player_id).then((pods) => {
            if (pods.length > 0) {
                navigate(`/pod/${pods[0].id}`, {replace: true});
            }
        });
    }, [user]);

    return <Typography>No pods yet. Create your first pod to get started.</Typography>;
}

function PodView(): ReactElement {
    return <Typography>Coming soon</Typography>;
}

function JoinView(): ReactElement {
    return <Typography>Coming soon</Typography>;
}

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            {
                index: true,
                element: <RequireAuth><HomeView /></RequireAuth>,
            },
            {
                path: "login",
                element: <LoginPage />,
            },
            {
                path: "join",
                element: <JoinView />,
            },
            {
                path: "pod/:podId",
                element: <RequireAuth><PodView /></RequireAuth>,
            },
            {
                path: "pod/:podId/new-game",
                element: <RequireAuth><NewGameView /></RequireAuth>,
                loader: GetNewDeckInfo,
                action: createGame,
            },
            {
                path: "pod/:podId/game/:gameId",
                element: <RequireAuth><GameView /></RequireAuth>,
                loader: GetGame,
            },
            {
                path: "player/:playerId",
                element: <RequireAuth><PlayerView /></RequireAuth>,
                loader: GetPlayer,
            },
            {
                path: "player/:playerId/deck/:deckId",
                element: <RequireAuth><DeckView /></RequireAuth>,
                loader: GetDeck,
            },
        ],
    },
]);

createRoot(document.getElementById("root") as HTMLElement).render(
    <StrictMode>
        <CssBaseline enableColorScheme />
        <AuthProvider>
            <RouterProvider router={router} />
        </AuthProvider>
    </StrictMode>
);

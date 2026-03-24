import { ReactElement, StrictMode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, Navigate, RouterProvider, useNavigate } from "react-router-dom";
import { Typography } from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider } from "@mui/material/styles";
import theme from "./theme";

import { AuthProvider, useAuth } from "./auth";
import {
    GetDeck,
    GetPlayer,
    GetPodsForPlayer,
} from "./http";
import ErrorPage from "./routes/error"
import DeckView from "./routes/deck";
import GameView, { gameLoader } from "./routes/game";
import JoinView from "./routes/join";
import LoginPage from "./routes/login";
import NewGameView, { newGameLoader, createGame } from "./routes/new";
import PlayerView from "./routes/player";
import PodView, { podLoader } from "./routes/pod";
import RequireAuth from "./routes/RequireAuth";
import Root from "./routes/root";

import "./styles.css";

// TODO: Move HomeView into its own file
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

    // TODO: Will need to link to where you can actually add one
    // TODO: Loading blip always shows this before data comes in
    return <Typography>No pods yet. Create your first pod to get started.</Typography>;
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
                loader: podLoader,
            },
            {
                path: "pod/:podId/new-game",
                element: <RequireAuth><NewGameView /></RequireAuth>,
                loader: newGameLoader,
                action: createGame,
            },
            {
                path: "pod/:podId/game/:gameId",
                element: <RequireAuth><GameView /></RequireAuth>,
                loader: gameLoader,
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
        <ThemeProvider theme={theme}>
            <CssBaseline enableColorScheme />
            <AuthProvider>
                <RouterProvider router={router} />
            </AuthProvider>
        </ThemeProvider>
    </StrictMode>
);

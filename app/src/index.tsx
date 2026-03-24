import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";
import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider } from "@mui/material/styles";
import theme from "./theme";

import { AuthProvider } from "./auth";
import {
    GetDeck,
    GetPlayer,
} from "./http";
import ErrorPage from "./routes/error"
import DeckView from "./routes/deck";
import GameView, { gameLoader } from "./routes/game";
import HomeView from "./routes/home";
import JoinView from "./routes/join";
import LoginPage from "./routes/login";
import NewGameView, { newGameLoader, createGame } from "./routes/new";
import PlayerView from "./routes/player";
import PodView, { podLoader } from "./routes/pod";
import RequireAuth from "./routes/RequireAuth";
import Root from "./routes/root";

import "./styles.css";

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

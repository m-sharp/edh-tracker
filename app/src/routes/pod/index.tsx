import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import { Box, Button, Typography } from "@mui/material";

import { useAuth } from "../../auth";
import { GetDecksForPod, GetGamesForPod, GetPod, GetPlayersForPod } from "../../http";
import { Deck, Game, PaginatedResponse, PlayerWithRole, Pod } from "../../types";
import TabbedLayout from "../../components/TabbedLayout";
import PodDecksTab from "./DecksTab";
import PodPlayersTab from "./PlayersTab";
import PodGamesTab from "./GamesTab";
import PodSettingsTab from "./SettingsTab";

interface PodLoaderData {
    pod: Pod;
    players: PlayerWithRole[];
    decks: PaginatedResponse<Deck>;
    games: PaginatedResponse<Game>;
}

export async function podLoader({ params }: LoaderFunctionArgs): Promise<PodLoaderData> {
    const podId = Number(params.podId);
    const [pod, players, decks, games] = await Promise.all([
        GetPod(podId),
        GetPlayersForPod(podId),
        GetDecksForPod(podId, 25, 0),
        GetGamesForPod(podId, 25, 0),
    ]);
    return { pod, players, decks, games };
}

export default function PodView(): ReactElement {
    const { pod, players, decks, games } = useLoaderData() as PodLoaderData;
    const { user } = useAuth();

    const currentUserRole = user
        ? players.find((p) => p.id === user.player_id)?.role ?? null
        : null;
    const isManager = currentUserRole === "manager";

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between", wrap: "true", mb: 2 }}>
                <Typography variant="h4">{pod.name}</Typography>
                <Button variant="contained" component={Link} to={`/pod/${pod.id}/new-game`}>
                    New Game
                </Button>
            </Box>
            <TabbedLayout
                queryKey="podTab"
                tabs={[
                    { id: "decks", label: "Decks", content: <PodDecksTab decks={decks} podId={pod.id} /> },
                    { id: "players", label: "Players", content: <PodPlayersTab players={players} podId={pod.id} isManager={isManager} /> },
                    { id: "games", label: "Games", content: <PodGamesTab games={games} podId={pod.id} /> },
                    { id: "settings", label: "Settings", content: <PodSettingsTab pod={pod} />, hidden: !isManager },
                ]}
            />
        </Box>
    );
}

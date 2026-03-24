import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box, Typography } from "@mui/material";

import { useAuth } from "../../auth";
import TabbedLayout from "../../components/TabbedLayout";
import { Record } from "../../components/stats";
import { Player } from "../../types";
import PlayerOverviewTab from "./OverviewTab";
import PlayerDecksTab from "./DecksTab";
import PlayerGamesTab from "./GamesTab";
import PlayerSettingsTab from "./SettingsTab";

export default function PlayerView(): ReactElement {
    const player = useLoaderData() as Player;
    const { user } = useAuth();

    const isOwner = user?.player_id === player.id;

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Typography variant="h4" sx={{ mb: 1 }}>{player.name}</Typography>
            <Record record={player.stats.record} />
            <TabbedLayout
                queryKey="playerTab"
                tabs={[
                    { id: "overview", label: "Overview", content: <PlayerOverviewTab player={player} /> },
                    { id: "decks", label: "Decks", content: <PlayerDecksTab playerId={player.id} /> },
                    { id: "games", label: "Games", content: <PlayerGamesTab playerId={player.id} /> },
                    { id: "settings", label: "Settings", content: <PlayerSettingsTab player={player} />, hidden: !isOwner },
                ]}
            />
        </Box>
    );
}

import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box, Typography } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";

import { useAuth } from "../../auth";
import TabbedLayout from "../../components/TabbedLayout";
import { Record } from "../../components/stats";
import { Deck } from "../../types";
import DeckOverviewTab from "./OverviewTab";
import DeckGamesTab from "./GamesTab";
import DeckSettingsTab from "./SettingsTab";

export default function DeckView(): ReactElement {
    const deck = useLoaderData() as Deck;
    const { user } = useAuth();

    const isOwner = user?.player_id === deck.player_id;

    const commanderLabel = deck.commanders
        ? deck.commanders.partner_commander_name
            ? `${deck.commanders.commander_name} / ${deck.commanders.partner_commander_name}`
            : deck.commanders.commander_name
        : null;

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Typography variant="h4" sx={{ mb: 0.5 }}>{deck.name}</Typography>
            {commanderLabel && (
                <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>{commanderLabel}</Typography>
            )}
            <Record record={deck.stats.record} />
            {deck.retired && (
                <Box sx={{ display: "flex", alignItems: "center", gap: 0.5, mt: 1 }}>
                    <DeleteIcon fontSize="small" />
                    <Typography variant="body2">Retired</Typography>
                </Box>
            )}
            <TabbedLayout
                queryKey="deckTab"
                tabs={[
                    { id: "overview", label: "Overview", content: <DeckOverviewTab deck={deck} /> },
                    { id: "games", label: "Games", content: <DeckGamesTab deckId={deck.id} commanderName={deck.commanders?.commander_name} /> },
                    { id: "settings", label: "Settings", content: <DeckSettingsTab deck={deck} />, hidden: !isOwner },
                ]}
            />
        </Box>
    );
}

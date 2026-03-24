import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { Box, Typography } from "@mui/material";

import { Record } from "../../components/stats";
import { Deck } from "../../types";

interface DeckOverviewTabProps {
    deck: Deck;
}

export default function DeckOverviewTab({ deck }: DeckOverviewTabProps): ReactElement {
    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
            <Box sx={{ display: "flex", flexDirection: "row", flexWrap: "wrap", gap: 2 }}>
                <Typography variant="body1"><strong>Games Played:</strong> {deck.stats.games}</Typography>
                <Typography variant="body1"><strong>Total Kills:</strong> {deck.stats.kills}</Typography>
                <Typography variant="body1"><strong>Total Points:</strong> {deck.stats.points}</Typography>
            </Box>
            <Record record={deck.stats.record} />
            <Typography variant="body1"><strong>Format:</strong> {deck.format_name}</Typography>
            <Typography variant="body1">
                <strong>Owner:</strong> <Link to={`/player/${deck.player_id}`}>{deck.player_name}</Link>
            </Typography>
            <Typography variant="body2" color="text.secondary">
                Created: {new Date(deck.created_at).toLocaleString()}
            </Typography>
        </Box>
    );
}

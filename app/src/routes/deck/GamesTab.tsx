import { ReactElement } from "react";
import { Box, Skeleton, Typography } from "@mui/material";

import { AsyncComponentHelper } from "../../components/common";
import { GetGamesForDeck } from "../../http";
import { MatchesDisplay } from "../../components/matches";

interface DeckGamesTabProps {
    deckId: number;
    commanderName?: string;
}

export default function DeckGamesTab({ deckId, commanderName }: DeckGamesTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetGamesForDeck(deckId));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} />;
    }
    if (error) {
        return <Typography color="error">Error loading games: {error.message}</Typography>;
    }

    return (
        <Box sx={{ height: 500, width: "100%" }}>
            <MatchesDisplay games={data} targetCommander={commanderName} />
        </Box>
    );
}

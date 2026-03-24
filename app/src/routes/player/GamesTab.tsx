import { ReactElement } from "react";
import { Box, Skeleton, Typography } from "@mui/material";

import { AsyncComponentHelper } from "../../components/common";
import { GetGamesForPlayer } from "../../http";
import { MatchesDisplay } from "../../components/matches";
import { Game } from "../../types";

interface PlayerGamesTabProps {
    playerId: number;
}

export default function PlayerGamesTab({ playerId }: PlayerGamesTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetGamesForPlayer(playerId));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={400} />;
    }
    if (error) {
        return (
            <Box sx={{ p: 2 }}>
                <Typography variant="body2" color="error">Could not load games. Refresh to try again.</Typography>
            </Box>
        );
    }

    if (data && data.length === 0) {
        return (
            <Box sx={{ p: 2 }}>
                <Typography variant="body2" color="text.secondary">No games yet.</Typography>
            </Box>
        );
    }

    return (
        <Box sx={{ height: 600, width: "100%" }}>
            <MatchesDisplay games={data as Game[]} />
        </Box>
    );
}

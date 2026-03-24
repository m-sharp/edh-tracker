import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { Box, Skeleton, Typography } from "@mui/material";

import { AsyncComponentHelper } from "../../components/common";
import { GetPodsForPlayer } from "../../http";
import { Record } from "../../components/stats";
import { Player, Pod } from "../../types";

interface PlayerOverviewTabProps {
    player: Player;
}

export default function PlayerOverviewTab({ player }: PlayerOverviewTabProps): ReactElement {
    const { data: pods, loading, error } = AsyncComponentHelper(GetPodsForPlayer(player.id));

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <Box sx={{ display: "flex", flexDirection: "row", justifyContent: "space-evenly", flexWrap: "wrap", gap: 2, px: 1, py: 1 }}>
                <Typography variant="body1"><strong>Games Played:</strong> {player.stats.games}</Typography>
                <Typography variant="body1"><strong>Total Kills:</strong> {player.stats.kills}</Typography>
                <Typography variant="body1"><strong>Total Points:</strong> {player.stats.points}</Typography>
            </Box>
            <Box>
                <Typography variant="h6">Pods</Typography>
                {loading && <Skeleton variant="text" />}
                {error && <Typography variant="body2" color="error">Could not load pods. Refresh to try again.</Typography>}
                {pods && pods.length === 0 && (
                    <Typography variant="body2">No pods yet.</Typography>
                )}
                {pods && pods.map((pod: Pod) => (
                    <Box key={pod.id}>
                        <Link to={`/pod/${pod.id}`}>{pod.name}</Link>
                    </Box>
                ))}
            </Box>
            <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
                <Typography variant="body2" color="text.secondary">Created: {new Date(player.created_at).toLocaleString()}</Typography>
            </Box>
        </Box>
    );
}

import { ReactElement, useState } from "react";
import { Link } from "react-router-dom";
import { Box, Button, Chip, List, ListItem, ListItemText } from "@mui/material";

import { useAuth } from "../../auth";
import { DeletePodPlayer, GetPlayersForPod, PatchPodPlayerRole } from "../../http";
import { PlayerWithRole } from "../../types";

interface PodPlayersTabProps {
    players: PlayerWithRole[];
    podId: number;
    isManager: boolean;
}

export default function PodPlayersTab({ players: initialPlayers, podId, isManager }: PodPlayersTabProps): ReactElement {
    const { user } = useAuth();
    const [players, setPlayers] = useState(initialPlayers);

    const refetchPlayers = async () => {
        const updated = await GetPlayersForPod(podId);
        setPlayers(updated);
    };

    const handlePromote = async (playerId: number) => {
        await PatchPodPlayerRole(podId, playerId);
        await refetchPlayers();
    };

    const handleRemove = async (playerId: number) => {
        await DeletePodPlayer(podId, playerId);
        await refetchPlayers();
    };

    // TODO: Use icons w/ tooltips for promote/remove buttons?
    // TODO: Title case Manager vs Member roles coming back from backend
    return (
        <List>
            {players.map((p) => (
                <ListItem
                    key={p.id}
                    secondaryAction={
                        isManager && user?.player_id !== p.id ? (
                            <Box sx={{ display: "flex", gap: 1 }}>
                                {p.role === "member" && (
                                    <Button size="small" onClick={() => handlePromote(p.id)} sx={{ minHeight: 44 }}>
                                        Promote
                                    </Button>
                                )}
                                <Button size="small" color="error" onClick={() => handleRemove(p.id)} sx={{ minHeight: 44 }}>
                                    Remove
                                </Button>
                            </Box>
                        ) : null
                    }
                >
                    <ListItemText
                        primary={<Link to={`/player/${p.id}`}>{p.name}</Link>}
                        secondary={
                            <Chip
                                label={p.role === "manager" ? "Manager" : "Member"}
                                size="small"
                                sx={{ mt: 0.5 }}
                            />
                        }
                    />
                </ListItem>
            ))}
        </List>
    );
}

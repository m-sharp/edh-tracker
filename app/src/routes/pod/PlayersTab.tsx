import { ReactElement, useState } from "react";
import { Link } from "react-router-dom";
import {
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    List,
    ListItem,
    ListItemText
} from "@mui/material";

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
    const [confirmAction, setConfirmAction] = useState<{ type: "promote" | "remove"; player: PlayerWithRole } | null>(null);

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

    const handleConfirm = async () => {
        if (!confirmAction) return;
        if (confirmAction.type === "promote") {
            await handlePromote(confirmAction.player.id);
        } else {
            await handleRemove(confirmAction.player.id);
        }
        setConfirmAction(null);
    };

    // TODO: Use icons w/ tooltips for promote/remove buttons?
    // TODO: Title case Manager vs Member roles coming back from backend
    return (
        <>
            <List>
                {players.map((p) => (
                    <ListItem
                        key={p.id}
                        secondaryAction={
                            isManager && user?.player_id !== p.id ? (
                                <Box sx={{ display: "flex", gap: 1 }}>
                                    {p.role === "member" && (
                                        <Button
                                            variant="contained"
                                            size="small"
                                            onClick={() => setConfirmAction({ type: "promote", player: p })}
                                            sx={{ minHeight: 44 }}
                                        >
                                            Promote
                                        </Button>
                                    )}
                                    <Button
                                        variant="contained"
                                        size="small"
                                        color="error"
                                        onClick={() => setConfirmAction({ type: "remove", player: p })}
                                        sx={{ minHeight: 44 }}
                                    >
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

            <Dialog open={confirmAction !== null} onClose={() => setConfirmAction(null)}>
                <DialogTitle>
                    {confirmAction?.type === "promote" ? "Promote player?" : "Remove player?"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmAction?.type === "promote"
                            ? `Promote ${confirmAction.player.name} to manager? Managers can edit pod settings and manage members.`
                            : `Remove ${confirmAction?.player.name} from this pod? This action cannot be undone.`}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmAction(null)}>Cancel</Button>
                    <Button
                        color={confirmAction?.type === "promote" ? "primary" : "error"}
                        onClick={handleConfirm}
                    >
                        Confirm
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}

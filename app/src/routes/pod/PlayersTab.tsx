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
    Paper,
    Typography
} from "@mui/material";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import PersonOffIcon from "@mui/icons-material/PersonOff";

import { useAuth } from "../../auth";
import { DeletePodPlayer, GetPlayersForPod, PatchPodPlayerRole } from "../../http";
import { Record } from "../../components/stats";
import { TooltipIconButton } from "../../components/TooltipIcon";
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

    return (
        <>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
                {players.map((p) => (
                    <Paper key={p.id} elevation={2} sx={{ p: 2 }}>
                        <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                                <Typography variant="body1">
                                    <Link to={`/player/${p.id}`}>{p.name}</Link>
                                </Typography>
                                {p.role === "manager" && (
                                    <Chip label="Manager" size="small" />
                                )}
                            </Box>
                            {isManager && user?.player_id !== p.id && (
                                <Box sx={{ display: "flex", gap: 1 }}>
                                    {p.role === "member" && (
                                        <TooltipIconButton
                                            title="Promote to Manager"
                                            onClick={() => setConfirmAction({ type: "promote", player: p })}
                                            icon={<PersonAddIcon />}
                                        />
                                    )}
                                    <TooltipIconButton
                                        title="Remove from pod"
                                        onClick={() => setConfirmAction({ type: "remove", player: p })}
                                        icon={<PersonOffIcon />}
                                        color="error"
                                    />
                                </Box>
                            )}
                        </Box>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                            <Record record={p.stats.record} /> {"\u2022"} {p.stats.points} pts {"\u2022"} {p.stats.kills} kills
                        </Typography>
                    </Paper>
                ))}
            </Box>

            <Dialog open={confirmAction !== null} onClose={() => setConfirmAction(null)}>
                <DialogTitle>
                    {confirmAction?.type === "promote"
                        ? `Promote ${confirmAction.player.name}?`
                        : `Remove ${confirmAction?.player.name}?`}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmAction?.type === "promote"
                            ? `Promote ${confirmAction.player.name} to manager? Managers can edit pod settings and manage members.`
                            : `Remove ${confirmAction?.player.name} from this pod? This action cannot be undone.`}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmAction(null)}>Never mind</Button>
                    <Button
                        variant="contained"
                        color={confirmAction?.type === "promote" ? "primary" : "error"}
                        onClick={handleConfirm}
                    >
                        {confirmAction?.type === "promote" ? "Make Manager" : "Remove"}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}

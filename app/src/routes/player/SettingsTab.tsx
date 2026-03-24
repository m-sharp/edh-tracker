import { ReactElement, useState } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    Skeleton,
    TextField,
    Typography,
} from "@mui/material";

import { AsyncComponentHelper } from "../../components/common";
import {
    GetPodsForPlayer,
    PatchPlayer,
    PostPod,
    PostPodLeave,
} from "../../http";
import { Player, Pod } from "../../types";

interface PlayerSettingsTabProps {
    player: Player;
}

export default function PlayerSettingsTab({ player }: PlayerSettingsTabProps): ReactElement {
    const navigate = useNavigate();
    const [name, setName] = useState(player.name);
    const [nameError, setNameError] = useState<string | null>(null);
    const [newPodName, setNewPodName] = useState("");
    const [createPodError, setCreatePodError] = useState<string | null>(null);
    const [leaveConfirmPodId, setLeaveConfirmPodId] = useState<number | null>(null);
    const [leaveError, setLeaveError] = useState<string | null>(null);
    const { data: pods, loading: podsLoading, error: podsError } = AsyncComponentHelper(GetPodsForPlayer(player.id));

    const handleSaveName = async () => {
        setNameError(null);
        try {
            await PatchPlayer(player.id, name);
            window.location.reload();
        } catch {
            setNameError("Failed to update name. Try again.");
        }
    };

    const handleLeave = async () => {
        if (leaveConfirmPodId === null) return;
        setLeaveError(null);
        try {
            await PostPodLeave(leaveConfirmPodId);
            setLeaveConfirmPodId(null);
            window.location.reload();
        } catch (e: any) {
            setLeaveConfirmPodId(null);
            if (e?.status === 403) {
                setLeaveError("Promote another member to manager before leaving.");
            } else {
                setLeaveError("Failed to leave pod. Try again.");
            }
        }
    };

    const handleCreatePod = async () => {
        setCreatePodError(null);
        try {
            const pod = await PostPod(newPodName);
            navigate(`/pod/${pod.id}`);
        } catch {
            setCreatePodError("Failed to create pod. Try again.");
        }
    };

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, maxWidth: 500 }}>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Display Name</Typography>
                <Box sx={{ display: "flex", gap: 1 }}>
                    <TextField
                        label="Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        size="small"
                    />
                    <Button variant="contained" onClick={handleSaveName}>Save Name</Button>
                </Box>
                {nameError && <Typography color="error" variant="body2">{nameError}</Typography>}
            </Box>

            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Your Pods</Typography>
                {leaveError && <Typography color="error" variant="body2">{leaveError}</Typography>}
                {podsLoading && <Skeleton variant="text" />}
                {podsError && <Typography color="error">Error loading pods.</Typography>}
                {pods && pods.length === 0 && (
                    <Typography variant="body2">No pods yet.</Typography>
                )}
                {pods && pods.map((pod: Pod) => (
                    <Box key={pod.id} sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                        <Link to={`/pod/${pod.id}`}>{pod.name}</Link>
                        <Button size="small" color="error" sx={{ minHeight: 44 }} onClick={() => setLeaveConfirmPodId(pod.id)}>
                            Leave
                        </Button>
                    </Box>
                ))}
            </Box>

            <Divider sx={{ my: 2 }} />

            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Create New Pod</Typography>
                <Box sx={{ display: "flex", gap: 1 }}>
                    <TextField
                        label="Pod Name"
                        value={newPodName}
                        onChange={(e) => setNewPodName(e.target.value)}
                        size="small"
                    />
                    <Button
                        variant="contained"
                        onClick={handleCreatePod}
                        disabled={!newPodName.trim()}
                    >
                        Create
                    </Button>
                </Box>
                {createPodError && <Typography color="error" variant="body2">{createPodError}</Typography>}
            </Box>

            <Dialog open={leaveConfirmPodId !== null} onClose={() => setLeaveConfirmPodId(null)}>
                <DialogTitle>Leave pod?</DialogTitle>
                <DialogContent>
                    <Typography>Are you sure you want to leave this pod?</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setLeaveConfirmPodId(null)}>Cancel</Button>
                    <Button color="error" onClick={handleLeave}>Leave</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

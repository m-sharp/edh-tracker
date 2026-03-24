import { ReactElement, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    TextField,
    Typography,
} from "@mui/material";

import { DeletePod, PatchPod, PostPodInvite } from "../../http";
import { Pod } from "../../types";

interface PodSettingsTabProps {
    pod: Pod;
}

export default function PodSettingsTab({ pod }: PodSettingsTabProps): ReactElement {
    const navigate = useNavigate();
    const [name, setName] = useState(pod.name);
    const [nameError, setNameError] = useState<string | null>(null);
    const [inviteLink, setInviteLink] = useState<string | null>(null);
    const [inviteError, setInviteError] = useState<string | null>(null);
    const [deleteOpen, setDeleteOpen] = useState(false);

    const handleSaveName = async () => {
        setNameError(null);
        try {
            await PatchPod(pod.id, name);
            window.location.reload();
        } catch {
            setNameError("Failed to update name.");
        }
    };

    const handleGenerateInvite = async () => {
        setInviteError(null);
        try {
            const { invite_code } = await PostPodInvite(pod.id);
            setInviteLink(`${window.location.origin}/join?code=${invite_code}`);
        } catch {
            setInviteError("Failed to generate invite link.");
        }
    };

    const handleDelete = async () => {
        await DeletePod(pod.id);
        navigate("/", { replace: true });
    };

    // TODO: Icon w/ tooltip for Save & Copy
    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, maxWidth: 500 }}>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                    <TextField
                        label="Pod Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        size="small"
                        sx={{ flex: 1, minWidth: 160 }}
                    />
                    <Button variant="contained" onClick={handleSaveName}>Save Pod Name</Button>
                </Box>
                {nameError && <Typography color="error" variant="body2">{nameError}</Typography>}
            </Box>
            <Box>
                <Button variant="outlined" onClick={handleGenerateInvite}>
                    Generate Invite Link
                </Button>
                {inviteError && <Typography color="error" variant="body2">{inviteError}</Typography>}
                {inviteLink && (
                    <Box sx={{ mt: 1, display: "flex", gap: 1, alignItems: "center", flexWrap: "wrap" }}>
                        <Typography variant="body2" sx={{ wordBreak: "break-all", flex: 1, minWidth: 160 }}>
                            {inviteLink}
                        </Typography>
                        <Button size="small" onClick={() => navigator.clipboard.writeText(inviteLink)}>
                            Copy
                        </Button>
                    </Box>
                )}
            </Box>
            <Box>
                <Button variant="outlined" color="error" onClick={() => setDeleteOpen(true)}>
                    Delete Pod
                </Button>
            </Box>
            <Dialog open={deleteOpen} onClose={() => setDeleteOpen(false)}>
                <DialogTitle>Delete "{pod.name}"?</DialogTitle>
                <DialogContent>
                    <Typography>This action cannot be undone.</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteOpen(false)}>Cancel</Button>
                    <Button color="error" onClick={handleDelete}>Delete</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

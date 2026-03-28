import { ReactElement, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Box,
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    TextField,
    Typography
} from "@mui/material";
import { useAuth } from "../../auth";
import { GetPodsForPlayer, PostPod } from "../../http";
import { Pod } from "../../types";

export default function HomeView(): ReactElement {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [pods, setPods] = useState<Pod[]>([]);
    const [createPodOpen, setCreatePodOpen] = useState(false);
    const [podName, setPodName] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const [createError, setCreateError] = useState<string | null>(null);

    useEffect(() => {
        if (!user) return;
        GetPodsForPlayer(user.player_id).then((result) => {
            if (result.length > 0) {
                navigate(`/pod/${result[0].id}`, { replace: true });
            } else {
                setPods(result);
                setLoading(false);
            }
        }).catch(() => setLoading(false));
    }, [user, navigate]);

    const handleCreatePod = async () => {
        setCreateError(null);
        setSubmitting(true);
        try {
            const { id } = await PostPod(podName);
            navigate(`/pod/${id}`);
        } catch {
            setCreateError("Failed to create pod. Try again.");
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <Box sx={{ display: "flex", justifyContent: "center", pt: 4 }}>
                <CircularProgress />
            </Box>
        );
    }

    if (pods.length === 0) {
        return (
            <>
                <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: { xs: 4, sm: 8 }, gap: 2 }}>
                    <Typography variant="h4">Welcome to EDH Tracker</Typography>
                    <Typography variant="body1" color="text.secondary" textAlign="center">
                        Create your first pod or ask a friend for an invite link.
                    </Typography>
                    <Button variant="contained" onClick={() => setCreatePodOpen(true)}>
                        Create a Pod
                    </Button>
                </Box>
                <Dialog open={createPodOpen} onClose={() => setCreatePodOpen(false)} maxWidth="xs" fullWidth>
                    <DialogTitle>Create a New Pod</DialogTitle>
                    <DialogContent>
                        <TextField
                            label="Pod Name"
                            value={podName}
                            onChange={(e) => setPodName(e.target.value)}
                            fullWidth
                            autoFocus
                            sx={{ mt: 1 }}
                            disabled={submitting}
                        />
                        {createError && <Typography color="error" variant="body2" sx={{ mt: 1 }}>{createError}</Typography>}
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => setCreatePodOpen(false)} disabled={submitting}>Discard</Button>
                        <Button
                            variant="contained"
                            disabled={!podName.trim() || submitting}
                            onClick={handleCreatePod}
                        >
                            {submitting ? <CircularProgress size={20} /> : "Create Pod"}
                        </Button>
                    </DialogActions>
                </Dialog>
            </>
        );
    }

    return <></>;
}

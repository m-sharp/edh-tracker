import { ReactElement, ReactNode, useEffect, useState } from "react";
import { Link, Outlet, useNavigate, useParams } from "react-router-dom";
import {
    AppBar,
    Avatar,
    Box,
    Button,
    CircularProgress,
    Container,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    MenuItem,
    Select,
    SelectChangeEvent,
    TextField,
    Toolbar,
    Typography
} from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { useAuth } from "../auth";
import { GetPodsForPlayer, PostPod } from "../http";
import { Pod } from "../types";
import SvgIconPlayingCards from "../components/SvgIconPlayingCards";
import { TooltipIconButton } from "../components/TooltipIcon";

export default function Root(): ReactElement {
    return (
        <Box sx={{ display: "flex", width: "auto" }}>
            <DrawerAppBar />
            <Container id="detail" component="main" sx={{ p: 3, width: "90%", bgcolor: "background.default", mt: 12, mb: 5 }} maxWidth="xl">
                <Outlet />
            </Container>
        </Box>
    );
}

function DrawerAppBar(): ReactElement {
    const { user, loading, logout } = useAuth();
    const navigate = useNavigate();

    let authSection: ReactNode = null;
    if (!loading) {
        if (!user) {
            authSection = (
                <Button href="/api/auth/google" sx={{ color: "white" }}>
                    Sign in with Google
                </Button>
            );
        } else {
            authSection = (
                <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    <PodSelector playerId={user.player_id} />
                    <Link to={`/player/${user.player_id}`} style={{ display: "flex", alignItems: "center", textDecoration: "none", color: "white" }}>
                        <Avatar
                            src={user.avatar_url ?? undefined}
                            alt={user.display_name ?? "User"}
                            sx={{ width: 32, height: 32, mr: 1 }}
                        />
                        <Typography variant="body2">{user.display_name}</Typography>
                    </Link>
                    <TooltipIconButton
                        title={"Logout"}
                        onClick={() => logout().then(() => navigate("/login"))}
                        icon={<LogoutIcon />}
                    />
                </Box>
            );
        }
    }

    return (
        <AppBar position="fixed">
            <Container maxWidth="xl">
                <Toolbar disableGutters>
                    <Box sx={{ display: "flex", mr: 2 }}>
                        <Link to="/" style={{ display: "flex" }}>
                            <SvgIconPlayingCards />
                        </Link>
                    </Box>
                    <Typography
                        variant="h6"
                        noWrap
                        sx={{
                            mr: 2,
                            display: { xs: "none", sm: "flex" },
                            fontWeight: 700,
                            letterSpacing: ".3rem",
                        }}
                    >
                        <Link to={`/`} style={{textDecoration: "none", color: "white"}}>EDH Tracker</Link>
                    </Typography>
                    <Box sx={{ flexGrow: 1 }} />
                    {authSection}
                </Toolbar>
            </Container>
        </AppBar>
    );
}

function PodSelector({ playerId }: { playerId: number }): ReactElement {
    const [pods, setPods] = useState<Pod[]>([]);
    const { podId } = useParams<{ podId?: string }>();
    const navigate = useNavigate();
    const [createPodOpen, setCreatePodOpen] = useState(false);
    const [podName, setPodName] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const [createError, setCreateError] = useState<string | null>(null);

    const selectedPodId = podId ?? localStorage.getItem("lastPodId") ?? "";

    useEffect(() => {
        GetPodsForPlayer(playerId).then(setPods).catch(() => setPods([]));
    }, [playerId]);

    const handleChange = (e: SelectChangeEvent) => {
        const id = e.target.value;
        if (id === "create-new") {
            setCreatePodOpen(true);
            return;
        }
        localStorage.setItem("lastPodId", id);
        navigate(`/pod/${id}`);
    };

    const handleCreatePod = async () => {
        setCreateError(null);
        setSubmitting(true);
        try {
            const { id } = await PostPod(podName);
            setCreatePodOpen(false);
            setPodName("");
            setSubmitting(false);
            navigate(`/pod/${id}`);
        } catch {
            setCreateError("Failed to create pod. Try again.");
            setSubmitting(false);
        }
    };

    return (
        <>
            <Select
                value={selectedPodId}
                onChange={handleChange}
                displayEmpty
                size="small"
                sx={{ color: "white", mr: 1 }}
            >
                {pods.length === 0
                    ? <MenuItem value="" disabled>No pods</MenuItem>
                    : pods.map(p => <MenuItem key={p.id} value={String(p.id)}>{p.name}</MenuItem>)
                }
                <Divider />
                <MenuItem value="create-new" sx={{ color: "primary.main" }}>Create new pod</MenuItem>
            </Select>
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

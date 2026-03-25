import { ReactElement, ReactNode, useEffect, useState } from "react";
import { Link, Outlet, useNavigate, useParams } from "react-router-dom";
import {
    AppBar,
    Avatar,
    Box,
    Button,
    Container,
    IconButton,
    MenuItem,
    Select,
    SelectChangeEvent,
    Toolbar,
    Tooltip,
    Typography
} from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { useAuth } from "../auth";
import { GetPodsForPlayer } from "../http";
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

// TODO: Add mobile menu icon and link menu
// TODO: Mobile view for all tables
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
                    <Box sx={{ display: "flex", mr: 2 }}><SvgIconPlayingCards /></Box>
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

    const selectedPodId = podId ?? localStorage.getItem("lastPodId") ?? "";

    useEffect(() => {
        GetPodsForPlayer(playerId).then(setPods).catch(() => setPods([]));
    }, [playerId]);

    const handleChange = (e: SelectChangeEvent) => {
        const id = e.target.value;
        localStorage.setItem("lastPodId", id);
        navigate(`/pod/${id}`);
    };

    return (
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
        </Select>
    );
}


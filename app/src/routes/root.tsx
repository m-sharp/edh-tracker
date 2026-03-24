import { ReactElement, ReactNode, useEffect, useState } from "react";
import { Link, Outlet, useNavigate, useParams } from "react-router-dom";
import {
    AppBar,
    Avatar,
    Box,
    Button,
    Container,
    MenuItem,
    Select,
    SelectChangeEvent,
    SvgIcon,
    Toolbar,
    Typography
} from "@mui/material";
import { useAuth } from "../auth";
import { GetPodsForPlayer } from "../http";
import { Pod } from "../types";

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
                    <Button
                        color="inherit"
                        onClick={() => logout().then(() => navigate("/login"))}
                    >
                        Logout
                    </Button>
                </Box>
            );
        }
    }

    return (
        <AppBar position="fixed">
            <Container maxWidth="xl">
                <Toolbar disableGutters>
                    <SvgIconPlayingCards />
                    <Typography
                        variant="h6"
                        noWrap
                        sx={{
                            mr: 2,
                            display: "flex",
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

function SvgIconPlayingCards(): ReactElement {
    // Via https://fonts.google.com/icons?selected=Material%20Symbols%20Outlined%3Aplaying_cards%3AFILL%400%3Bwght%40400%3BGRAD%400%3Bopsz%4024
    return (
        <SvgIcon sx={{ display: "flex", mr: 2 }}>
            <svg
                xmlns="http://www.w3.org/2000/svg"
                height="24"
                viewBox="0 -960 960 960"
                width="24"
                strokeWidth={1.5}
                stroke="currentColor"
                fill="white"
            >
                <path d="m608-368 46-166-142-98-46 166 142 98ZM160-207l-33-16q-31-13-42-44.5t3-62.5l72-156v279Zm160 87q-33 0-56.5-24T240-201v-239l107 294q3 7 5 13.5t7 12.5h-39Zm206-5q-31 11-62-3t-42-45L245-662q-11-31 3-61.5t45-41.5l301-110q31-11 61.5 3t41.5 45l178 489q11 31-3 61.5T827-235L526-125Zm-28-75 302-110-179-490-301 110 178 490Zm62-300Z" />
            </svg>
        </SvgIcon>
    );
}

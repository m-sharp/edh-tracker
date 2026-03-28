import { ReactElement, useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { Box, Button, CircularProgress, Typography } from "@mui/material";

import { useAuth } from "../auth";
import {API_BASE_URL, PostPodJoin} from "../http";
import SvgIconPlayingCards from "../components/SvgIconPlayingCards";

export default function JoinView(): ReactElement {
    const [searchParams] = useSearchParams();
    const code = searchParams.get("code");
    const { user, loading } = useAuth();
    const navigate = useNavigate();

    const [error, setError] = useState<string | null>(null);
    const [joining, setJoining] = useState(false);

    useEffect(() => {
        if (loading) {
            return;
        }

        if (!user) {
            const redirectUrl = `/join${code ? `?code=${code}` : ""}`;
            window.location.href = `${API_BASE_URL}/api/auth/google?redirect=${encodeURIComponent(redirectUrl)}`;
            return;
        }

        if (!code) {
            return;
        }

        setJoining(true);
        PostPodJoin(code)
            .then((pod) => navigate(`/pod/${pod.id}`, { replace: true }))
            .catch((err) => {
                setError(err.message || "Invalid or expired invite code.");
                setJoining(false);
            });
    }, [loading, user, code]);

    if (!code) {
        return (
            <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 8, gap: 2 }}>
                <SvgIconPlayingCards fontSize={40} />
                <Typography variant="h6">No invite code provided.</Typography>
                <Button component={Link} to="/" variant="outlined" size="medium">Go home</Button>
            </Box>
        );
    }

    if (error) {
        return (
            <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 8, gap: 2 }}>
                <SvgIconPlayingCards fontSize={40} />
                <Typography variant="h6">Something went wrong</Typography>
                <Typography variant="body1" color="error">{error}</Typography>
                <Button component={Link} to="/" variant="outlined" size="medium">Go home</Button>
            </Box>
        );
    }

    return (
        <Box sx={{ display: "flex", justifyContent: "center", pt: 8 }}>
            <CircularProgress />
        </Box>
    );
}

import { ReactElement, useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { Box, CircularProgress, Typography } from "@mui/material";

import { useAuth } from "../auth";
import {API_BASE_URL, PostPodJoin} from "../http";

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

    if (loading || joining) {
        return (
            <Box sx={{ display: "flex", justifyContent: "center", pt: 8 }}>
                <CircularProgress />
            </Box>
        );
    }

    if (!code) {
        return (
            <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 8, gap: 2 }}>
                <Typography variant="h6">No invite code provided.</Typography>
                <Link to="/">Go home</Link>
            </Box>
        );
    }

    if (error) {
        return (
            <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 8, gap: 2 }}>
                <Typography variant="h6" color="error">{error}</Typography>
                <Link to="/">Go home</Link>
            </Box>
        );
    }

    // TODO: This is weird ultimate return. Bad setup here with all these if checks
    return (
        <Box sx={{ display: "flex", justifyContent: "center", pt: 8 }}>
            <CircularProgress />
        </Box>
    );
}

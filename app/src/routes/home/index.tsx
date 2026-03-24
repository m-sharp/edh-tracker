import { ReactElement, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Box, CircularProgress, Typography } from "@mui/material";
import { useAuth } from "../../auth";
import { GetPodsForPlayer } from "../../http";
import { Pod } from "../../types";

export default function HomeView(): ReactElement {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [pods, setPods] = useState<Pod[]>([]);

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

    if (loading) {
        return (
            <Box sx={{ display: "flex", justifyContent: "center", pt: 4 }}>
                <CircularProgress />
            </Box>
        );
    }

    if (pods.length === 0) {
        return (
            <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", pt: 4 }}>
                <Typography>No pods yet. Create your first pod or ask a manager for an invite link.</Typography>
            </Box>
        );
    }

    return <></>;
}

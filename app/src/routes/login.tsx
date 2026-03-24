import { ReactElement } from "react";
import { Box, Button, Typography } from "@mui/material";

import { API_BASE_URL } from "../http";
import SvgIconPlayingCards from "../components/SvgIconPlayingCards";

export default function LoginPage(): ReactElement {
    return (
        <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
            <Box sx={{ mb: 2 }}>
                <SvgIconPlayingCards fontSize={48} />
            </Box>
            <Typography variant="h4" gutterBottom>EDH Tracker</Typography>
            <Button variant="contained" href={`${API_BASE_URL}/api/auth/google`} size="large">
                Sign in with Google
            </Button>
        </Box>
    );
}

import { ReactElement } from "react";
import { Box, Button, Typography } from "@mui/material";

import { API_BASE_URL } from "../http";
import SvgIconPlayingCards from "../components/SvgIconPlayingCards";

export default function LoginPage(): ReactElement {
    return (
        <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column", justifyContent: "flex-start", alignItems: "center", pt: { xs: 4, sm: 8 } }}>
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

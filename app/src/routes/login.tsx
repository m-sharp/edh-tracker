import { ReactElement } from "react";
import { Box, Button, Typography } from "@mui/material";

import { API_BASE_URL } from "../http";

export default function LoginPage(): ReactElement {
    return (
        <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", mt: 8 }}>
            <Typography variant="h4" gutterBottom>EDH Tracker</Typography>
            <Button variant="contained" href={`${API_BASE_URL}/api/auth/google`} size="large">
                Sign in with Google
            </Button>
        </Box>
    );
}

import { ReactElement, ReactNode } from "react";
import { Box, CircularProgress } from "@mui/material";
import { Navigate } from "react-router-dom";
import { useAuth } from "../auth";

export default function RequireAuth({ children }: { children: ReactNode }): ReactElement {
    const { user, loading } = useAuth();
    if (loading) {
        return (
            <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", pt: 4 }}>
                <CircularProgress />
            </Box>
        );
    }
    if (!user) {
        return <Navigate to="/login" replace/>
    };
    return <>{children}</>;
}

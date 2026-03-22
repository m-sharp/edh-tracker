import { ReactElement, ReactNode } from "react";
import { CircularProgress } from "@mui/material";
import { Navigate } from "react-router-dom";
import { useAuth } from "../auth";

export default function RequireAuth({ children }: { children: ReactNode }): ReactElement {
    const { user, loading } = useAuth();
    if (loading) {
        return <CircularProgress/>
    }
    if (!user) {
        return <Navigate to="/login" replace/>
    };
    return <>{children}</>;
}

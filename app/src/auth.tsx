import { createContext, ReactElement, ReactNode, useContext, useEffect, useState } from "react";
import { GetMe, Logout } from "./http";

export interface AuthUser {
    id: number;
    player_id: number;
    display_name: string | null;
    avatar_url: string | null;
}

interface AuthContextValue {
    user: AuthUser | null;
    loading: boolean;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue>({
    user: null,
    loading: true,
    logout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }): ReactElement {
    const [user, setUser] = useState<AuthUser | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        GetMe()
            .then(setUser)
            .catch(() => setUser(null))
            .finally(() => setLoading(false));
    }, []);

    const logout = async () => {
        await Logout();
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, loading, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth(): AuthContextValue {
    return useContext(AuthContext);
}

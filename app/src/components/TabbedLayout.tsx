import { ReactElement } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Box, CircularProgress, Tab, Tabs } from "@mui/material";

interface TabConfig {
    id: string;
    label: string;
    content: ReactElement;
    hidden?: boolean;
}

interface TabbedLayoutProps {
    queryKey: string;
    tabs: TabConfig[];
    loading?: boolean;
}

export default function TabbedLayout({ queryKey, tabs, loading }: TabbedLayoutProps): ReactElement {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const visibleTabs = tabs.filter((t) => !t.hidden);
    const activeId = searchParams.get(queryKey);
    const activeIndex = Math.max(0, visibleTabs.findIndex((t) => t.id === activeId));
    const activeTab = visibleTabs[activeIndex];

    const handleChange = (_: React.SyntheticEvent, newIndex: number) => {
        const params = new URLSearchParams(searchParams);
        params.set(queryKey, visibleTabs[newIndex].id);
        navigate(`?${params.toString()}`, { replace: true });
    };

    return (
        <Box>
            <Tabs
                value={activeIndex}
                onChange={handleChange}
                variant="scrollable"
                scrollButtons="auto"
                sx={{ mb: 2 }}
            >
                {visibleTabs.map((t) => (
                    <Tab key={t.id} label={t.label} />
                ))}
            </Tabs>
            {loading ? (
                <Box sx={{ display: "flex", justifyContent: "center", py: 4 }}>
                    <CircularProgress />
                </Box>
            ) : (
                activeTab?.content
            )}
        </Box>
    );
}

import { ReactElement, useState } from "react";
import { Link, useLoaderData, useNavigate } from "react-router-dom";
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Skeleton,
    Tab,
    Tabs,
    TextField,
    Typography,
} from "@mui/material";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { useAuth } from "../auth";
import { AsyncComponentHelper } from "../common";
import {
    GetDecksForPlayer,
    GetGamesForPlayer,
    GetPodsForPlayer,
    PatchPlayer,
    PostPod,
    PostPodLeave,
} from "../http";
import { MatchesDisplay } from "../matches";
import { CommanderColumn, Record, StatColumns } from "../stats";
import { Game, Player, Pod } from "../types";

export default function PlayerView(): ReactElement {
    const player = useLoaderData() as Player;
    const { user } = useAuth();
    // TODO: Use common tab component
    const [tab, setTab] = useState(0);

    const isOwnProfile = user?.player_id === player.id;

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Typography variant="h4" sx={{ mb: 1 }}>{player.name}</Typography>
            <Record record={player.stats.record} />
            <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 2, mt: 2 }}>
                <Tab label="Overview" />
                <Tab label="Decks" />
                <Tab label="Games" />
                {isOwnProfile && <Tab label="Settings" />}
            </Tabs>
            {tab === 0 && <PlayerOverviewTab player={player} />}
            {tab === 1 && <PlayerDecksTab player={player} />}
            {tab === 2 && <PlayerGamesTab player={player} />}
            {tab === 3 && isOwnProfile && <PlayerSettingsTab player={player} />}
        </Box>
    );
}

interface PlayerTabProps {
    player: Player;
}

// TODO: File system restructure
function PlayerOverviewTab({ player }: PlayerTabProps): ReactElement {
    const { data: pods, loading, error } = AsyncComponentHelper(GetPodsForPlayer(player.id));

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <Box sx={{ display: "flex", flexDirection: "row", justifyContent: "space-evenly", py: 1 }}>
                <span><strong>Games Played:</strong> {player.stats.games}</span>
                <span><strong>Total Kills:</strong> {player.stats.kills}</span>
                <span><strong>Total Points:</strong> {player.stats.points}</span>
            </Box>
            <Box>
                <Typography variant="h6">Pods</Typography>
                {/* TODO: Probably want some common skeleton and error text handling for tab components - would be nice to have everywhere */}
                {loading && <Skeleton variant="text" />}
                {error && <Typography color="error">Error loading pods: {error.message}</Typography>}
                {pods && pods.length === 0 && (
                    <Typography variant="body2">No pods yet.</Typography>
                )}
                {pods && pods.map((pod: Pod) => (
                    <Box key={pod.id}>
                        <Link to={`/pod/${pod.id}`}>{pod.name}</Link>
                    </Box>
                ))}
            </Box>
            <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
                <em>Player created at: {new Date(player.created_at).toLocaleString()}</em>
            </Box>
        </Box>
    );
}

function PlayerDecksTab({ player }: PlayerTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetDecksForPlayer(player.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={750} />;
    }
    if (error) {
        return <span>Error Loading Player's Decks: {error.message}</span>;
    }

    const columns = [
        CommanderColumn,
        ...StatColumns,
        { field: "retired", headerName: "Is Retired", type: "boolean", width: 100 },
    ];

    return (
        <Box style={{ height: 750, width: "100%" }}>
            <DataGrid
                rows={data}
                columns={columns}
                slots={{ toolbar: GridToolbar }}
                initialState={{ sorting: { sortModel: [{ field: "name", sort: "asc" }] } }}
            />
        </Box>
    );
}

function PlayerGamesTab({ player }: PlayerTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetGamesForPlayer(player.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={400} />;
    }
    if (error) {
        return <span>Error Loading Player's Games: {error.message}</span>;
    }

    return (
        <Box sx={{ height: 600, width: "100%" }}>
            <MatchesDisplay games={data as Game[]} />
        </Box>
    );
}

function PlayerSettingsTab({ player }: PlayerTabProps): ReactElement {
    const navigate = useNavigate();
    const [name, setName] = useState(player.name);
    const [nameError, setNameError] = useState<string | null>(null);
    const [newPodName, setNewPodName] = useState("");
    const [leaveConfirmPodId, setLeaveConfirmPodId] = useState<number | null>(null);
    const [leaveError, setLeaveError] = useState<string | null>(null);
    const { data: pods, loading: podsLoading, error: podsError } = AsyncComponentHelper(GetPodsForPlayer(player.id));

    const handleSaveName = async () => {
        setNameError(null);
        try {
            await PatchPlayer(player.id, name);
            window.location.reload();
        } catch {
            setNameError("Failed to update name.");
        }
    };

    const handleLeave = async () => {
        if (leaveConfirmPodId === null) return;
        setLeaveError(null);
        try {
            await PostPodLeave(leaveConfirmPodId);
            setLeaveConfirmPodId(null);
            window.location.reload();
        } catch (e: any) {
            setLeaveConfirmPodId(null);
            if (e?.status === 403) {
                setLeaveError("Promote another member to manager before leaving.");
            } else {
                setLeaveError("Failed to leave pod.");
            }
        }
    };

    // TODO: This doesn't make sense hinden behind player settings - should live in a tab within the pod page
    const handleCreatePod = async () => {
        const pod = await PostPod(newPodName);
        navigate(`/pod/${pod.id}`);
    };

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, maxWidth: 500 }}>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Display Name</Typography>
                <Box sx={{ display: "flex", gap: 1 }}>
                    <TextField
                        label="Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        size="small"
                    />
                    <Button variant="contained" onClick={handleSaveName}>Save</Button>
                </Box>
                {nameError && <Typography color="error" variant="body2">{nameError}</Typography>}
            </Box>

            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Your Pods</Typography>
                {leaveError && <Typography color="error" variant="body2">{leaveError}</Typography>}
                {podsLoading && <Skeleton variant="text" />}
                {podsError && <Typography color="error">Error loading pods.</Typography>}
                {pods && pods.length === 0 && (
                    <Typography variant="body2">No pods yet.</Typography>
                )}
                {pods && pods.map((pod: Pod) => (
                    <Box key={pod.id} sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                        <Link to={`/pod/${pod.id}`}>{pod.name}</Link>
                        <Button size="small" color="error" onClick={() => setLeaveConfirmPodId(pod.id)}>
                            Leave
                        </Button>
                    </Box>
                ))}
            </Box>

            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Create New Pod</Typography>
                <Box sx={{ display: "flex", gap: 1 }}>
                    <TextField
                        label="Pod Name"
                        value={newPodName}
                        onChange={(e) => setNewPodName(e.target.value)}
                        size="small"
                    />
                    <Button
                        variant="contained"
                        onClick={handleCreatePod}
                        disabled={!newPodName.trim()}
                    >
                        Create
                    </Button>
                </Box>
            </Box>

            <Dialog open={leaveConfirmPodId !== null} onClose={() => setLeaveConfirmPodId(null)}>
                <DialogTitle>Leave pod?</DialogTitle>
                <DialogContent>
                    <Typography>Are you sure you want to leave this pod?</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setLeaveConfirmPodId(null)}>Cancel</Button>
                    <Button color="error" onClick={handleLeave}>Leave</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

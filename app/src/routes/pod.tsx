import { ReactElement, useState } from "react";
import { Link, useLoaderData, useNavigate, useParams } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import {
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    List,
    ListItem,
    ListItemText,
    Tab,
    Tabs,
    TextField,
    Typography,
} from "@mui/material";
import { DataGrid, GridColDef, GridPaginationModel, GridToolbar } from "@mui/x-data-grid";

import { useAuth } from "../auth";
import {
    DeletePod,
    DeletePodPlayer,
    GetDecksForPod,
    GetGamesForPod,
    GetPod,
    GetPlayersForPod,
    PatchPod,
    PatchPodPlayerRole,
    PostPodInvite,
} from "../http";
import { StatColumns } from "../stats";
import { Deck, Game, PaginatedResponse, PlayerWithRole, Pod } from "../types";

interface PodLoaderData {
    pod: Pod;
    players: PlayerWithRole[];
    decks: PaginatedResponse<Deck>;
    games: PaginatedResponse<Game>;
}

export async function podLoader({ params }: LoaderFunctionArgs): Promise<PodLoaderData> {
    const podId = Number(params.podId);
    const [pod, players, decks, games] = await Promise.all([
        GetPod(podId),
        GetPlayersForPod(podId),
        GetDecksForPod(podId, 25, 0),
        GetGamesForPod(podId, 25, 0),
    ]);
    return { pod, players, decks, games };
}

// TODO: Restructure files -> app/src/routes/pod/view.jsx|decks.jsx|players.jsx|new.jsx|...
export default function PodView(): ReactElement {
    const { pod, players, decks, games } = useLoaderData() as PodLoaderData;
    const { user } = useAuth();
    // TODO: General helper component for rendering tabs - takes in a prop that maps tab name -> component to show when tab selected
    const [tab, setTab] = useState(0);

    // TODO: This is ugly, helper function that just returns an isManager bool
    const currentUserRole = user
        ? players.find((p) => p.id === user.player_id)?.role ?? null
        : null;
    const isManager = currentUserRole === "manager";

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Typography variant="h4" sx={{ mb: 2 }}>{pod.name}</Typography>
            <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 2 }}>
                <Tab label="Decks" />
                <Tab label="Players" />
                <Tab label="Games" />
                {isManager && <Tab label="Settings" />}
            </Tabs>
            {tab === 0 && <PodDecksTab initialData={decks} />}
            {tab === 1 && <PodPlayersTab players={players} isManager={isManager} />}
            {tab === 2 && <PodGamesTab initialData={games} />}
            {tab === 3 && isManager && <PodSettingsTab pod={pod} />}
        </Box>
    );
}

interface PodDecksTabProps {
    initialData: PaginatedResponse<Deck>;
}

function PodDecksTab({ initialData }: PodDecksTabProps): ReactElement {
    const { podId } = useParams();
    const [rows, setRows] = useState<Deck[]>(initialData.items);
    const [rowCount, setRowCount] = useState(initialData.total);
    const [loading, setLoading] = useState(false);
    const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({ page: 0, pageSize: 25 });

    const handlePaginationChange = async (model: GridPaginationModel) => {
        setPaginationModel(model);
        setLoading(true);
        const data = await GetDecksForPod(Number(podId), model.pageSize, model.page * model.pageSize);
        setRows(data.items);
        setRowCount(data.total);
        setLoading(false);
    };

    const columns: GridColDef[] = [
        {
            field: "name",
            headerName: "Deck",
            flex: 1,
            minWidth: 200,
            renderCell: (params) => (
                <Link to={`/player/${params.row.player_id}/deck/${params.row.id}`}>
                    {params.row.name}
                    {params.row.commanders && ` (${params.row.commanders.commander_name})`}
                </Link>
            ),
        },
        { field: "format_name", headerName: "Format", width: 120 },
        ...StatColumns,
    ];

    return (
        <Box sx={{ height: 600, width: "100%" }}>
            <DataGrid
                rows={rows}
                columns={columns}
                loading={loading}
                paginationMode="server"
                rowCount={rowCount}
                paginationModel={paginationModel}
                onPaginationModelChange={handlePaginationChange}
                slots={{ toolbar: GridToolbar }}
            />
        </Box>
    );
}

interface PodPlayersTabProps {
    players: PlayerWithRole[];
    isManager: boolean;
}

function PodPlayersTab({ players: initialPlayers, isManager }: PodPlayersTabProps): ReactElement {
    const { podId } = useParams();
    const { user } = useAuth();
    const [players, setPlayers] = useState(initialPlayers);

    const refetchPlayers = async () => {
        const updated = await GetPlayersForPod(Number(podId));
        setPlayers(updated);
    };

    const handlePromote = async (playerId: number) => {
        await PatchPodPlayerRole(Number(podId), playerId);
        await refetchPlayers();
    };

    const handleRemove = async (playerId: number) => {
        await DeletePodPlayer(Number(podId), playerId);
        await refetchPlayers();
    };

    // TODO: Use icons w/ tooltips for promote/remove buttons?
    // TODO: Title case Manager vs Member roles coming back from backend
    return (
        <List>
            {players.map((p) => (
                <ListItem
                    key={p.id}
                    secondaryAction={
                        isManager && user?.player_id !== p.id ? (
                            <Box sx={{ display: "flex", gap: 1 }}>
                                {p.role === "member" && (
                                    <Button size="small" onClick={() => handlePromote(p.id)}>
                                        Promote
                                    </Button>
                                )}
                                <Button size="small" color="error" onClick={() => handleRemove(p.id)}>
                                    Remove
                                </Button>
                            </Box>
                        ) : null
                    }
                >
                    <ListItemText
                        primary={<Link to={`/player/${p.id}`}>{p.name}</Link>}
                        secondary={
                            <Chip
                                label={p.role === "manager" ? "Manager" : "Member"}
                                size="small"
                                sx={{ mt: 0.5 }}
                            />
                        }
                    />
                </ListItem>
            ))}
        </List>
    );
}

interface PodGamesTabProps {
    initialData: PaginatedResponse<Game>;
}

function PodGamesTab({ initialData }: PodGamesTabProps): ReactElement {
    const { podId } = useParams();
    const navigate = useNavigate();
    const [rows, setRows] = useState<Game[]>(initialData.items);
    const [rowCount, setRowCount] = useState(initialData.total);
    const [loading, setLoading] = useState(false);
    const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({ page: 0, pageSize: 25 });

    const handlePaginationChange = async (model: GridPaginationModel) => {
        setPaginationModel(model);
        setLoading(true);
        const data = await GetGamesForPod(Number(podId), model.pageSize, model.page * model.pageSize);
        setRows(data.items);
        setRowCount(data.total);
        setLoading(false);
    };

    const columns: GridColDef[] = [
        {
            field: "id",
            headerName: "Game #",
            width: 100,
            renderCell: (params) => (
                <Link to={`/pod/${podId}/game/${params.row.id}`}>#{params.row.id}</Link>
            ),
        },
        { field: "description", headerName: "Description", flex: 1, minWidth: 200 },
        {
            field: "created_at",
            headerName: "Date",
            width: 180,
            valueFormatter: (params) => new Date(params.value).toLocaleDateString(),
        },
        {
            field: "results",
            headerName: "Participants",
            width: 120,
            valueGetter: (params) => params.row.results?.length ?? 0,
        },
    ];

    return (
        <Box>
            <Box sx={{ display: "flex", justifyContent: "flex-end", mb: 2 }}>
                <Button variant="contained" onClick={() => navigate(`/pod/${podId}/new-game`)}>
                    New Game
                </Button>
            </Box>
            <Box sx={{ height: 600, width: "100%" }}>
                <DataGrid
                    rows={rows}
                    columns={columns}
                    loading={loading}
                    paginationMode="server"
                    rowCount={rowCount}
                    paginationModel={paginationModel}
                    onPaginationModelChange={handlePaginationChange}
                    slots={{ toolbar: GridToolbar }}
                />
            </Box>
        </Box>
    );
}

interface PodSettingsTabProps {
    pod: Pod;
}

function PodSettingsTab({ pod }: PodSettingsTabProps): ReactElement {
    const navigate = useNavigate();
    const [name, setName] = useState(pod.name);
    const [inviteLink, setInviteLink] = useState<string | null>(null);
    const [deleteOpen, setDeleteOpen] = useState(false);

    const handleSaveName = async () => {
        await PatchPod(pod.id, name);
        window.location.reload();
    };

    const handleGenerateInvite = async () => {
        const { invite_code } = await PostPodInvite(pod.id);
        setInviteLink(`${window.location.origin}/join?code=${invite_code}`);
    };

    const handleDelete = async () => {
        await DeletePod(pod.id);
        navigate("/", { replace: true });
    };

    // TODO: Icon w/ tooltip for Save & Copy
    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, maxWidth: 500 }}>
            <Box sx={{ display: "flex", gap: 1 }}>
                <TextField
                    label="Pod Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    size="small"
                />
                <Button variant="contained" onClick={handleSaveName}>Save</Button>
            </Box>
            <Box>
                <Button variant="outlined" onClick={handleGenerateInvite}>
                    Generate Invite Link
                </Button>
                {inviteLink && (
                    <Box sx={{ mt: 1, display: "flex", gap: 1, alignItems: "center" }}>
                        <Typography variant="body2" sx={{ wordBreak: "break-all" }}>
                            {inviteLink}
                        </Typography>
                        <Button size="small" onClick={() => navigator.clipboard.writeText(inviteLink)}>
                            Copy
                        </Button>
                    </Box>
                )}
            </Box>
            <Box>
                <Button variant="outlined" color="error" onClick={() => setDeleteOpen(true)}>
                    Delete Pod
                </Button>
            </Box>
            <Dialog open={deleteOpen} onClose={() => setDeleteOpen(false)}>
                <DialogTitle>Delete "{pod.name}"?</DialogTitle>
                <DialogContent>
                    <Typography>This action cannot be undone.</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteOpen(false)}>Cancel</Button>
                    <Button color="error" onClick={handleDelete}>Delete</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

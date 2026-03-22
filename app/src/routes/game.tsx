import { ReactElement, useState } from "react";
import { Link, useLoaderData, useNavigate } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import {
    Autocomplete,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    IconButton,
    Stack,
    TextField,
    Typography,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import { DataGrid, GridColDef } from "@mui/x-data-grid";

import { useAuth } from "../auth";
import {
    DeleteGame,
    DeleteGameResult,
    GetDecksForPod,
    GetGame,
    GetPlayersForPod,
    GetPod,
    PatchGame,
    PatchGameResult,
    PostGameResult,
} from "../http";
import { Deck, Game, GameResult, PlayerWithRole, Pod } from "../types";

interface GameLoaderData {
    game: Game;
    pod: Pod;
    players: PlayerWithRole[];
    decks: Deck[];
}

export async function gameLoader(args: LoaderFunctionArgs): Promise<GameLoaderData> {
    const podId = Number(args.params.podId);
    const [game, pod, players, decks] = await Promise.all([
        GetGame(args),
        GetPod(podId),
        GetPlayersForPod(podId),
        GetDecksForPod(podId, 1000, 0),
    ]);
    return { game, pod, players, decks: decks.items };
}

interface GameDescriptionProps {
    game: Game;
    isManager: boolean;
}

function GameDescription({ game, isManager }: GameDescriptionProps): ReactElement {
    const [editing, setEditing] = useState(false);
    const [value, setValue] = useState(game.description);

    async function handleSave() {
        await PatchGame(game.id, value);
        window.location.reload();
    }

    if (editing) {
        return (
            <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 2 }}>
                <TextField
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
                    label="Description"
                    size="small"
                />
                <Button onClick={handleSave} variant="contained" size="small">Save</Button>
                <Button
                    onClick={() => {
                        setValue(game.description);
                        setEditing(false);
                    }}
                    size="small"
                >
                    Cancel
                </Button>
            </Stack>
        );
    }

    // TODO: Make sure icon button has a tooltip
    return (
        <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 2 }}>
            {game.description
                ? <Typography>Description: {game.description}</Typography>
                : <Typography color="text.secondary"><em>No description</em></Typography>
            }
            {isManager && (
                <IconButton size="small" onClick={() => setEditing(true)}>
                    <EditIcon fontSize="small" />
                </IconButton>
            )}
        </Stack>
    );
}

interface EditResultModalProps {
    result: GameResult;
    decks: Deck[];
    onClose: () => void;
}

function EditResultModal({ result, decks, onClose }: EditResultModalProps): ReactElement {
    const [place, setPlace] = useState(result.place);
    const [killCount, setKillCount] = useState(result.kill_count);
    const [deckId, setDeckId] = useState(result.deck_id);

    async function handleSave() {
        await PatchGameResult(result.id, { place, kill_count: killCount, deck_id: deckId });
        window.location.reload();
    }

    return (
        <Dialog open onClose={onClose}>
            <DialogTitle>Edit Result</DialogTitle>
            <DialogContent>
                <Stack spacing={2} sx={{ mt: 1, minWidth: 300 }}>
                    <TextField
                        label="Place"
                        type="number"
                        value={place}
                        onChange={(e) => setPlace(Number(e.target.value))}
                        size="small"
                    />
                    <TextField
                        label="Kills"
                        type="number"
                        value={killCount}
                        onChange={(e) => setKillCount(Number(e.target.value))}
                        size="small"
                    />
                    <Autocomplete
                        options={decks}
                        getOptionLabel={(d) => d.name}
                        value={decks.find((d) => d.id === deckId) ?? null}
                        onChange={(_, d) => {
                            if (d) {
                                setDeckId(d.id);
                            }
                        }}
                        renderInput={(params) => <TextField {...params} label="Deck" size="small" />}
                    />
                </Stack>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSave} variant="contained">Save</Button>
            </DialogActions>
        </Dialog>
    );
}

interface RemoveResultDialogProps {
    result: GameResult;
    onClose: () => void;
}

function RemoveResultDialog({ result, onClose }: RemoveResultDialogProps): ReactElement {
    async function handleRemove() {
        await DeleteGameResult(result.id);
        window.location.reload();
    }

    return (
        <Dialog open onClose={onClose}>
            <DialogTitle>Remove result?</DialogTitle>
            <DialogContent>
                <Typography>Remove {result.deck_name} from this game? This action cannot be undone.</Typography>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button color="error" onClick={handleRemove}>Remove</Button>
            </DialogActions>
        </Dialog>
    );
}

interface AddResultModalProps {
    gameId: number;
    podId: number;
    decks: Deck[];
    players: PlayerWithRole[];
    onClose: () => void;
}

function AddResultModal({ gameId, decks, players, onClose }: AddResultModalProps): ReactElement {
    const [playerId, setPlayerId] = useState<number | null>(null);
    const [deckId, setDeckId] = useState<number | null>(null);
    const [place, setPlace] = useState(1);
    const [killCount, setKillCount] = useState(0);

    async function handleAdd() {
        if (!playerId || !deckId) return;
        await PostGameResult({ game_id: gameId, deck_id: deckId, player_id: playerId, place, kill_count: killCount });
        window.location.reload();
    }

    // TODO: When adding a game result, does Player actually matter?
    return (
        <Dialog open onClose={onClose}>
            <DialogTitle>Add Result</DialogTitle>
            <DialogContent>
                <Stack spacing={2} sx={{ mt: 1, minWidth: 300 }}>
                    <Autocomplete
                        options={players}
                        getOptionLabel={(p) => p.name}
                        onChange={(_, p) => setPlayerId(p?.id ?? null)}
                        renderInput={(params) => <TextField {...params} label="Player" size="small" />}
                    />
                    <Autocomplete
                        options={decks}
                        getOptionLabel={(d) => d.name}
                        onChange={(_, d) => setDeckId(d?.id ?? null)}
                        renderInput={(params) => <TextField {...params} label="Deck" size="small" />}
                    />
                    <TextField
                        label="Place"
                        type="number"
                        value={place}
                        onChange={(e) => setPlace(Number(e.target.value))}
                        size="small"
                    />
                    <TextField
                        label="Kills"
                        type="number"
                        value={killCount}
                        onChange={(e) => setKillCount(Number(e.target.value))}
                        size="small"
                    />
                </Stack>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleAdd} variant="contained" disabled={!playerId || !deckId}>Add</Button>
            </DialogActions>
        </Dialog>
    );
}

interface GameResultsGridProps {
    game: Game;
    pod: Pod;
    isManager: boolean;
    decks: Deck[];
    players: PlayerWithRole[];
}

function GameResultsGrid({ game, pod, isManager, decks, players }: GameResultsGridProps): ReactElement {
    const [editTarget, setEditTarget] = useState<GameResult | null>(null);
    const [removeTarget, setRemoveTarget] = useState<GameResult | null>(null);
    const [addOpen, setAddOpen] = useState(false);

    const baseColumns: Array<GridColDef> = [
        {
            field: "place",
            headerName: "Place",
            type: "number",
            minWidth: 80,
        },
        {
            field: "deck_name",
            headerName: "Deck",
            renderCell: (params) => (
                <Link to={`/player/${params.row.player_id}/deck/${params.row.deck_id}`}>{params.row.deck_name}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        {
            field: "commander_name",
            headerName: "Commander",
            flex: 1,
            renderCell: (params) => params.row.commander_name ?? "—",
        },
        {
            field: "kill_count",
            headerName: "Kills",
            type: "number",
            minWidth: 80,
        },
        {
            field: "points",
            headerName: "Points",
            type: "number",
            minWidth: 80,
        },
    ];

    const managerColumns: Array<GridColDef> = [
        {
            field: "_edit",
            headerName: "",
            width: 60,
            sortable: false,
            renderCell: (params) => (
                <IconButton size="small" onClick={() => setEditTarget(params.row as GameResult)}>
                    <EditIcon fontSize="small" />
                </IconButton>
            ),
        },
        {
            field: "_remove",
            headerName: "",
            width: 60,
            sortable: false,
            renderCell: (params) => (
                <IconButton size="small" color="error" onClick={() => setRemoveTarget(params.row as GameResult)}>
                    <DeleteIcon fontSize="small" />
                </IconButton>
            ),
        },
    ];

    const columns = isManager ? [...baseColumns, ...managerColumns] : baseColumns;

    return (
        <Box sx={{ width: "100%" }}>
            <Box sx={{ height: 355, width: "100%" }}>
                <DataGrid
                    rows={game.results}
                    columns={columns}
                    initialState={{
                        sorting: {
                            sortModel: [{ field: "place", sort: "asc" }],
                        },
                    }}
                />
            </Box>
            {isManager && (
                <Box sx={{ mt: 1 }}>
                    <Button onClick={() => setAddOpen(true)}>Add Result</Button>
                </Box>
            )}
            {editTarget && (
                <EditResultModal result={editTarget} decks={decks} onClose={() => setEditTarget(null)} />
            )}
            {removeTarget && (
                <RemoveResultDialog result={removeTarget} onClose={() => setRemoveTarget(null)} />
            )}
            {addOpen && (
                <AddResultModal
                    gameId={game.id}
                    podId={pod.id}
                    decks={decks}
                    players={players}
                    onClose={() => setAddOpen(false)}
                />
            )}
        </Box>
    );
}

interface DeleteGameButtonProps {
    gameId: number;
    podId: number;
}

function DeleteGameButton({ gameId, podId }: DeleteGameButtonProps): ReactElement {
    const navigate = useNavigate();
    const [open, setOpen] = useState(false);

    async function handleDelete() {
        await DeleteGame(gameId);
        navigate(`/pod/${podId}`);
    }

    return (
        <>
            <Button color="error" variant="outlined" onClick={() => setOpen(true)} sx={{ mt: 2 }}>
                Delete Game
            </Button>
            <Dialog open={open} onClose={() => setOpen(false)}>
                <DialogTitle>Delete game?</DialogTitle>
                <DialogContent>
                    <Typography>This action cannot be undone.</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpen(false)}>Cancel</Button>
                    <Button color="error" onClick={handleDelete}>Delete</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}

// TODO: Restructure file system
export default function GameView(): ReactElement {
    const { game, pod, players, decks } = useLoaderData() as GameLoaderData;
    const { user } = useAuth();
    const isManager = players.some((p) => p.id === user?.player_id && p.role === "manager");

    return (
        <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
            <h1>{pod.name} — Game #{game.id}</h1>
            <em>{new Date(game.created_at).toLocaleString()}</em>
            <GameDescription game={game} isManager={isManager} />
            <GameResultsGrid game={game} pod={pod} isManager={isManager} decks={decks} players={players} />
            {isManager && <DeleteGameButton gameId={game.id} podId={pod.id} />}
        </Box>
    );
}

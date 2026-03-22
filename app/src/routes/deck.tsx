import { ReactElement, useState } from "react";
import { Link, useLoaderData, useNavigate } from "react-router-dom";
import {
    Autocomplete,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    MenuItem,
    Select,
    SelectChangeEvent,
    Skeleton,
    Tab,
    Tabs,
    TextField,
    Typography,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";

import { useAuth } from "../auth";
import { AsyncComponentHelper } from "../common";
import {
    DeleteDeck,
    GetCommanders,
    GetFormats,
    GetGamesForDeck,
    PatchDeck,
} from "../http";
import { MatchesDisplay } from "../matches";
import { Record } from "../stats";
import { Commander, Deck, Format } from "../types";

// TODO: File system restructuring
export default function DeckView(): ReactElement {
    const deck = useLoaderData() as Deck;
    const { user } = useAuth();
    const [tab, setTab] = useState(0);

    const isOwner = user?.player_id === deck.player_id;

    // TODO This is ugly, move to helper function that uses ifs and early returns
    const commanderLabel = deck.commanders
        ? deck.commanders.partner_commander_name
            ? `${deck.commanders.commander_name} / ${deck.commanders.partner_commander_name}`
            : deck.commanders.commander_name
        : null;

    return (
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
            <Typography variant="h4" sx={{ mb: 0 }}>{deck.name}</Typography>
            {commanderLabel && (
                <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>{commanderLabel}</Typography>
            )}
            <Record record={deck.stats.record} />
            {deck.retired && (
                <Box sx={{ display: "flex", alignItems: "center", gap: 0.5, mt: 1 }}>
                    <DeleteIcon fontSize="small" />
                    <Typography variant="body2">Retired</Typography>
                </Box>
            )}
            {/* TODO: Common tabs component */}
            <Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 2, mt: 2 }}>
                <Tab label="Overview" />
                <Tab label="Games" />
                {isOwner && <Tab label="Settings" />}
            </Tabs>
            {tab === 0 && <DeckOverviewTab deck={deck} />}
            {tab === 1 && <DeckGamesTab deck={deck} />}
            {tab === 2 && isOwner && <DeckSettingsTab deck={deck} />}
        </Box>
    );
}

interface DeckTabProps {
    deck: Deck;
}

function DeckOverviewTab({ deck }: DeckTabProps): ReactElement {
    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
            <Box sx={{ display: "flex", flexDirection: "row", gap: 4 }}>
                <span><strong>Games Played:</strong> {deck.stats.games}</span>
                <span><strong>Total Kills:</strong> {deck.stats.kills}</span>
                <span><strong>Total Points:</strong> {deck.stats.points}</span>
            </Box>
            <Typography variant="body1"><strong>Format:</strong> {deck.format_name}</Typography>
            <Typography variant="body1">
                <strong>Owner:</strong> <Link to={`/player/${deck.player_id}`}>{deck.player_name}</Link>
            </Typography>
            <Typography variant="body2" color="text.secondary">
                Created: {new Date(deck.created_at).toLocaleString()}
            </Typography>
        </Box>
    );
}

function DeckGamesTab({ deck }: DeckTabProps): ReactElement {
    const { data, loading, error } = AsyncComponentHelper(GetGamesForDeck(deck.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} />;
    }
    if (error) {
        return <Typography color="error">Error loading games: {error.message}</Typography>;
    }

    return (
        <Box sx={{ height: 500, width: "100%" }}>
            <MatchesDisplay games={data} targetCommander={deck.commanders?.commander_name} />
        </Box>
    );
}

// TODO: Need to rebuild context on working with the frontend after all of this
function DeckSettingsTab({ deck }: DeckTabProps): ReactElement {
    const navigate = useNavigate();

    const { data: formats, loading: formatsLoading } = AsyncComponentHelper(GetFormats());
    const { data: commanders, loading: commandersLoading } = AsyncComponentHelper(GetCommanders());

    const [name, setName] = useState(deck.name);
    const [nameError, setNameError] = useState<string | null>(null);

    const [formatId, setFormatId] = useState(deck.format_id);
    const [formatError, setFormatError] = useState<string | null>(null);

    const [commanderId, setCommanderId] = useState<number | null>(deck.commanders?.commander_id ?? null);
    const [partnerCommanderId, setPartnerCommanderId] = useState<number | null>(
        deck.commanders?.partner_commander_id ?? null
    );
    const [commanderError, setCommanderError] = useState<string | null>(null);

    const [retireConfirm, setRetireConfirm] = useState(false);
    const [deleteConfirm, setDeleteConfirm] = useState(false);
    const [actionError, setActionError] = useState<string | null>(null);

    const handleSaveName = async () => {
        setNameError(null);
        try {
            await PatchDeck(deck.id, { name });
            window.location.reload();
        } catch {
            setNameError("Failed to update name.");
        }
    };

    const handleSaveFormat = async () => {
        setFormatError(null);
        try {
            await PatchDeck(deck.id, { format_id: formatId });
            window.location.reload();
        } catch {
            setFormatError("Failed to update format.");
        }
    };

    const handleSaveCommanders = async () => {
        setCommanderError(null);
        if (commanderId === null) {
            setCommanderError("A primary commander is required.");
            return;
        }
        try {
            await PatchDeck(deck.id, {
                commander_id: commanderId,
                partner_commander_id: partnerCommanderId ?? undefined,
            });
            window.location.reload();
        } catch {
            setCommanderError("Failed to update commanders.");
        }
    };

    // TODO: Check retirement behavior - discuss where retired decks should and should not show. E.g., should show on player->decks view, should not show on pod->decks
    // TODO: Should be filterable if we're showing decks in a grid
    const handleRetire = async () => {
        setRetireConfirm(false);
        setActionError(null);
        try {
            await PatchDeck(deck.id, { retired: !deck.retired });
            if (deck.retired) {
                window.location.reload();
            } else {
                navigate(`/player/${deck.player_id}`);
            }
        } catch {
            setActionError("Failed to update retired status.");
        }
    };

    const handleDelete = async () => {
        setDeleteConfirm(false);
        setActionError(null);
        try {
            await DeleteDeck(deck.id);
            navigate(`/player/${deck.player_id}`);
        } catch {
            setActionError("Failed to delete deck.");
        }
    };

    const currentCommander = commanders?.find((c: Commander) => c.id === commanderId) ?? null;
    const currentPartner = commanders?.find((c: Commander) => c.id === partnerCommanderId) ?? null;

    return (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3, maxWidth: 500 }}>
            {/* Edit Name */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Deck Name</Typography>
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

            {/* Edit Format */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Format</Typography>
                {formatsLoading ? (
                    <Skeleton variant="rectangular" height={40} width={200} />
                ) : (
                    <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
                        <Select
                            value={String(formatId)}
                            onChange={(e: SelectChangeEvent) => setFormatId(Number(e.target.value))}
                            size="small"
                        >
                            {(formats ?? []).map((f: Format) => (
                                <MenuItem key={f.id} value={String(f.id)}>{f.name}</MenuItem>
                            ))}
                        </Select>
                        <Button variant="contained" onClick={handleSaveFormat}>Save</Button>
                    </Box>
                )}
                {formatError && <Typography color="error" variant="body2">{formatError}</Typography>}
            </Box>

            {/* Edit Commanders */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Commanders</Typography>
                {commandersLoading ? (
                    <Skeleton variant="rectangular" height={80} width={300} />
                ) : (
                    <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                        <Autocomplete
                            options={commanders ?? []}
                            getOptionLabel={(c: Commander) => c.name}
                            getOptionKey={(c: Commander) => c.id}
                            value={currentCommander}
                            onChange={(_, value) => setCommanderId(value?.id ?? null)}
                            renderInput={(params) => <TextField {...params} label="Commander" size="small" />}
                            sx={{ width: 300 }}
                        />
                        <Autocomplete
                            options={commanders ?? []}
                            getOptionLabel={(c: Commander) => c.name}
                            getOptionKey={(c: Commander) => c.id}
                            value={currentPartner}
                            onChange={(_, value) => setPartnerCommanderId(value?.id ?? null)}
                            renderInput={(params) => <TextField {...params} label="Partner (optional)" size="small" />}
                            sx={{ width: 300 }}
                        />
                        <Button variant="contained" onClick={handleSaveCommanders} sx={{ width: "fit-content" }}>
                            Save
                        </Button>
                    </Box>
                )}
                {commanderError && <Typography color="error" variant="body2">{commanderError}</Typography>}
            </Box>

            {/* Retire / Un-retire */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">{deck.retired ? "Un-retire Deck" : "Retire Deck"}</Typography>
                <Box>
                    <Button
                        variant="outlined"
                        color={deck.retired ? "primary" : "warning"}
                        onClick={() => setRetireConfirm(true)}
                    >
                        {deck.retired ? "Un-retire" : "Retire"}
                    </Button>
                </Box>
            </Box>

            {/* Delete */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Delete Deck</Typography>
                <Box>
                    <Button variant="outlined" color="error" onClick={() => setDeleteConfirm(true)}>
                        Delete
                    </Button>
                </Box>
            </Box>

            {actionError && <Typography color="error" variant="body2">{actionError}</Typography>}

            <Dialog open={retireConfirm} onClose={() => setRetireConfirm(false)}>
                <DialogTitle>{deck.retired ? "Un-retire deck?" : "Retire deck?"}</DialogTitle>
                <DialogContent>
                    <Typography>
                        {deck.retired
                            ? "Mark this deck as active again?"
                            : "Mark this deck as retired? It will no longer appear in active deck lists."}
                    </Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setRetireConfirm(false)}>Cancel</Button>
                    <Button color={deck.retired ? "primary" : "warning"} onClick={handleRetire}>
                        {deck.retired ? "Un-retire" : "Retire"}
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog open={deleteConfirm} onClose={() => setDeleteConfirm(false)}>
                <DialogTitle>Delete deck?</DialogTitle>
                <DialogContent>
                    <Typography>This will permanently delete <strong>{deck.name}</strong>. This cannot be undone.</Typography>
                    <Typography>If you want to preserve overall record data, retire this deck instead.</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteConfirm(false)}>Cancel</Button>
                    <Button color="error" onClick={handleDelete}>Delete</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
}

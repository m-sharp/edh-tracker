import { ReactElement, useState } from "react";
import { useNavigate } from "react-router-dom";
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
    TextField,
    Typography,
} from "@mui/material";

import { AsyncComponentHelper } from "../../components/common";
import { TooltipIcon } from "../../components/TooltipIcon";
import {
    DeleteDeck,
    GetCommanders,
    GetFormats,
    PatchDeck,
} from "../../http";
import { Commander, Deck, Format } from "../../types";

interface DeckSettingsTabProps {
    deck: Deck;
}

export default function DeckSettingsTab({ deck }: DeckSettingsTabProps): ReactElement {
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
                    <Button variant="contained" onClick={handleSaveName}>Save Name</Button>
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
                        <Button variant="contained" onClick={handleSaveFormat}>Save Format</Button>
                    </Box>
                )}
                {formatError && <Typography color="error" variant="body2">{formatError}</Typography>}
            </Box>

            {/* Edit Commanders */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                    <Typography variant="h6">Commanders</Typography>
                    <TooltipIcon title="This is for changing an existing deck's commander. To add a new deck, use the Add Deck button instead." />
                </Box>
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
                            fullWidth
                        />
                        <Autocomplete
                            options={commanders ?? []}
                            getOptionLabel={(c: Commander) => c.name}
                            getOptionKey={(c: Commander) => c.id}
                            value={currentPartner}
                            onChange={(_, value) => setPartnerCommanderId(value?.id ?? null)}
                            renderInput={(params) => <TextField {...params} label="Partner (optional)" size="small" />}
                            fullWidth
                        />
                        <Button variant="contained" onClick={handleSaveCommanders} sx={{ width: "fit-content" }}>
                            Save Commanders
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
                        sx={{ minHeight: 44 }}
                    >
                        {deck.retired ? "Un-retire" : "Retire"}
                    </Button>
                </Box>
            </Box>

            {/* Delete */}
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                <Typography variant="h6">Delete Deck</Typography>
                <Box>
                    <Button
                        variant="outlined"
                        color="error"
                        onClick={() => setDeleteConfirm(true)}
                        sx={{ minHeight: 44 }}
                    >
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

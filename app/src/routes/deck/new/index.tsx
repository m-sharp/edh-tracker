import { ReactElement, useState } from "react";
import { Link, useLoaderData, useNavigate } from "react-router-dom";
import {
    Autocomplete,
    Box,
    Button,
    CircularProgress,
    Container,
    createFilterOptions,
    MenuItem,
    Paper,
    Select,
    SelectChangeEvent,
    TextField,
    Typography,
} from "@mui/material";
import { useAuth } from "../../../auth";
import { GetCommanders, GetFormats, PostCommander, PostDeck } from "../../../http";
import { Commander, Format, NewDeckData, NewDeckRequest } from "../../../types";

const filter = createFilterOptions<Commander>();

export async function newDeckLoader(): Promise<NewDeckData> {
    const [formats, commanders] = await Promise.all([GetFormats(), GetCommanders()]);
    return { formats, commanders };
}

export default function NewDeckView(): ReactElement {
    const { formats, commanders: initialCommanders } = useLoaderData() as NewDeckData;
    const { user } = useAuth();
    const navigate = useNavigate();

    const [name, setName] = useState("");
    const [formatId, setFormatId] = useState<number | "">("");
    const [commanders, setCommanders] = useState<Commander[]>(initialCommanders);
    const [commanderId, setCommanderId] = useState<number | null>(null);
    const [commanderInput, setCommanderInput] = useState("");
    const [partnerCommanderId, setPartnerCommanderId] = useState<number | null>(null);
    const [partnerInput, setPartnerInput] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const selectedFormat = formats.find((f) => f.id === formatId);
    const isCommander = selectedFormat?.name === "Commander";

    const canSubmit = name.trim() !== ""
        && formatId !== ""
        && (!isCommander || commanderId !== null);

    const handleCommanderSelect = async (
        value: Commander | string | null,
        inputValue: string,
        setId: (id: number | null) => void,
    ) => {
        if (value === null) {
            setId(null);
            return;
        }
        if (typeof value === "string" || (typeof value === "object" && value.id === -1)) {
            // Create new commander
            const newName = typeof value === "string" ? value : inputValue.replace(/^Create "/, "").replace(/"$/, "");
            try {
                const res = await PostCommander(newName);
                if (!res.ok) throw new Error("Failed");
                const { id } = await res.json();
                setCommanders((prev) => [...prev, { id, name: newName }]);
                setId(id);
            } catch {
                setError("Failed to create commander. Try again.");
            }
            return;
        }
        setId(value.id);
    };

    const handleSubmit = async () => {
        setError(null);
        setSubmitting(true);
        try {
            const body: NewDeckRequest = {
                name: name.trim(),
                format_id: formatId as number,
            };
            if (isCommander && commanderId !== null) {
                body.commander_id = commanderId;
                if (partnerCommanderId !== null) {
                    body.partner_commander_id = partnerCommanderId;
                }
            }
            const { id } = await PostDeck(body);
            navigate(`/player/${user!.player_id}/deck/${id}`);
        } catch {
            setError("Failed to create deck. Try again.");
            setSubmitting(false);
        }
    };

    const renderCommanderAutocomplete = (
        label: string,
        currentId: number | null,
        inputValue: string,
        setInputValue: (v: string) => void,
        setId: (id: number | null) => void,
        required: boolean,
    ) => {
        const currentValue = commanders.find((c) => c.id === currentId) ?? null;
        return (
            <Autocomplete
                options={commanders}
                freeSolo
                value={currentValue}
                inputValue={inputValue}
                onInputChange={(_, v) => setInputValue(v)}
                getOptionLabel={(opt) => typeof opt === "string" ? opt : opt.name}
                filterOptions={(options, params) => {
                    const filtered = filter(options, params);
                    const { inputValue: iv } = params;
                    const isExisting = options.some((opt) => iv === opt.name);
                    if (iv !== "" && !isExisting) {
                        filtered.push({ id: -1, name: `Create "${iv}"` } as Commander);
                    }
                    return filtered;
                }}
                onChange={(_, value) => handleCommanderSelect(value, inputValue, setId)}
                renderInput={(params) => (
                    <TextField {...params} label={label} required={required} fullWidth />
                )}
                fullWidth
                disabled={submitting}
            />
        );
    };

    return (
        <Container maxWidth="sm">
            <Paper elevation={2} sx={{ p: 3 }}>
                <Typography variant="h4" sx={{ mb: 3 }}>New Deck</Typography>

                <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
                    <TextField
                        label="Deck Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        fullWidth
                        disabled={submitting}
                    />

                    <Select
                        value={formatId === "" ? "" : String(formatId)}
                        onChange={(e: SelectChangeEvent) => {
                            const val = Number(e.target.value);
                            setFormatId(val);
                            // Reset commander fields when format changes
                            if (formats.find((f) => f.id === val)?.name !== "Commander") {
                                setCommanderId(null);
                                setPartnerCommanderId(null);
                                setCommanderInput("");
                                setPartnerInput("");
                            }
                        }}
                        displayEmpty
                        fullWidth
                        disabled={submitting}
                    >
                        <MenuItem value="" disabled>Select Format</MenuItem>
                        {formats.map((f: Format) => (
                            <MenuItem key={f.id} value={String(f.id)}>{f.name}</MenuItem>
                        ))}
                    </Select>

                    {isCommander && (
                        <>
                            {renderCommanderAutocomplete(
                                "Commander",
                                commanderId,
                                commanderInput,
                                setCommanderInput,
                                setCommanderId,
                                true,
                            )}
                            {renderCommanderAutocomplete(
                                "Partner Commander (optional)",
                                partnerCommanderId,
                                partnerInput,
                                setPartnerInput,
                                setPartnerCommanderId,
                                false,
                            )}
                        </>
                    )}

                    {error && <Typography color="error" variant="body2">{error}</Typography>}

                    <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 1, flexDirection: { xs: "column-reverse", sm: "row" } }}>
                        <Button
                            variant="outlined"
                            component={Link}
                            to={`/player/${user?.player_id ?? ""}`}
                            disabled={submitting}
                        >
                            Discard
                        </Button>
                        <Button
                            variant="contained"
                            disabled={!canSubmit || submitting}
                            onClick={handleSubmit}
                        >
                            {submitting ? <CircularProgress size={20} /> : "Create Deck"}
                        </Button>
                    </Box>
                </Box>
            </Paper>
        </Container>
    );
}

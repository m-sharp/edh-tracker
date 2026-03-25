import { ReactElement, useState } from "react";
import { Form, redirect, useLoaderData, useParams, useSubmit } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import {
    Autocomplete,
    Box,
    Button,
    Card,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    SelectChangeEvent,
    TextField,
    Typography,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import CloseIcon from "@mui/icons-material/Close";
import PublishIcon from "@mui/icons-material/Publish";
import { GetAllDecksForPod, GetFormats, PostGame } from "../../http";
import { Deck, Format, NewGame, NewGameData } from "../../types";
import { TooltipIconButton } from "../../components/TooltipIcon";

export async function newGameLoader({ params }: LoaderFunctionArgs): Promise<NewGameData> {
    const podId = Number(params.podId);
    const [decks, formats] = await Promise.all([
        GetAllDecksForPod(podId),
        GetFormats(),
    ]);
    return { decks, formats };
}

export async function createGame({ request, params }: { request: Request; params: Record<string, string | undefined> }): Promise<Response> {
    const body = await request.json();

    const resp = await PostGame(body);
    if (!resp.ok) {
        throw new Error("Failed to create new game record: received " + resp.status + " " + resp.statusText);
    }

    const { id } = await resp.json();
    return redirect(`/pod/${params.podId}/game/${id}`);
}

interface CardState {
    key: number;
    deckId: number | null;
    place: string;
    kills: string;
}

const initialCards: CardState[] = [
    { key: 0, deckId: null, place: "", kills: "" },
    { key: 1, deckId: null, place: "", kills: "" },
];

export default function View(): ReactElement {
    const data = useLoaderData() as NewGameData;
    const { podId } = useParams();
    const submit = useSubmit();

    const [formatID, setFormatID] = useState<number>(0);
    const [cards, setCards] = useState<CardState[]>(initialCards);
    const [nextKey, setNextKey] = useState(2);
    const [desc, setDesc] = useState("");
    const [showDescription, setShowDescription] = useState(false);

    const filteredDecks = formatID !== 0
        ? data.decks.filter((d) => d.format_id === formatID)
        : data.decks;

    const addCard = () => {
        setCards((prev) => [...prev, { key: nextKey, deckId: null, place: "", kills: "" }]);
        setNextKey((k) => k + 1);
    };

    const removeCard = (key: number) => {
        setCards((prev) => prev.filter((c) => c.key !== key));
    };

    const updateCard = (key: number, field: keyof Omit<CardState, "key">, value: any) => {
        setCards((prev) => prev.map((c) => c.key === key ? { ...c, [field]: value } : c));
    };

    const isSubmittable =
        formatID !== 0 &&
        cards.length >= 2 &&
        cards.every((c) => c.deckId !== null && c.place !== "" && c.kills !== "");

    const handleSubmit = (e: React.MouseEvent) => {
        e.preventDefault();
        const result: NewGame = {
            description: desc,
            format_id: formatID,
            pod_id: Number(podId),
            results: cards.map((c) => ({
                deck_id: c.deckId!,
                place: Number(c.place),
                kill_count: Number(c.kills),
            })),
        };
        submit(result as { [key: string]: any }, { method: "post", encType: "application/json" });
    };

    return (
        <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", width: "100%", pb: 4 }}>
            <Typography variant="h4">Add New Game</Typography>
            <Form>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 2, width: "100%", maxWidth: 600 }}>
                    <FormControl fullWidth required>
                        <InputLabel id="format-label">Format</InputLabel>
                        <Select
                            labelId="format-label"
                            label="Format"
                            value={formatID === 0 ? "" : String(formatID)}
                            onChange={(event: SelectChangeEvent) => {
                                setFormatID(Number(event.target.value));
                            }}
                        >
                            {data.formats.map((format: Format) => (
                                <MenuItem key={format.id} value={format.id}>{format.name}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>

                    <Button
                        variant="text"
                        size="small"
                        onClick={() => setShowDescription((prev) => !prev)}
                        sx={{ alignSelf: "flex-start" }}
                    >
                        {showDescription ? "- Hide description" : "+ Add description"}
                    </Button>
                    {showDescription && (
                        <TextField
                            multiline
                            fullWidth
                            label="Description"
                            placeholder="Add a game description!"
                            value={desc}
                            onChange={(e) => setDesc(e.target.value)}
                        />
                    )}

                    {cards.map((card) => (
                        <Card key={card.key} variant="outlined" sx={{ width: "100%", p: 2 }}>
                            <Box sx={{ display: "flex", alignItems: "flex-start", gap: 1 }}>
                                {/* Left: fields column */}
                                <Box sx={{ flex: 1, display: "flex", flexDirection: "column", gap: 1.5 }}>
                                    <Autocomplete
                                        fullWidth
                                        options={filteredDecks}
                                        getOptionLabel={(deck: Deck) => `${deck.name} (${deck.player_name})`}
                                        getOptionKey={(deck: Deck) => deck.id}
                                        value={filteredDecks.find((d) => d.id === card.deckId) ?? null}
                                        onChange={(_, value) => updateCard(card.key, "deckId", value?.id ?? null)}
                                        noOptionsText="No decks available for this format."
                                        renderInput={(params) => <TextField {...params} label="Deck" required />}
                                    />
                                    <Box sx={{ display: "flex", gap: 2 }}>
                                        <TextField
                                            type="number"
                                            label="Place"
                                            required
                                            sx={{ flex: 1 }}
                                            value={card.place}
                                            onChange={(e) => updateCard(card.key, "place", e.target.value)}
                                            inputProps={{ min: 1, max: cards.length }}
                                        />
                                        <TextField
                                            type="number"
                                            label="Kills"
                                            required
                                            sx={{ flex: 1 }}
                                            value={card.kills}
                                            onChange={(e) => updateCard(card.key, "kills", e.target.value)}
                                            inputProps={{ min: 0, max: cards.length }}
                                        />
                                    </Box>
                                </Box>
                                {/* Right: remove button */}
                                <TooltipIconButton
                                    title={cards.length <= 2 ? "Minimum 2 entries required" : "Remove"}
                                    icon={<CloseIcon />}
                                    onClick={() => removeCard(card.key)}
                                    color="error"
                                    size="small"
                                    disabled={cards.length <= 2}
                                    sx={{ minHeight: 44, minWidth: 44 }}
                                />
                            </Box>
                        </Card>
                    ))}

                    <Button
                        variant="outlined"
                        startIcon={<AddIcon />}
                        sx={{ alignSelf: "flex-start", mt: 1, minHeight: 44 }}
                        onClick={addCard}
                    >
                        Add Deck
                    </Button>

                    <Button
                        variant="contained"
                        size="large"
                        startIcon={<PublishIcon />}
                        sx={{ mt: 2, minHeight: 48 }}
                        fullWidth
                        onClick={handleSubmit}
                        disabled={!isSubmittable}
                    >
                        Submit Game
                    </Button>
                </Box>
            </Form>
        </Box>
    );
}

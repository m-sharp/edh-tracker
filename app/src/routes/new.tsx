import { Dispatch, ReactElement, SetStateAction, SyntheticEvent, useState } from "react";
import { Form, useLoaderData, useSubmit } from "react-router-dom";
import {
    Autocomplete,
    Box,
    Button,
    MenuItem,
    Select,
    SelectChangeEvent,
    TextField
} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import PublishIcon from '@mui/icons-material/Publish';

import { PostGame } from "../http";
import { Deck, Format, NewGame, NewGameData, NewGameResult, Player } from "../types";

interface CreateActionProps {
    request: Request;
}

// ToDo: Validation
export async function createGame({request}: CreateActionProps): Promise<null> {
    const body = await request.json()

    const resp = await PostGame(body);
    if ( !resp.ok ) {
        throw new Error("Failed to create new game record: received " + resp.status + " " + resp.statusText);
    }

    // ToDo: Doesn't trigger any reload, probably need to return an object back?
    // ToDo: Toast alerts?
    return null;
}

interface ResultsMap {
    [key: number]: NewGameResult;
}

export default function View(): ReactElement {
    const newGameInfo = useLoaderData() as NewGameData;
    const submit = useSubmit();

    const [formatID, setFormatID] = useState<number>(0);
    const [numPlayers, setNumPlayers] = useState<number>(2);
    const initResults: ResultsMap = {};
    for (let i=0; i<numPlayers; i++) {
        initResults[i] = {
            deck_id: 0,
            kill_count: 0,
            place: 0,
            player_id: 0,
        };
    }

    const [results, setResults] = useState<ResultsMap>(initResults);

    const filteredDecks = formatID !== 0
        ? newGameInfo.decks.filter((d) => d.format_id === formatID)
        : newGameInfo.decks;

    let inputs: Array<ReactElement> = [];
    for (let i=0; i<numPlayers; i++) {
        inputs.push(<GameInput
            key={i}
            resultKey={i}
            players={newGameInfo.players}
            decks={filteredDecks}
            numOfPlayers={numPlayers}
            setResults={setResults}
        />)
    }

    const [desc, setDesc] = useState<string>("");

    // TODO: Better styling
    // TODO: Honor player_id on backend
    return (
        <Box id="newGameForm" sx={{display: "flex", flexDirection: "column", alignItems: "center", width: "100%"}}>
            <h1>Add New Game</h1>
            <Form method="post">
                <Box sx={{display: "flex", flexDirection: "column", alignItems: "center", gap: "1em"}}>
                    <Select
                        label={"Format"}
                        id={"format-select"}
                        value={formatID === 0 ? "" : String(formatID)}
                        displayEmpty
                        onChange={(event: SelectChangeEvent) => {
                            setFormatID(Number(event.target.value));
                        }}
                    >
                        <MenuItem value="" disabled>Select Format</MenuItem>
                        {newGameInfo.formats.map((format: Format) => (
                            <MenuItem key={format.id} value={format.id}>{format.name}</MenuItem>
                        ))}
                    </Select>
                    <TextField
                        multiline
                        label={"Game Description"}
                        placeholder={"Add a game description!"}
                        name={"description"}
                        sx={{width: "50%"}}
                        value={desc}
                        onChange={(e) => {
                            setDesc(e.target.value);
                        }}
                    />
                    {inputs.map((input) => input)}
                    <Button
                        variant={"outlined"}
                        type={"button"}
                        size={"medium"}
                        startIcon={<AddIcon />}
                        sx={{marginTop: 2}}
                        onClick={(e) => {
                            e.preventDefault();
                            setNumPlayers((prev) => { return prev + 1})
                        }}
                    >
                        Add Player
                    </Button>
                </Box>
                <Button
                    variant={"contained"}
                    type={"submit"}
                    size={"large"}
                    startIcon={<PublishIcon />}
                    sx={{marginTop: 2}}
                    onClick={(e) => {
                        e.preventDefault();
                        const result: NewGame = {
                            description: desc,
                            format_id: formatID,
                            pod_id: 0,
                            results: Object.values(results),
                        }
                        submit(result as {[key: string]: any}, {method: "post", encType: "application/json"});
                    }}
                >
                    Submit
                </Button>
            </Form>
        </Box>
    );
}

interface GameInputProps {
    resultKey: number;
    players: Array<Player>;
    decks: Array<Deck>;
    numOfPlayers: number;
    setResults: Dispatch<SetStateAction<ResultsMap>>;
}

type OnChangeFn = (setter: Dispatch<SetStateAction<number>>, val: number) => void;

function GameInput({resultKey, players, decks, numOfPlayers, setResults}: GameInputProps): ReactElement {
    const [playerID, setPlayerID] = useState<number>(0);
    const [deckID, setDeckID] = useState<number>(0);
    const [place, setPlace] = useState<number>(0);
    const [kills, setKills] = useState<number>(0);

    const onChange: OnChangeFn = (setter, val) => {
        // Call the individual setter
        setter(val);

        // Update the overall result
        const newResult = {
            deck_id: deckID,
            player_id: playerID,
            place: place,
            kill_count: kills,
        }
        setResults((prev) => {
            return {
                ...prev,
                [resultKey]: newResult,
            }
        });
    }

    return (
        <Box sx={{display: "flex", flexDirection: "column", alignItems: "center", width: "100%"}}>
            <Box sx={{display: "flex", flexDirection: "row", justifyContent: "space-evenly", width: "100%", gap: "1em"}}>
                <Select
                    label={"Player"}
                    id={`player-select-${resultKey}`}
                    onChange={(event: SelectChangeEvent) => {
                        onChange(setPlayerID, Number(event.target.value))
                    }}
                >
                    {players.map((player) => <MenuItem key={player.id} value={player.id}>{player.name}</MenuItem>)}
                </Select>
                <Autocomplete
                    autoComplete
                    id={`deck-select-${resultKey}`}
                    sx={{width: 300}}
                    disablePortal
                    options={decks}
                    getOptionLabel={(deck: Deck) => `${deck.name}${deck.commanders ? ` (${deck.commanders.commander_name})` : ""} — ${deck.player_name}`}
                    getOptionKey={(deck: Deck) => deck.id}
                    onChange={(event: SyntheticEvent, value: Deck | null, _reason: string) => {
                        if (value === null) {
                            return;
                        }

                        onChange(setDeckID, value.id);
                    }}
                    renderInput={(params) => <TextField {...params} label={`Deck ${resultKey}`} required />}
                />
                <TextField
                    required
                    name={`deck${resultKey}Place`}
                    type={"number"}
                    helperText={"Finishing Place"}
                    inputProps={{min: 1, max: numOfPlayers}}
                    onChange={(event) => {
                        onChange(setPlace, Number(event.target.value));
                    }}
                />
                <TextField
                    required
                    name={`deck${resultKey}Kills`}
                    type={"number"}
                    helperText={"Kill Points"}
                    inputProps={{min: 0, max: numOfPlayers}}
                    onChange={(event) => {
                        onChange(setKills, Number(event.target.value));
                    }}
                />
            </Box>
        </Box>
    );
}

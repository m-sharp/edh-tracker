import { Dispatch, ReactElement, SetStateAction, SyntheticEvent, useState } from "react";
import { Form, useLoaderData, useSubmit } from "react-router-dom";
import { Autocomplete, Box, Button, TextField } from "@mui/material";
import PublishIcon from '@mui/icons-material/Publish';

import { PostGame } from "../http";
import { Deck, NewGame, NewGameResult } from "../types";

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
    [key: string]: NewGameResult;
}

export default function View(): ReactElement {
    const decks = useLoaderData() as Array<Deck>;
    const submit = useSubmit();

    const [results, setResults] = useState<ResultsMap>({});
    const [desc, setDesc] = useState<string>("");

    return (
        <Box id="newGameForm" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <h1>Add New Game</h1>
            <Form method="post">
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
                <DeckInputs decks={decks} keyName={"One"} results={results} setResults={setResults} />
                <DeckInputs decks={decks} keyName={"Two"} results={results} setResults={setResults} />
                {/*ToDo: Button to add more players*/}
                <DeckInputs decks={decks} keyName={"Three"} results={results} setResults={setResults} />
                <DeckInputs decks={decks} keyName={"Four"} results={results} setResults={setResults} />
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

interface DeckInputProps {
    decks: Array<Deck>;
    // ToDo: Generate a unique key on DeckInputs render?
    keyName: string;
    results: ResultsMap;
    setResults: Dispatch<SetStateAction<ResultsMap>>;
}

function DeckInputs({decks, keyName, results, setResults}: DeckInputProps): ReactElement {
    return (
        <Box sx={{display: "flex", flexDirection: "column", alignItems: "center", width: "100%"}}>
            <h2>Deck {keyName}</h2>
            <Box sx={{display: "flex", flexDirection: "row", justifyContent: "space-evenly", width: "100%"}}>
                <Autocomplete
                    autoComplete
                    sx={{width: 300}}
                    disablePortal
                    options={decks}
                    getOptionLabel={(deck: Deck) => `${deck.commander} (${deck.player_name})`}
                    getOptionKey={(deck: Deck) => deck.id}
                    onChange={(event: SyntheticEvent, value: Deck | null, _reason: string) => {
                        if (value !== null) {
                            // ToDo: Cleaner approach for all three onChanges?
                            if (!(keyName in results)) {
                                const newResult: NewGameResult = {deck_id: value.id, place: 0, kill_count: 0};
                                setResults({...results, [keyName]: newResult});
                            } else {
                                results[keyName].deck_id = value.id;
                                setResults(results);
                            }
                        }
                    }}
                    renderInput={(params) => <TextField {...params} label="Deck" required />}
                />
                <TextField
                    required
                    name={`deck${keyName}Place`}
                    type={"number"}
                    helperText={"Finishing Place"}
                    inputProps={{min: 1, max: 4}}
                    onChange={(event) => {
                        if (!(keyName in results)) {
                            const newResult: NewGameResult = {deck_id: 0, place: Number(event.target.value), kill_count: 0};
                            setResults({...results, [keyName]: newResult});
                        } else {
                            results[keyName].place = Number(event.target.value);
                            setResults(results);
                        }
                    }}
                />
                <TextField
                    required
                    name={`deck${keyName}Kills`}
                    type={"number"}
                    helperText={"Kill Points"}
                    inputProps={{min: 0, max: 4}}
                    onChange={(event) => {
                        if (!(keyName in results)) {
                            const newResult: NewGameResult = {deck_id: 0, place: 0, kill_count: Number(event.target.value)};
                            setResults({...results, [keyName]: newResult});
                        } else {
                            results[keyName].kill_count = Number(event.target.value);
                            setResults(results);
                        }
                    }}
                />
            </Box>
        </Box>
    );
}

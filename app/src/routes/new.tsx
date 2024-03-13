import { ReactElement, SyntheticEvent, useState } from "react";
import { Form, useLoaderData } from "react-router-dom";
import { Autocomplete, Box, Button, TextField } from "@mui/material";
import PublishIcon from '@mui/icons-material/Publish';

import { PostGame } from "../http";
import { Deck } from "../types";

interface CreateActionProps {
    request: Request;
}

// ToDo: Validation on formData.get()s
export async function createGame({request}: CreateActionProps): Promise<null> {
    const formData = await request.formData();
    const resp = await PostGame({
        "description": formData.get("description") as string,
        "results": [
            {
                "deck_id": parseInt(formData.get("deckOneId") as string, 10),
                "place": parseInt(formData.get("deckOnePlace") as string, 10),
                "kill_count": parseInt(formData.get("deckOneKills") as string, 10)
            },
            {
                "deck_id": parseInt(formData.get("deckTwoId") as string, 10),
                "place": parseInt(formData.get("deckTwoPlace") as string, 10),
                "kill_count": parseInt(formData.get("deckTwoKills") as string, 10)
            },
            {
                "deck_id": parseInt(formData.get("deckThreeId") as string, 10),
                "place": parseInt(formData.get("deckThreePlace") as string, 10),
                "kill_count": parseInt(formData.get("deckThreeKills") as string, 10)
            },
            {
                "deck_id": parseInt(formData.get("deckFourId") as string, 10),
                "place": parseInt(formData.get("deckFourPlace") as string, 10),
                "kill_count": parseInt(formData.get("deckFourKills") as string, 10)
            }
        ]
    });

    if ( !resp.ok ) {
        throw new Error("Failed to create new game record: received " + resp.status + " " + resp.statusText);
    }

    // ToDo: Doesn't trigger any reload, probably need to return an object back?
    // ToDo: Toast alerts?
    return null;
}

export default function View(): ReactElement {
    const decks = useLoaderData() as Array<Deck>;

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
                />
                <DeckInputs decks={decks} keyName={"One"} />
                <DeckInputs decks={decks} keyName={"Two"} />
                <DeckInputs decks={decks} keyName={"Three"} />
                <DeckInputs decks={decks} keyName={"Four"} />
                <Button
                    variant={"contained"}
                    type={"submit"}
                    size={"large"}
                    startIcon={<PublishIcon />}
                    sx={{marginTop: 2}}
                >
                    Submit
                </Button>
            </Form>
        </Box>
    );
}

interface DeckInputProps {
    decks: Array<Deck>;
    keyName: string;
}

function DeckInputs({decks, keyName}: DeckInputProps): ReactElement {
    const [deckValue, setValue] = useState<number>(0);

    return (
        <Box sx={{display: "flex", flexDirection: "column", alignItems: "center", width: "100%"}}>
            <h2>Deck {keyName}</h2>
            <Box sx={{display: "flex", flexDirection: "row", justifyContent: "space-evenly", width: "100%"}}>
                <Autocomplete
                    autoComplete
                    sx={{width: 300}}
                    disablePortal
                    options={decks}
                    getOptionLabel={(deck: Deck) => deck.commander}
                    getOptionKey={(deck: Deck) => deck.id}
                    onChange={(event: SyntheticEvent, value: Deck | null, _reason: string) => {
                        if (value !== null) {
                            setValue(value.id);
                        }
                    }}
                    renderInput={(params) => <TextField {...params} label="Deck" required />}
                />
                {/*Hold the Deck ID value in a hidden input*/}
                <input name={`deck${keyName}Id`} value={deckValue} hidden={true} style={{display: "none"}} />
                <TextField
                    required
                    name={`deck${keyName}Place`}
                    type={"number"}
                    helperText={"Finishing Place"}
                    inputProps={{min: 1, max: 4}}
                />
                <TextField
                    required
                    name={`deck${keyName}Kills`}
                    type={"number"}
                    helperText={"Kill Points"}
                    inputProps={{min: 0, max: 4}}
                />
            </Box>
        </Box>
    );
}

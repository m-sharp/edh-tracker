import { ReactElement } from "react";
import { Form } from "react-router-dom";
import { Box, Button, TextField } from "@mui/material";
import PublishIcon from '@mui/icons-material/Publish';

interface CreateActionProps {
    request: Request;
}

export async function createGame({request}: CreateActionProps): Promise<null> {
    const formData = await request.formData();
    const resp = await fetch(`http://localhost:8080/api/game`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "description": formData.get("description"),
            "results": [
                {
                    "game_id": -1,
                    "deck_id": parseInt(formData.get("deckOneId") as string, 10),
                    "place": parseInt(formData.get("deckOnePlace") as string, 10),
                    "kill_count": parseInt(formData.get("deckOneKills") as string, 10)
                },
                {
                    "game_id": -1,
                    "deck_id": parseInt(formData.get("deckTwoId") as string, 10),
                    "place": parseInt(formData.get("deckTwoPlace") as string, 10),
                    "kill_count": parseInt(formData.get("deckTwoKills") as string, 10)
                },
                {
                    "game_id": -1,
                    "deck_id": parseInt(formData.get("deckThreeId") as string, 10),
                    "place": parseInt(formData.get("deckThreePlace") as string, 10),
                    "kill_count": parseInt(formData.get("deckThreeKills") as string, 10)
                },
                {
                    "game_id": -1,
                    "deck_id": parseInt(formData.get("deckFourId") as string, 10),
                    "place": parseInt(formData.get("deckFourPlace") as string, 10),
                    "kill_count": parseInt(formData.get("deckFourKills") as string, 10)
                }
            ]
        }),
    });

    if ( !resp.ok ) {
        throw new Error("Failed to create new game record: received " + resp.status + " " + resp.statusText);
    }

    // ToDo: Doesn't trigger any reload, probably need to return an object back?
    return null;
}

export default function View(): ReactElement {
    // TODO: Required fields & validation?
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
                <DeckInputs keyName={"One"} />
                <DeckInputs keyName={"Two"} />
                <DeckInputs keyName={"Three"} />
                <DeckInputs keyName={"Four"} />
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
    keyName: string;
}

function DeckInputs({keyName}: DeckInputProps): ReactElement {
    return (
        <Box sx={{display: "flex", flexDirection: "column", alignItems: "center", width: "100%"}}>
            <h2>Deck {keyName}</h2>
            <Box sx={{display: "flex", flexDirection: "row", justifyContent: "space-evenly", width: "100%"}}>
                <TextField
                    required
                    name={`deck${keyName}Id`}
                    type={"number"}
                    helperText={"Deck ID"}
                />
                <TextField
                    required
                    name={`deck${keyName}Place`}
                    type={"number"}
                    helperText={"Finishing Place"}
                />
                <TextField
                    required
                    name={`deck${keyName}Kills`}
                    type={"number"}
                    helperText={"Kill Points"}
                />
            </Box>
        </Box>
    );
}

import { ReactElement } from "react";
import { Form } from "react-router-dom";
import { Box } from "@mui/material";

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
    // TODO: Should have a repeater for each deck instance
    // TODO: Subcomponent for deck input?
    // TODO: Styling
    // TODO: Textarea for description
    return (
        <Box id="newGameForm" sx={{display: "flex"}}>
            <h1>Add New Game Record</h1>
            <Form method="post">
                <input
                    type="text"
                    name="description"
                    placeholder="Add a game description!"
                />
                <Box>
                    <h2>Deck One</h2>
                    <span>ID: <input type="number" name="deckOneId" /></span>
                    <span>Place: <input type="number" name="deckOnePlace" min="1" max="4" /></span>
                    <span>Kills: <input type="number" name="deckOneKills" max="4" /></span>
                </Box>
                <Box>
                    <h2>Deck Two</h2>
                    <span>ID: <input type="number" name="deckTwoId" /></span>
                    <span>Place: <input type="number" name="deckTwoPlace" min="1" max="4" /></span>
                    <span>Kills: <input type="number" name="deckTwoKills" max="4" /></span>
                </Box>
                <Box>
                    <h2>Deck Three</h2>
                    <span>ID: <input type="number" name="deckThreeId" /></span>
                    <span>Place: <input type="number" name="deckThreePlace" min="1" max="4" /></span>
                    <span>Kills: <input type="number" name="deckThreeKills" max="4" /></span>
                </Box>
                <Box>
                    <h2>Deck Four</h2>
                    <span>ID: <input type="number" name="deckFourId" /></span>
                    <span>Place: <input type="number" name="deckFourPlace" min="1" max="4" /></span>
                    <span>Kills: <input type="number" name="deckFourKills" max="4" /></span>
                </Box>
                <button type="submit">New</button>
            </Form>
        </Box>
    );
}
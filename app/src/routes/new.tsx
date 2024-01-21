import { ReactElement } from "react";
import { Form } from "react-router-dom";

export async function createGame(): Promise<null> {
    const resp = await fetch(`http://localhost:8080/api/game`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "description": "This was an even harder fought game",
            "results": [
                {
                    "game_id": -1,
                    "deck_id": 12,
                    "place": 1,
                    "kill_count": 3
                },
                {
                    "game_id": -1,
                    "deck_id": 15,
                    "place": 2,
                    "kill_count": 0
                },
                {
                    "game_id": -1,
                    "deck_id": 29,
                    "place": 3,
                    "kill_count": 0
                },
                {
                    "game_id": -1,
                    "deck_id": 49,
                    "place": 4,
                    "kill_count": 0
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
    return (
        <Form method="post">
            <button type="submit">New</button>
        </Form>
    );
}
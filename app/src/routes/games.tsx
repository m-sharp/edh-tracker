import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";

import { MatchesDisplay, Game } from "../matches";

export async function getGames(): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games`);
    return await res.json();
}

export default function View(): ReactElement {
    const games = useLoaderData() as Array<Game>;

    return (
        <div id="games">
            <MatchesDisplay games={games} />
        </div>
    );
}

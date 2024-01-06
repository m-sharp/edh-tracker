import { useLoaderData } from "react-router-dom";

import { MatchesDisplay } from "../matches";

export async function getGames() {
    const res = await fetch(`http://localhost:8080/api/games`);
    return await res.json();
}

export default function Games() {
    const games = useLoaderData();

    return (
        <div id="games">
            <MatchesDisplay games={games} />
        </div>
    );
}

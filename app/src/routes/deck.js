import { Link, useLoaderData } from "react-router-dom";
import Skeleton from "@mui/material/Skeleton";

import { AsyncComponentHelper } from "../common";
import { MatchesDisplay } from "../matches";
import { Record } from "../stats";

export async function getDeck({ params }) {
    const res = await fetch(`http://localhost:8080/api/deck?deck_id=${params.deckId}`);
    return res.json();
}

export default function Deck() {
    const deck = useLoaderData();

    return (
        <div id="deck">
            <h1><Link to={`/player/${deck.player_id}`}>{deck.player_name}&apos;s</Link> {deck.commander}</h1>
            <p>Retired: {deck.retired.toString()}</p>
            <p>Created At: {new Date(deck.ctime).toLocaleString()}</p>
            <p>Games Played: {deck.games}</p>
            <p>Record: <Record record={deck.record}/></p>
            <p>Total Kills: {deck.kills}</p>
            <p>Total Points: {deck.points}</p>
            <MatchUpsForDeck deck={deck} />
        </div>
    );
}

function MatchUpsForDeck({ deck }) {
    const {data, loading, error} = AsyncComponentHelper(getGamesForDeck(deck.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} width={"75%"} />;
    }
    if (error) {
        return <span>Error Loading Deck's Games: {error.message}</span>;
    }

    return (
        <MatchesDisplay games={data} targetCommander={deck.commander} />
    );
}

async function getGamesForDeck(id) {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`);
    return await res.json();
}

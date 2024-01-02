import {useEffect, useState} from "react";
import {Link, useLoaderData} from "react-router-dom";

import {Record} from "../common";
import {MatchUpDisplay} from "./games";

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
            <p>Created At: {deck.ctime}</p>
            <p>Games Played: {deck.games}</p>
            <p>Record: <Record record={deck.record}/></p>
            <p>Total Kills: {deck.kills}</p>
            <p>Total Points: {deck.points}</p>
            <MatchUpsForDeck deck={deck} />
        </div>
    );
}

function MatchUpsForDeck({ deck }) {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const games = await getGamesForDeck(deck.id);
                setData(games);
                setLoading(false);
            } catch (error) {
                setError(error);
                setLoading(false);
            }
        }

        fetchData()
    }, []);

    if (loading) {
        // ToDo: Get a spinner
        return <span>Loading...</span>
    }
    if (error) {
        return <span>Error: {error.message}</span>
    }

    // ToDo: Repeat of <Games />?
    return (
        <div id="games">
            <ul>
                {data.map(game => (
                    <li key={game.id}>
                        <MatchUpDisplay game={game} />
                    </li>
                ))}
            </ul>
        </div>
    );
}

async function getGamesForDeck(id) {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`);
    const games = await res.json();

    return games.map((game) => ({
        id: game.id.toString(),
        description: game.description,
        ctime: game.ctime,
        results: game.results,
    }));
}

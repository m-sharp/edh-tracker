import {Link, useLoaderData} from "react-router-dom";

export async function getDecks() {
    const res = await fetch(`http://localhost:8080/api/decks`);
    const decks = await res.json();

    return decks.map((deck) => ({
        id: deck.id.toString(),
        player_id: deck.player_id,
        commander: deck.commander,
        retired: deck.retired,
        ctime: deck.ctime,
    }));
}

export default function Decks() {
    const decks = useLoaderData();

    return (
        <div id="decks">
            <ul>
                {decks.map(deck => (
                    <li key={deck.id}>
                        <Link to={`/deck/${deck.id}`}>{deck.commander}</Link>
                    </li>
                ))}
            </ul>
        </div>
    );
}

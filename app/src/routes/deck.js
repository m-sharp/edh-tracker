import { useLoaderData } from "react-router-dom";

export async function getDeck({ params }) {
    const res = await fetch(`http://localhost:8080/api/deck?deck_id=${params.deckId}`);
    return res.json();
}

export default function Deck() {
    const deck = useLoaderData();

    return (
        <div id="deck">
            {/*ToDo: Make this "Mike's Giada...*/}
            <h1>Deck #{deck.id} - {deck.commander}</h1>
            <p>PlayerID: {deck.player_id}</p>
            <p>Retired: {deck.retired.toString()}</p>
            <p>Created At: {deck.ctime}</p>
        </div>
    )
}

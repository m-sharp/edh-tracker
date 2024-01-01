import {Link, useLoaderData} from "react-router-dom";

export async function getPlayer({ params }) {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${params.playerId}`);
    return res.json();
}

async function getDecksForPlayer(id) {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    const decks = await res.json();

    return decks.map((deck) => ({
        id: deck.id.toString(),
        player_id: deck.player_id,
        commander: deck.commander,
        retired: deck.retired,
        ctime: deck.ctime,
    }));
}

export default function Player() {
    const player = useLoaderData();

    // ToDo: Need to figure out how to get at this - can't do an await here as the route element can't be a promise
    // const playerDecks = await getDecksForPlayer(player.id);
    // const decks = playerDecks.map(deck =>
    //     <li key={deck.id}>
    //         <Link to={`/deck/${deck.id}`}>{deck.commander}</Link>
    //     </li>
    // );

    return (
        <div id="player">
            <h1>{player.name}&apos;s Page!</h1>
            <p>Player created time: {player.ctime}</p>
            <p>Player Total Kills: {player.kills}</p>
            <p>Player Record: {player.record[1]} / {player.record[2]} / {player.record[3]} / {player.record[4]}</p>
            {/*<p>Player's Decks:</p>*/}
            {/*<ol>{decks}</ol>*/}
        </div>
    );
}
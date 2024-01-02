import {Link, useLoaderData} from "react-router-dom";

export async function getGame({ params }) {
    const res = await fetch(`http://localhost:8080/api/game?game_id=${params.gameId}`);
    return res.json();
}

export default function Game() {
    const game = useLoaderData();

    return (
        <div id="game">
            <h1>Game #{game.id} Results</h1>
            <p>Description: {game.description}</p>
            {game.results.map(result => (
                <ResultDisplay result={result} />
            ))}
        </div>
    );
}

function ResultDisplay({ result }) {
    return (
        <div>
            <h2><Link to={`/deck/${result.deck_id}`}>{result.commander}</Link></h2>
            <ul>
                <li>Place: {result.place}</li>
                <li>Kills: {result.kill_count}</li>
                <li>Points: {result.points}</li>
            </ul>
        </div>
    );
}

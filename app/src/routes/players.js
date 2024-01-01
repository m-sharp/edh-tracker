import { Link, useLoaderData } from "react-router-dom";

export async function getPlayers() {
    const res = await fetch(`http://localhost:8080/api/players`);
    const players = await res.json();

    return players.map((player) => ({
        id: player.id.toString(),
        name: player.name,
        ctime: player.ctime,
    }));
}

export default function Players() {
    const players = useLoaderData();

    const playerItems = players.map(player =>
        <li key={player.id}>
            <Link to={`/player/${player.id}`}>{player.name}</Link>
        </li>
    );

    return (
        <div id="players">
            <ul>{playerItems}</ul>
        </div>
    );
}

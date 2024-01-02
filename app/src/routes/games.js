import {Link, useLoaderData} from "react-router-dom";

export async function getGames() {
    const res = await fetch(`http://localhost:8080/api/games`);
    const games = await res.json();

    return games.map((game) => ({
        id: game.id.toString(),
        description: game.description,
        ctime: game.ctime,
        results: game.results,
    }));
}

export default function Games() {
    const games = useLoaderData();

    return (
        <div id="games">
            <ul>
                {games.map(game => (
                    <li key={game.id}>
                        <MatchUpDisplay game={game} />
                    </li>
                ))}
            </ul>
        </div>
    );
}

export function MatchUpDisplay({ game }) {
    return (
      <span>
          <span><Link to={`/game/${game.id}`}>{game.ctime}</Link> - </span>
          {game.results.map(result => (
              <span>
                <span color={ result.place === 1 ? "green" : "black"}>{result.commander}</span>
                {/*TODO: Hide with css on last child*/}
                <span> VS </span>
              </span>
          ))}
      </span>
    );
}

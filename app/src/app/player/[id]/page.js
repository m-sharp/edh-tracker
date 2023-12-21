import { usePathname } from 'next/navigation'

export async function generateStaticParams() {
    const res = await fetch("http://localhost:8080/api/players");
    const players = await res.json();

    return players.map((player) => ({
        id: player.id.toString(),
        name: player.name,
        ctime: player.ctime,
    }));
}

async function getPlayer(id) {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${id}`);
    return res.json();
}

export default async function Page({ params }) {
    const player = await getPlayer(params.id);

    return (
        <main>
            <h1>Hello, Player {player.name}&apos;s Page!</h1>
            <p>Player created time: {player.ctime}</p>
        </main>
    );
}

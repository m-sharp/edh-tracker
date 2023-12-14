import Link from 'next/link'

export default function Page() {
    return (
        <main>
            <h1>Hello, Players Page!</h1>
            <Link href="/player/1">Player 1</Link>
        </main>
    )
}

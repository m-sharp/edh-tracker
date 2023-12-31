import { Outlet, Link } from "react-router-dom";

export default function Root() {
    return (
        <>
            <nav>
                <ul>
                    <li><Link to={`/players`}>Players</Link></li>
                    <li><Link to={`/decks`}>Decks</Link></li>
                </ul>
            </nav>
            <hr />
            <div id="detail">
                <Outlet />
            </div>
        </>
    )
}

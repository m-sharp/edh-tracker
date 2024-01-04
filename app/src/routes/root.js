import { Outlet, Link } from "react-router-dom";

export default function Root() {
    // ToDo: https://mui.com/material-ui/react-app-bar/
    // ToDo: https://mui.com/material-ui/react-drawer/
    return (
        <>
            <nav>
                <ul>
                    <li><Link to={`/players`}>Players</Link></li>
                    <li><Link to={`/decks`}>Decks</Link></li>
                    <li><Link to={`/games`}>Games</Link></li>
                </ul>
            </nav>
            <hr />
            <div id="detail">
                <Outlet />
            </div>
        </>
    )
}

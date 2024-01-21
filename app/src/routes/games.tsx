import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";

import { MatchesDisplay, Game } from "../matches";

export async function getGames(): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games`);
    return await res.json();
}

export default function View(): ReactElement {
    const games = useLoaderData() as Array<Game>;

    return (
        <Box id="games" style={{height: 500}}>
            <MatchesDisplay games={games} />
        </Box>
    );
}

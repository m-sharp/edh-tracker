import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";

import { MatchesDisplay } from "../matches";
import { Game } from "../types";

export default function View(): ReactElement {
    const games = useLoaderData() as Array<Game>;

    return (
        <Box id="games" style={{height: 500}}>
            <MatchesDisplay games={games} />
        </Box>
    );
}

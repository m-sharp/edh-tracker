import { ReactElement } from "react";
import { useLoaderData } from "react-router-dom";
import { Box } from "@mui/material";

import { MatchesDisplay } from "../matches";
import { Game } from "../types";

// TODO: Big clean up after all pages revamped?
export default function View(): ReactElement {
    // TODO: Should be getting games for a given pod. Will be supplanted by the new pods view described in TODOs in @app/src/routes/decks.tsx
    // TODO: Currently errors as no query string is provided

    const games = useLoaderData() as Array<Game>;

    return (
        <Box id="games" style={{height: 500}}>
            <MatchesDisplay games={games} />
        </Box>
    );
}

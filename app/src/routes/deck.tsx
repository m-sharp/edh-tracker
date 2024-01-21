import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import DeleteIcon from '@mui/icons-material/Delete';
import { Box } from "@mui/material";
import Skeleton from "@mui/material/Skeleton";

import { AsyncComponentHelper } from "../common";
import { Game, MatchesDisplay } from "../matches";
import { Record, RecordDict } from "../stats";

import "../styles.css";

export interface Deck {
    id: number;
    player_id: number;
    player_name: string;
    commander: string;
    retired: boolean;
    ctime: string;
    record: RecordDict;
    games: number;
    kills: number;
    points: number;
}

export async function getDeck({ params }: LoaderFunctionArgs): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/deck?deck_id=${params.deckId}`);
    return res.json();
}

export default function View(): ReactElement {
    const deck = useLoaderData() as Deck;

    return (
        <Box id="deck" sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
            <Box sx={{display: "flex", flexDirection: "column", alignItems: "center"}}>
                <h1>{deck.commander}</h1>
                <Record record={deck.record} />
                <span className={"info-column-item"}>Owner - <Link to={`/player/${deck.player_id}`}>{deck.player_name}</Link></span>
                {!deck.retired &&
                    <Box sx={{display: "flex"}}><DeleteIcon /><span> Retired</span></Box>
                }
            </Box>
            <Box sx={{width: "100%", display: "flex", flexDirection: "row", justifyContent: "space-evenly", py: 3}}>
                <span><strong>Games Played:</strong> {deck.games}</span>
                <span><strong>Total Kills:</strong> {deck.kills}</span>
                <span><strong>Total Points:</strong> {deck.points}</span>
            </Box>
            <MatchUpsForDeck deck={deck} />
            <Box sx={{width: "100%", display: "flex", justifyContent: "flex-end", pt: 1}}>
                <em>Deck created at: {new Date(deck.ctime).toLocaleString()}</em>
            </Box>
        </Box>
    );
}

function MatchUpsForDeck({ deck }: MatchUpsForDeckProps): ReactElement {
    const {data, loading, error} = AsyncComponentHelper(getGamesForDeck(deck.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} />;
    }
    if (error) {
        return <span>Error Loading Deck's Games: {error.message}</span>;
    }

    return (
        <Box style={{height: 500, width: "100%"}}>
            <MatchesDisplay games={data} targetCommander={deck.commander} />
        </Box>
    );
}

interface MatchUpsForDeckProps {
    deck: Deck;
}

async function getGamesForDeck(id: number): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`);
    return await res.json();
}

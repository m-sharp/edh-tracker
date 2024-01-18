import { ReactElement } from "react";
import { Link, useLoaderData } from "react-router-dom";
import { LoaderFunctionArgs } from "@remix-run/router/utils";
import Skeleton from "@mui/material/Skeleton";

import { AsyncComponentHelper } from "../common";
import { Game, MatchesDisplay } from "../matches";
import { Record, RecordDict } from "../stats";

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
        <div id="deck">
            <h1><Link to={`/player/${deck.player_id}`}>{deck.player_name}&apos;s</Link> {deck.commander}</h1>
            <p>Retired: {deck.retired.toString()}</p>
            <p>Created At: {new Date(deck.ctime).toLocaleString()}</p>
            <p>Games Played: {deck.games}</p>
            <p>Record: <Record record={deck.record}/></p>
            <p>Total Kills: {deck.kills}</p>
            <p>Total Points: {deck.points}</p>
            <MatchUpsForDeck deck={deck} />
        </div>
    );
}

interface MatchUpsForDeckProps {
    deck: Deck;
}

function MatchUpsForDeck({ deck }: MatchUpsForDeckProps): ReactElement {
    const {data, loading, error} = AsyncComponentHelper(getGamesForDeck(deck.id));

    if (loading) {
        return <Skeleton variant="rounded" animation="wave" height={500} width={"75%"} />;
    }
    if (error) {
        return <span>Error Loading Deck's Games: {error.message}</span>;
    }

    return (
        <MatchesDisplay games={data} targetCommander={deck.commander} />
    );
}

async function getGamesForDeck(id: number): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`);
    return await res.json();
}

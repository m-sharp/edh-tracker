import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { Commander, Deck, Format, Game, NewGame, NewGameData, Player } from "./types";

// ToDo: These endpoints either need to be relative or configurable somehow
// Player Methods
export async function GetPlayer({ params }: LoaderFunctionArgs): Promise<Player> {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${params.playerId}`);
    return res.json();
}

export async function GetPlayers(): Promise<Array<Player>> {
    const res = await fetch(`http://localhost:8080/api/players`);
    return await res.json();
}

export async function GetDecksForPlayer(id: number): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`);
    return await res.json();
}

// Deck Methods
export async function GetDeck({ params }: LoaderFunctionArgs): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/deck?deck_id=${params.deckId}`);
    return res.json();
}

export async function GetDecks(): Promise<Array<Deck>> {
    const res = await fetch(`http://localhost:8080/api/decks`);
    return await res.json();
}

export async function GetGamesForDeck(id: number): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`);
    return await res.json();
}

// Game Methods
export async function GetGame({ params }: LoaderFunctionArgs): Promise<Game> {
    const res = await fetch(`http://localhost:8080/api/game?game_id=${params.gameId}`);
    return res.json();
}

export async function GetGames(): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games`);
    return await res.json();
}

export async function PostGame(newGame: NewGame): Promise<Response> {
    return await fetch(`http://localhost:8080/api/game`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(newGame),
    });
}

// Format Methods
export async function GetFormats(): Promise<Array<Format>> {
    const res = await fetch(`http://localhost:8080/api/formats`);
    return await res.json();
}

// Commander Methods
export async function GetCommander(id: number): Promise<Commander> {
    const res = await fetch(`http://localhost:8080/api/commander?commander_id=${id}`);
    return await res.json();
}

export async function PostCommander(name: string): Promise<Response> {
    return await fetch(`http://localhost:8080/api/commander`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name }),
    });
}

// Loaders
export async function GetNewDeckInfo(): Promise<NewGameData> {
    const [decks, players, formats] = await Promise.all([
        GetDecks(),
        GetPlayers(),
        GetFormats(),
    ]);

    return {
        decks,
        players,
        formats,
    }
}

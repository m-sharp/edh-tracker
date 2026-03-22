import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { AuthUser } from "./auth";
import {
    Commander,
    Deck,
    DeckUpdateFields,
    Format,
    Game,
    GameResultUpdateFields,
    NewGame,
    NewGameData,
    NewGameResultWithGame,
    PaginatedResponse,
    Player,
    PlayerWithRole,
    Pod,
} from "./types";

// ToDo: These endpoints either need to be relative or configurable somehow

// Auth Methods
export async function GetMe(): Promise<AuthUser> {
    const res = await fetch(`http://localhost:8080/api/auth/me`, { credentials: "include" });
    if (!res.ok) throw new Error("Unauthenticated");
    return res.json();
}

export async function Logout(): Promise<void> {
    await fetch(`http://localhost:8080/api/auth/logout`, {
        method: "POST",
        credentials: "include",
    });
}

// Player Methods
export async function GetPlayer({ params }: LoaderFunctionArgs): Promise<Player> {
    const res = await fetch(`http://localhost:8080/api/player?player_id=${params.playerId}`, { credentials: "include" });
    return res.json();
}

export async function GetPlayers(): Promise<Array<Player>> {
    const res = await fetch(`http://localhost:8080/api/players`, { credentials: "include" });
    return await res.json();
}

export async function GetPlayersForPod(podId: number): Promise<Array<PlayerWithRole>> {
    const res = await fetch(`http://localhost:8080/api/players?pod_id=${podId}`, { credentials: "include" });
    return await res.json();
}

export async function GetDecksForPlayer(id: number): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/decks?player_id=${id}`, { credentials: "include" });
    return await res.json();
}

export async function GetGamesForPlayer(playerId: number): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?player_id=${playerId}`, { credentials: "include" });
    return await res.json();
}

export async function PatchPlayer(playerId: number, name: string): Promise<void> {
    await fetch(`http://localhost:8080/api/player?player_id=${playerId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
}

// Deck Methods
export async function GetDeck({ params }: LoaderFunctionArgs): Promise<Deck> {
    const res = await fetch(`http://localhost:8080/api/deck?deck_id=${params.deckId}`, { credentials: "include" });
    return res.json();
}

export async function GetDecks(): Promise<Array<Deck>> {
    const res = await fetch(`http://localhost:8080/api/decks`, { credentials: "include" });
    return await res.json();
}

export async function GetDecksForPod(podId: number, limit: number, offset: number): Promise<PaginatedResponse<Deck>> {
    const res = await fetch(`http://localhost:8080/api/decks?pod_id=${podId}&limit=${limit}&offset=${offset}`, { credentials: "include" });
    return await res.json();
}

export async function GetGamesForDeck(id: number): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?deck_id=${id}`, { credentials: "include" });
    return await res.json();
}

export async function PatchDeck(deckId: number, fields: DeckUpdateFields): Promise<void> {
    await fetch(`http://localhost:8080/api/deck?deck_id=${deckId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(fields),
    });
}

export async function DeleteDeck(deckId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/deck?deck_id=${deckId}`, {
        method: "DELETE",
        credentials: "include",
    });
}

// Game Methods
export async function GetGame({ params }: LoaderFunctionArgs): Promise<Game> {
    const res = await fetch(`http://localhost:8080/api/game?game_id=${params.gameId}`, { credentials: "include" });
    return res.json();
}

export async function GetGames(): Promise<Array<Game>> {
    const res = await fetch(`http://localhost:8080/api/games`, { credentials: "include" });
    return await res.json();
}

export async function GetGamesForPod(podId: number, limit: number, offset: number): Promise<PaginatedResponse<Game>> {
    const res = await fetch(`http://localhost:8080/api/games?pod_id=${podId}&limit=${limit}&offset=${offset}`, { credentials: "include" });
    return await res.json();
}

export async function PostGame(newGame: NewGame): Promise<Response> {
    return await fetch(`http://localhost:8080/api/game`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(newGame),
    });
}

export async function PatchGame(gameId: number, description: string): Promise<void> {
    await fetch(`http://localhost:8080/api/game?game_id=${gameId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ description }),
    });
}

export async function DeleteGame(gameId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/game?game_id=${gameId}`, {
        method: "DELETE",
        credentials: "include",
    });
}

export async function PostGameResult(result: NewGameResultWithGame): Promise<void> {
    await fetch(`http://localhost:8080/api/game/result`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(result),
    });
}

export async function PatchGameResult(resultId: number, fields: GameResultUpdateFields): Promise<void> {
    await fetch(`http://localhost:8080/api/game/result?result_id=${resultId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(fields),
    });
}

export async function DeleteGameResult(resultId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/game/result?result_id=${resultId}`, {
        method: "DELETE",
        credentials: "include",
    });
}

// Format Methods
export async function GetFormats(): Promise<Array<Format>> {
    const res = await fetch(`http://localhost:8080/api/formats`, { credentials: "include" });
    return await res.json();
}

// Commander Methods
export async function GetCommander(id: number): Promise<Commander> {
    const res = await fetch(`http://localhost:8080/api/commander?commander_id=${id}`, { credentials: "include" });
    return await res.json();
}

export async function PostCommander(name: string): Promise<Response> {
    return await fetch(`http://localhost:8080/api/commander`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
}

// Pod Methods
export async function GetPod(podId: number): Promise<Pod> {
    const res = await fetch(`http://localhost:8080/api/pod?pod_id=${podId}`, { credentials: "include" });
    return await res.json();
}

export async function GetPodsForPlayer(playerId: number): Promise<Array<Pod>> {
    const res = await fetch(`http://localhost:8080/api/pod?player_id=${playerId}`, { credentials: "include" });
    return await res.json();
}

export async function PostPod(name: string): Promise<Pod> {
    const res = await fetch(`http://localhost:8080/api/pod`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
    return await res.json();
}

export async function PatchPod(podId: number, name: string): Promise<void> {
    await fetch(`http://localhost:8080/api/pod?pod_id=${podId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
}

export async function DeletePod(podId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/pod?pod_id=${podId}`, {
        method: "DELETE",
        credentials: "include",
    });
}

export async function PostPodInvite(podId: number): Promise<{ invite_code: string }> {
    const res = await fetch(`http://localhost:8080/api/pod/invite`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId }),
    });
    return await res.json();
}

export async function PostPodJoin(inviteCode: string): Promise<Pod> {
    const res = await fetch(`http://localhost:8080/api/pod/join`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ invite_code: inviteCode }),
    });
    return await res.json();
}

export async function PostPodLeave(podId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/pod/leave`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId }),
    });
}

export async function PatchPodPlayerRole(podId: number, playerId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/pod/player?pod_id=${podId}&player_id=${playerId}`, {
        method: "PATCH",
        credentials: "include",
    });
}

export async function DeletePodPlayer(podId: number, playerId: number): Promise<void> {
    await fetch(`http://localhost:8080/api/pod/player?pod_id=${podId}&player_id=${playerId}`, {
        method: "DELETE",
        credentials: "include",
    });
}

// Loaders
// TODO: Need to revist this, probably changes a lot
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

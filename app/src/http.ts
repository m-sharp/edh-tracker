import { LoaderFunctionArgs } from "@remix-run/router/utils";

import { AuthUser } from "./auth";
import {
    Commander,
    Deck,
    DeckUpdateFields,
    Format,
    Game,
    GameResultUpdateFields,
    NewDeckRequest,
    NewGame,
    NewGameData,
    NewGameResultWithGame,
    PaginatedResponse,
    Player,
    PlayerWithRole,
    Pod,
} from "./types";

// TODO: How does this look in production?
export const API_BASE_URL = "http://localhost:8080";

// Auth Methods
export async function GetMe(): Promise<AuthUser> {
    const res = await fetch(`${API_BASE_URL}/api/auth/me`, { credentials: "include" });
    if (!res.ok) throw new Error("Unauthenticated");
    return res.json();
}

export async function Logout(): Promise<void> {
    await fetch(`${API_BASE_URL}/api/auth/logout`, {
        method: "POST",
        credentials: "include",
    });
}

// Player Methods
export async function GetPlayer({ params }: LoaderFunctionArgs): Promise<Player> {
    const res = await fetch(`${API_BASE_URL}/api/player?player_id=${params.playerId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return res.json();
}

export async function GetPlayers(): Promise<Array<Player>> {
    const res = await fetch(`${API_BASE_URL}/api/players`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetPlayersForPod(podId: number): Promise<Array<PlayerWithRole>> {
    const res = await fetch(`${API_BASE_URL}/api/players?pod_id=${podId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetDecksForPlayer(id: number): Promise<Deck> {
    const res = await fetch(`${API_BASE_URL}/api/decks?player_id=${id}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetGamesForPlayer(playerId: number): Promise<Array<Game>> {
    const res = await fetch(`${API_BASE_URL}/api/games?player_id=${playerId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PatchPlayer(playerId: number, name: string): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/player?player_id=${playerId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

// Deck Methods
export async function GetDeck({ params }: LoaderFunctionArgs): Promise<Deck> {
    const res = await fetch(`${API_BASE_URL}/api/deck?deck_id=${params.deckId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return res.json();
}

export async function GetDecks(): Promise<Array<Deck>> {
    const res = await fetch(`${API_BASE_URL}/api/decks`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetDecksForPod(podId: number, limit: number, offset: number): Promise<PaginatedResponse<Deck>> {
    const res = await fetch(`${API_BASE_URL}/api/decks?pod_id=${podId}&limit=${limit}&offset=${offset}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetGamesForDeck(id: number): Promise<Array<Game>> {
    const res = await fetch(`${API_BASE_URL}/api/games?deck_id=${id}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PatchDeck(deckId: number, fields: DeckUpdateFields): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/deck?deck_id=${deckId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(fields),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function DeleteDeck(deckId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/deck?deck_id=${deckId}`, {
        method: "DELETE",
        credentials: "include",
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

// Game Methods
export async function GetGame({ params }: LoaderFunctionArgs): Promise<Game> {
    const res = await fetch(`${API_BASE_URL}/api/game?game_id=${params.gameId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return res.json();
}

export async function GetGames(): Promise<Array<Game>> {
    const res = await fetch(`${API_BASE_URL}/api/games`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetGamesForPod(podId: number, limit: number, offset: number): Promise<PaginatedResponse<Game>> {
    const res = await fetch(`${API_BASE_URL}/api/games?pod_id=${podId}&limit=${limit}&offset=${offset}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PostGame(newGame: NewGame): Promise<Response> {
    return await fetch(`${API_BASE_URL}/api/game`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(newGame),
    });
}

export async function PatchGame(gameId: number, description: string): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/game?game_id=${gameId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ description }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function DeleteGame(gameId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/game?game_id=${gameId}`, {
        method: "DELETE",
        credentials: "include",
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function PostGameResult(result: NewGameResultWithGame): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/game/result`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(result),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function PatchGameResult(resultId: number, fields: GameResultUpdateFields): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/game/result?result_id=${resultId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(fields),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function DeleteGameResult(resultId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/game/result?result_id=${resultId}`, {
        method: "DELETE",
        credentials: "include",
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

// Format Methods
export async function GetFormats(): Promise<Array<Format>> {
    const res = await fetch(`${API_BASE_URL}/api/formats`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

// Commander Methods
export async function GetCommanders(): Promise<Array<Commander>> {
    const res = await fetch(`${API_BASE_URL}/api/commanders`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetCommander(id: number): Promise<Commander> {
    const res = await fetch(`${API_BASE_URL}/api/commander?commander_id=${id}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PostCommander(name: string): Promise<Response> {
    return await fetch(`${API_BASE_URL}/api/commander`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
}

export async function PostDeck(body: NewDeckRequest): Promise<{ id: number }> {
    const res = await fetch(`${API_BASE_URL}/api/deck`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(body),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

// Pod Methods
export async function GetPod(podId: number): Promise<Pod> {
    const res = await fetch(`${API_BASE_URL}/api/pod?pod_id=${podId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function GetPodsForPlayer(playerId: number): Promise<Array<Pod>> {
    const res = await fetch(`${API_BASE_URL}/api/pod?player_id=${playerId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PostPod(name: string): Promise<Pod> {
    const res = await fetch(`${API_BASE_URL}/api/pod`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PatchPod(podId: number, name: string): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/pod?pod_id=${podId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ name }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function DeletePod(podId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/pod?pod_id=${podId}`, {
        method: "DELETE",
        credentials: "include",
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function PostPodInvite(podId: number): Promise<{ invite_code: string }> {
    const res = await fetch(`${API_BASE_URL}/api/pod/invite`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PostPodJoin(inviteCode: string): Promise<Pod> {
    const res = await fetch(`${API_BASE_URL}/api/pod/join`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ invite_code: inviteCode }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

export async function PostPodLeave(podId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/pod/leave`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId }),
    });
    if (!res.ok) {
        const err: any = new Error(`Failed to leave pod: ${res.status}`);
        err.status = res.status;
        throw err;
    }
}

export async function PatchPodPlayerRole(podId: number, playerId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/pod/player`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId, player_id: playerId }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function DeletePodPlayer(podId: number, playerId: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/api/pod/player`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ pod_id: podId, player_id: playerId }),
    });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
}

export async function GetAllDecksForPod(podId: number): Promise<Array<Deck>> {
    const res = await fetch(`${API_BASE_URL}/api/decks?pod_id=${podId}`, { credentials: "include" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return await res.json();
}

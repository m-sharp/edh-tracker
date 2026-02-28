export interface RecordDict {
    [key: number]: number;
}

export interface Format {
    id: number;
    name: string;
}

export interface Commander {
    id: number;
    name: string;
}

export interface DeckCommanderInfo {
    commander_id: number;
    commander_name: string;
    partner_commander_id?: number;
    partner_commander_name?: string;
}

export interface Deck {
    id: number;
    player_id: number;
    player_name: string;
    name: string;
    format_id: number;
    format_name: string;
    commanders?: DeckCommanderInfo;
    retired: boolean;
    created_at: string;
    updated_at: string;
    record: RecordDict;
    games: number;
    kills: number;
    points: number;
}

export interface Player {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
    record: RecordDict;
    games: number;
    kills: number;
    points: number;
}

export interface Game {
    id: number;
    description: string;
    format_id: number;
    created_at: string;
    updated_at: string;
    results: Array<GameResult>;
}

export interface GameResult {
    id: number;
    game_id: number;
    deck_id: number;
    deck_name: string;
    commander_name?: string;
    partner_commander_name?: string;
    place: number;
    kill_count: number;
    points: number;
    created_at: string;
    updated_at: string;
}

export interface NewGameData {
    players: Array<Player>;
    decks: Array<Deck>;
    formats: Array<Format>;
}

export interface NewGame {
    description: string;
    format_id: number;
    pod_id: number;
    results: Array<NewGameResult>;
}

export interface NewGameResult {
    deck_id: number;
    player_id: number;
    place: number;
    kill_count: number;
}

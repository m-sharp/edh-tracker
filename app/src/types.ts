export interface RecordDict {
    [key: number]: number;
}

export interface Deck {
    id: number;
    player_id: number;
    player_name: string;
    commander: string;
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
    created_at: string;
    updated_at: string;
    results: Array<GameResult>;
}

export interface GameResult {
    id: number;
    game_id: number;
    deck_id: number;
    commander: string;
    place: number;
    kill_count: number;
    points: number;
    created_at: string;
    updated_at: string;
}

export interface NewGameData {
    players: Array<Player>;
    decks: Array<Deck>;
}

export interface NewGame {
    description: string;
    results: Array<NewGameResult>;
}

export interface NewGameResult {
    deck_id: number;
    player_id: number;
    place: number;
    kill_count: number;
}

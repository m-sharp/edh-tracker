import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { DataGrid, GridColDef, GridToolbar } from "@mui/x-data-grid";

import { CreatedAtColumn } from "./stats";

export interface Game {
    id: number;
    description: string;
    ctime: string;
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
}

interface MatchesDisplayProps {
    games: Array<Game>;
    targetCommander?: string;
}

export function MatchesDisplay({ games, targetCommander }: MatchesDisplayProps): ReactElement {
    const columns: Array<GridColDef> = [
        {
            field: "id",
            headerName: "Game #",
            renderCell: (params) => (
                <Link to={`/game/${params.row.id}`}>#{params.row.id}</Link>
            ),
            minWidth: 100,
        },
        {
            field: "results",
            headerName: "Commanders (In Place Order)",
            renderCell: (params) => (
                <MatchUpDisplay results={params.row.results} targetCommander={targetCommander || ""} />
            ),
            hideable: false,
            sortable: false,
            flex: 1,
        },
        CreatedAtColumn,
    ];

    return (
        <div id="matches" style={{height: 500, width: "75%"}}>
            <DataGrid
                rows={games}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "id", sort: "asc"}],
                    },
                }}
            />
        </div>
    );
}

interface MatchUpDisplayProps {
    results: Array<GameResult>;
    targetCommander: string;
}

export function MatchUpDisplay({ results, targetCommander }: MatchUpDisplayProps): ReactElement {
    results.sort((a, b) => {
        if (a.place < b.place) {
            return -1;
        } else if (a.place > b.place) {
            return 1;
        }

        return 0;
    });

    return (
        <span className="match">
            {results.map(result => (
                <span className="match-up" key={result.id}>
                    <span color={result.commander === targetCommander ? "blue" : "black"}>{result.commander}</span>
                    <span className="vs"> VS </span>
                </span>
            ))}
        </span>
    );
}

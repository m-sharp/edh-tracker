import { Link } from "react-router-dom";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { CreatedAtColumn } from "./stats";

export function MatchesDisplay({ games, targetCommander }) {
    const columns = [
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
                <MatchUpDisplay results={params.row.results} targetCommander={targetCommander} />
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

export function MatchUpDisplay({ results, targetCommander }) {
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

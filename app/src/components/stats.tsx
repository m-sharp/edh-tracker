import { ReactElement } from "react";
import { Link } from "react-router-dom";
import { GridColDef } from "@mui/x-data-grid";

import { RecordDict } from "../types";

interface RecordProps {
    record: RecordDict;
}

// Record takes a Record dictionary like {1: 10, 2: 12, 3: 7, 4: 5}
export function Record({ record }: RecordProps): ReactElement {
    const maxPlace = Math.max(...Object.keys(record).map(Number), 4);
    const parts = Array.from({ length: maxPlace }, (_, i) => record[i + 1] ?? 0);
    return <span className="record">{parts.join(" / ")}</span>;
}

// RecordComparator is a custom DataGrid comparator for Record. Iterates from place 1 to max place in either dict, returns the first non-zero difference.
export function RecordComparator(record1: RecordDict, record2: RecordDict): number {
    const maxPlace = Math.max(
        ...Object.keys(record1).map(Number),
        ...Object.keys(record2).map(Number),
        1
    );
    for (let place = 1; place <= maxPlace; place++) {
        const diff = (record1[place] ?? 0) - (record2[place] ?? 0);
        if (diff !== 0) return diff;
    }
    return 0;
}

// StatColumns returns a list of DataGrid column definitions for the common game stats
export const StatColumns: Array<GridColDef> = [
    {
        field: "record",
        headerName: "Record",
        valueGetter: (params) => params.row.stats?.record,
        renderCell: (params) => (
            <Record record={params.value} />
        ),
        sortComparator: RecordComparator,
        flex: 1,
        minWidth: 150,
    },
    {
        field: "kills",
        headerName: "Total Kills",
        type: "number",
        valueGetter: (params) => params.row.stats?.kills,
        minWidth: 125,
    },
    {
        field: "points",
        headerName: "Total Points",
        type: "number",
        valueGetter: (params) => params.row.stats?.points,
        minWidth: 150,
    },
    {
        field: "games",
        headerName: "Games Played",
        type: "number",
        valueGetter: (params) => params.row.stats?.games,
        minWidth: 150,
    },
];

// CommanderColumn is a DataGrid column definition for a deck name, formatted as a <Link />
export const CommanderColumn: GridColDef = {
    field: "name",
    headerName: "Deck",
    renderCell: (params) => (
        <Link to={`/player/${params.row.player_id}/deck/${params.row.id}`}>{params.row.name}</Link>
    ),
    hideable: false,
    flex: 1,
    minWidth: 200,
};

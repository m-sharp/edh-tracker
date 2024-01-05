import { Link } from "react-router-dom";

// Record takes a record dictionary like {1: 10, 2: 12, 3: 7, 4: 5}
export function Record({ record }) {
    let first = getter(record, 1);
    let second = getter(record, 2);
    let third = getter(record, 3);
    let fourth = getter(record, 4);

    return (
        <span id="record">{first} / {second} / {third} / {fourth}</span>
    )
}

// RecordComparator is a custom DataGrid comparator for a Record to enable sorting
export function RecordComparator(record1, record2) {
    const firsts = getter(record1, 1) - getter(record2, 1);
    if (firsts !== 0) {
        return firsts;
    }

    const seconds = getter(record1, 2) - getter(record2, 2);
    if (seconds !== 0) {
        return seconds;
    }

    const thirds = getter(record1, 3) - getter(record2, 3);
    if (thirds !== 0) {
        return thirds;
    }

    return getter(record1, 4) - getter(record2, 4);
}

function getter(m, target) {
    return m[target] || 0;
}

// StatColumns returns a list of DataGrid column definitions for the common game stats
export const StatColumns = [
    {
        field: "record",
        headerName: "Record",
        renderCell: (params) => (
            <Record record={params.row.record} />
        ),
        sortComparator: RecordComparator,
        minWidth: 150,
    },
    {
        field: "kills",
        headerName: "Total Kills",
        type: "number",
        minWidth: 100,
    },
    {
        field: "points",
        headerName: "Total Points",
        type: "number",
        minWidth: 100,
    },
    {
        field: "games",
        headerName: "Games Played",
        type: "number",
        minWidth: 100,
    },
];

// CommanderColumn is a DataGrid column definition for a commander name string, formatted as a <Link />
export const CommanderColumn = {
    field: "commander",
    headerName: "Commander",
    renderCell: (params) => (
        <Link to={`/deck/${params.row.id}`}>{params.row.commander}</Link>
    ),
    hideable: false,
    flex: 1,
};

// CreatedAtColumn is a DataGrid column definition for a ctime datetime.
export const CreatedAtColumn = {
    field: "ctime",
    headerName: "Created At",
    type: "dateTime",
    valueGetter: ({ value }) => value && new Date(value),
    minWidth: 225,
};

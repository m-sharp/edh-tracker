import { Link, useLoaderData } from "react-router-dom";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { CreatedAtColumn, StatColumns } from "../common";

export async function getPlayers() {
    const res = await fetch(`http://localhost:8080/api/players`);
    return await res.json();
}

export default function Players() {
    const players = useLoaderData();

    const columns = [
        {
            field: "name",
            headerName: "Player Name",
            renderCell: (params) => (
                <Link to={`/player/${params.row.id}`}>{params.row.name}</Link>
            ),
            hideable: false,
            flex: 1,
        },
        ...StatColumns,
        CreatedAtColumn,
    ];

    return (
        <div id="players" style={{height: 500, width: "75%"}}>
            <DataGrid
                rows={players}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "record", sort: "desc" }],
                    },
                }}
            />
        </div>
    );
}

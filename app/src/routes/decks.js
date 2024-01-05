import { useLoaderData } from "react-router-dom";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import { CommanderColumn, CreatedAtColumn, StatColumns } from "../common";

export async function getDecks() {
    const res = await fetch(`http://localhost:8080/api/decks`);
    return await res.json();
}

export default function Decks() {
    const decks = useLoaderData();

    const columns = [
        CommanderColumn,
        ...StatColumns,
        CreatedAtColumn,
    ];

    return (
        <div id="decks" style={{height: 500, width: "75%"}}>
            <DataGrid
                rows={decks}
                columns={columns}
                slots={{toolbar: GridToolbar}}
                initialState={{
                    sorting: {
                        sortModel: [{field: "points", sort: "desc"}],
                    },
                }}
            />
        </div>
    );
}

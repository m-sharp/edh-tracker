import { ReactElement, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridPaginationModel, GridToolbar } from "@mui/x-data-grid";

import { GetDecksForPod } from "../../http";
import { StatColumns } from "../../components/stats";
import { Deck, PaginatedResponse } from "../../types";

interface PodDecksTabProps {
    decks: PaginatedResponse<Deck>;
    podId: number;
}

export default function PodDecksTab({ decks: initialData, podId }: PodDecksTabProps): ReactElement {
    const [rows, setRows] = useState<Deck[]>(initialData.items);
    const [rowCount, setRowCount] = useState(initialData.total);
    const [loading, setLoading] = useState(false);
    const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({ page: 0, pageSize: 25 });

    const handlePaginationChange = async (model: GridPaginationModel) => {
        setPaginationModel(model);
        setLoading(true);
        const data = await GetDecksForPod(podId, model.pageSize, model.page * model.pageSize);
        setRows(data.items);
        setRowCount(data.total);
        setLoading(false);
    };

    const columns: GridColDef[] = [
        {
            field: "name",
            headerName: "Deck",
            flex: 1,
            minWidth: 200,
            renderCell: (params) => (
                <Link to={`/player/${params.row.player_id}/deck/${params.row.id}`}>
                    {params.row.name}
                    {params.row.commanders && ` (${params.row.commanders.commander_name})`}
                </Link>
            ),
        },
        { field: "format_name", headerName: "Format", width: 120 },
        ...StatColumns,
    ];

    return (
        <Box sx={{ height: { xs: 400, sm: 600 }, width: "100%" }}>
            <DataGrid
                rows={rows}
                columns={columns}
                loading={loading}
                paginationMode="server"
                rowCount={rowCount}
                paginationModel={paginationModel}
                onPaginationModelChange={handlePaginationChange}
                slots={{ toolbar: GridToolbar }}
            />
        </Box>
    );
}

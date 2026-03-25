import { ReactElement, useState } from "react";
import { Link } from "react-router-dom";
import { Box } from "@mui/material";
import { DataGrid, GridColDef, GridPaginationModel, GridToolbar } from "@mui/x-data-grid";

import { GetGamesForPod } from "../../http";
import { Game, PaginatedResponse } from "../../types";

interface PodGamesTabProps {
    games: PaginatedResponse<Game>;
    podId: number;
}

export default function PodGamesTab({ games: initialData, podId }: PodGamesTabProps): ReactElement {
    const [rows, setRows] = useState<Game[]>(initialData.items);
    const [rowCount, setRowCount] = useState(initialData.total);
    const [loading, setLoading] = useState(false);
    const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({ page: 0, pageSize: 25 });

    const handlePaginationChange = async (model: GridPaginationModel) => {
        setPaginationModel(model);
        setLoading(true);
        const data = await GetGamesForPod(podId, model.pageSize, model.page * model.pageSize);
        setRows(data.items);
        setRowCount(data.total);
        setLoading(false);
    };

    const columns: GridColDef[] = [
        {
            field: "id",
            headerName: "Game #",
            width: 100,
            renderCell: (params) => (
                <Link to={`/pod/${podId}/game/${params.row.id}`}>#{params.row.id}</Link>
            ),
        },
        { field: "description", headerName: "Description", flex: 1, minWidth: 200 },
        {
            field: "created_at",
            headerName: "Date",
            width: 180,
            valueFormatter: (params) => new Date(params.value).toLocaleDateString(),
        },
        {
            field: "results",
            headerName: "Participants",
            width: 120,
            valueGetter: (params) => params.row.results?.length ?? 0,
        },
    ];

    return (
        <Box>
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
        </Box>
    );
}

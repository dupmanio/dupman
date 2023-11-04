import * as React from "react";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import useSWR, { useSWRConfig } from "swr";

import { Button, Tooltip, IconButton, Box } from "@mui/material";
import {
  DataGrid,
  GridColDef,
  GridToolbarContainer,
  GridValueGetterParams,
  GridRenderCellParams,
} from "@mui/x-data-grid";

import AddIcon from "@mui/icons-material/Add";
import VisibilityIcon from "@mui/icons-material/Visibility";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

import PageLoader from "@/components/PageLoader";
import WebsiteFormDialog from "@/components/WebsiteFormDialog";
import WebsiteDeleteDialog from "@/components/WebsiteDeleteDialog";
import WebsiteStatusCell from "@/components/WebsiteStatusCell";
import { WebsiteRepository } from "@/lib/repositories/website";
import { formatISO } from "@/lib/util/time";
import { Website } from "@/types/entities/website";

function Websites() {
  const router = useRouter();

  const [rowCount, setRowCount] = useState<number>(0);
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 10,
    page: 0,
  });
  const [formDialogOpen, setFormDialogOpen] = useState<boolean>(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);
  const [currentWebsite, setCurrentWebsite] = useState<Website | null>(null);

  const { mutate } = useSWRConfig();

  const { data, isLoading } = useSWR(
    ["/website", paginationModel.page, paginationModel.pageSize],
    ([_, page, limit]) => WebsiteRepository.getAll(page + 1, limit),
  );

  useEffect(() => {
    setRowCount((prevRowCount) =>
      data?.pagination?.totalItems !== undefined
        ? data?.pagination?.totalItems
        : prevRowCount,
    );
  }, [data, setRowCount]);

  function DataGridToolbar() {
    const handleFormDialogClose = () => {
      setFormDialogOpen(false);
      void mutate(["/website", paginationModel.page, paginationModel.pageSize]);
    };

    return (
      <>
        <GridToolbarContainer>
          <Button
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => {
              setCurrentWebsite(null);
              setFormDialogOpen(true);
            }}
          >
            Add new
          </Button>
        </GridToolbarContainer>

        <WebsiteFormDialog
          website={currentWebsite}
          open={formDialogOpen}
          onClose={handleFormDialogClose}
        />
      </>
    );
  }

  const columns: GridColDef[] = [
    {
      field: "url",
      headerName: "URL",
      width: 300,
    },
    {
      field: "lastScanDate",
      headerName: "Last Scan Date",
      width: 180,
      valueGetter: (params: GridValueGetterParams) =>
        formatISO(params.row.status.updatedAt),
    },
    {
      field: "status",
      headerName: "Status",
      width: 70,
      sortable: false,
      filterable: false,
      disableColumnMenu: true,
      renderCell: (params: GridRenderCellParams) => (
        <WebsiteStatusCell status={params.row.status} />
      ),
    },
    {
      field: "actions",
      headerName: "Actions",
      width: 130,
      sortable: false,
      filterable: false,
      hideable: false,
      disableColumnMenu: true,
      renderCell: (params: GridRenderCellParams) => (
        <Box>
          <Tooltip title="View">
            <IconButton
              color="info"
              onClick={() => {
                router.push({
                  pathname: "/website/[id]",
                  query: { id: params.row.id },
                });
              }}
            >
              <VisibilityIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Edit">
            <IconButton
              color="info"
              onClick={() => {
                setCurrentWebsite(params.row);
                setFormDialogOpen(true);
              }}
            >
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete">
            <IconButton
              color="error"
              onClick={() => {
                setCurrentWebsite(params.row);
                setDeleteDialogOpen(true);
              }}
            >
              <DeleteIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  return (
    <>
      <WebsiteDeleteDialog
        website={currentWebsite}
        open={deleteDialogOpen}
        onClose={() => {
          setDeleteDialogOpen(false);
          void mutate([
            "/website",
            paginationModel.page,
            paginationModel.pageSize,
          ]);
        }}
      />

      {isLoading && <PageLoader size={40} />}

      {!isLoading && data && (
        <DataGrid
          autoHeight={true}
          rows={data?.data ?? []}
          rowCount={rowCount}
          columns={columns}
          paginationMode="server"
          paginationModel={paginationModel}
          onPaginationModelChange={setPaginationModel}
          pageSizeOptions={[5, 10, 50]}
          slots={{
            toolbar: DataGridToolbar,
          }}
        />
      )}
    </>
  );
}

export default Websites;

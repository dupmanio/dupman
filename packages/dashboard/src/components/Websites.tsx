import * as React from "react";
import { useEffect, useState } from "react";
import { format, parseISO } from "date-fns";
import useSWR, { useSWRConfig } from "swr";

import { Button } from "@mui/material";
import {
  DataGrid,
  GridColDef,
  GridToolbarContainer,
  GridValueGetterParams,
  GridRenderCellParams,
} from "@mui/x-data-grid";

import AddIcon from "@mui/icons-material/Add";

import PageLoader from "@/components/PageLoader";
import WebsiteFormDialog from "@/components/WebsiteFormDialog";
import WebsiteStatusCell from "@/components/WebsiteStatusCell";
import { WebsiteRepository } from "@/lib/repositories/website";

function Websites() {
  const [rowCount, setRowCount] = useState<number>(0);
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 10,
    page: 0,
  });

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
    const [open, setOpen] = useState<boolean>(false);

    const { mutate } = useSWRConfig();

    const handleClose = () => {
      setOpen(false);
      void mutate(["/website", paginationModel.page, paginationModel.pageSize]);
    };

    return (
      <>
        <GridToolbarContainer>
          <Button
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => setOpen(true)}
          >
            Add new
          </Button>
        </GridToolbarContainer>

        <WebsiteFormDialog open={open} onClose={handleClose} />
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
        format(parseISO(params.row.status.updatedAt), "dd/MM/yyyy HH:mm:ss"),
    },
    {
      field: "status",
      headerName: "Status",
      renderCell: (params: GridRenderCellParams) => (
        <WebsiteStatusCell status={params.row.status} />
      ),
    },
  ];

  return (
    <>
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

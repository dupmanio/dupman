import * as React from "react";
import {
  DataGrid,
  GridColDef,
  GridToolbarContainer,
  GridValueGetterParams,
} from "@mui/x-data-grid";
import { useEffect, useState } from "react";
import useSWR, { useSWRConfig } from "swr";
import { WebsiteRepository } from "@/lib/repositories/website";
import { format, parseISO } from "date-fns";
import PageLoader from "@/components/PageLoader";
import { Button } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import WebsiteFormDialog from "@/components/WebsiteFormDialog";

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
    { field: "id", headerName: "ID", width: 300, sortable: false },
    {
      field: "url",
      headerName: "URL",
      width: 200,
      resizable: true,
    },
    {
      field: "createdAt",
      headerName: "Created At",
      width: 160,
      valueGetter: (params: GridValueGetterParams) =>
        format(parseISO(params.row.createdAt), "dd/MM/yyyy HH:mm:ss"),
    },
    {
      field: "updatedAt",
      headerName: "Updated At",
      width: 160,
      valueGetter: (params: GridValueGetterParams) =>
        format(parseISO(params.row.updatedAt), "dd/MM/yyyy HH:mm:ss"),
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

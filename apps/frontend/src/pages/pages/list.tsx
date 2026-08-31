import React, { useCallback, useEffect, useMemo, useState } from "react";
import { CrudFilter, HttpError, useGo } from "@refinedev/core";
import {
  DeleteButton,
  EditButton,
  List,
  useDataGrid,
} from "@refinedev/mui";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import debounce from "lodash/debounce";
import {
  Alert,
  Box,
  Button,
  Stack,
  Typography,
} from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import { PageListToolbar } from "./page-list-toolbar";
import { PageStatusChip } from "./page-status-chip";
import { IPageListItem, PageStatusFilter } from "./page-types";

const formatDate = (value: string): string => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString("vi-VN", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const EmptyPageOverlay: React.FC = () => {
  const go = useGo();

  return (
    <Stack alignItems="center" justifyContent="center" spacing={1} sx={{ height: "100%", p: 4 }}>
      <Typography variant="subtitle1" fontWeight={700} color="text.primary">
        Chưa có trang nào
      </Typography>
      <Typography variant="body2" color="text.secondary" align="center">
        Tạo trang đầu tiên để bắt đầu quản lý nội dung tĩnh.
      </Typography>
      <Button
        size="small"
        variant="outlined"
        onClick={() => go({ to: "/pages/create" })}
        sx={{ mt: 1, textTransform: "none", fontWeight: 700 }}
      >
        Tạo trang
      </Button>
    </Stack>
  );
};

export const PageList: React.FC = () => {
  const [searchValue, setSearchValue] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [status, setStatus] = useState<PageStatusFilter>("");

  const { dataGridProps, setFilters, tableQuery } = useDataGrid<IPageListItem, HttpError>({
    resource: "pages",
    syncWithLocation: true,
    pagination: { mode: "server", pageSize: 10 },
    sorters: {
      initial: [{ field: "updated_at", order: "desc" }],
    },
  });

  const debouncedSetSearch = useMemo(
    () => debounce((value: string) => setDebouncedSearch(value.trim()), 300),
    [],
  );

  useEffect(() => () => debouncedSetSearch.cancel(), [debouncedSetSearch]);

  useEffect(() => {
    const nextFilters: CrudFilter[] = [];
    if (debouncedSearch) {
      nextFilters.push({ field: "title_like", operator: "eq", value: debouncedSearch });
    }
    if (status) {
      nextFilters.push({ field: "status", operator: "eq", value: status });
    }
    setFilters(nextFilters);
  }, [debouncedSearch, setFilters, status]);

  const handleSearchChange = useCallback((value: string) => {
    setSearchValue(value);
    debouncedSetSearch(value);
  }, [debouncedSetSearch]);

  const handleStatusChange = useCallback((value: PageStatusFilter) => {
    setStatus(value);
  }, []);

  const columns = useMemo<GridColDef<IPageListItem>[]>(
    () => [
      {
        field: "id",
        headerName: "ID",
        width: 76,
        type: "number",
      },
      {
        field: "title",
        headerName: "Tiêu đề",
        minWidth: 240,
        flex: 1.25,
        renderCell: ({ row }) => (
          <Typography variant="body2" fontWeight={700} color="text.primary" noWrap title={row.title}>
            {row.title}
          </Typography>
        ),
      },
      {
        field: "slug",
        headerName: "Slug",
        minWidth: 220,
        flex: 1,
        renderCell: ({ row }) => (
          <Typography variant="body2" color="text.secondary" noWrap title={row.slug}>
            /{row.slug}
          </Typography>
        ),
      },
      {
        field: "status",
        headerName: "Trạng thái",
        width: 150,
        renderCell: ({ row }) => <PageStatusChip status={row.status} />,
      },
      {
        field: "updatedAt",
        headerName: "Cập nhật lúc",
        width: 180,
        valueGetter: (_, row) => row.updatedAt,
        renderCell: ({ value }) => (
          <Typography variant="body2" color="text.secondary">
            {formatDate(String(value))}
          </Typography>
        ),
      },
      {
        field: "actions",
        headerName: "Thao tác",
        width: 126,
        sortable: false,
        align: "right",
        headerAlign: "right",
        renderCell: ({ row }) => (
          <Stack direction="row" spacing={1} alignItems="center" justifyContent="flex-end" width="100%">
            <EditButton
              hideText
              size="small"
              recordItemId={row.id}
              sx={{ color: "text.secondary", bgcolor: "action.hover" }}
            />
            <DeleteButton
              hideText
              size="small"
              recordItemId={row.id}
              confirmTitle="Xóa trang này?"
              confirmOkText="Xóa"
              confirmCancelText="Hủy"
              sx={{ color: "error.main", bgcolor: "action.hover" }}
            />
          </Stack>
        ),
      },
    ],
    [],
  );

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight={700} color="text.primary" letterSpacing="-0.02em">
          Danh sách Trang
        </Typography>
      </Box>

      {tableQuery.isError && (
        <Alert
          severity="error"
          sx={{ mb: 2 }}
          action={(
            <Button color="inherit" size="small" startIcon={<RefreshIcon />} onClick={() => void tableQuery.refetch()}>
              Thử lại
            </Button>
          )}
        >
          Không thể tải danh sách trang. Vui lòng thử lại.
        </Alert>
      )}

      <List
        title=""
        createButtonProps={{
          onClick: (event) => {
            event.preventDefault();
            window.location.href = "/pages/create";
          },
        }}
        wrapperProps={{
          sx: {
            bgcolor: "background.paper",
            borderRadius: 2,
            border: "1px solid",
            borderColor: "divider",
            boxShadow: "none",
            p: 0,
            overflow: "hidden",
          },
        }}
        headerProps={{
          sx: {
            p: 2,
            px: { xs: 2, md: 3 },
            borderBottom: "1px solid",
            borderColor: "divider",
          },
        }}
        headerButtons={(props) => (
          <PageListToolbar
            searchValue={searchValue}
            status={status}
            onSearchChange={handleSearchChange}
            onStatusChange={handleStatusChange}
            disabled={Boolean(dataGridProps.loading)}
          >
            {props.defaultButtons}
          </PageListToolbar>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          autoHeight
          rowHeight={68}
          density="compact"
          disableColumnMenu
          slots={{ noRowsOverlay: EmptyPageOverlay }}
          sx={{
            border: "none",
            "& .MuiDataGrid-columnHeaders": {
              bgcolor: "action.hover",
              color: "text.primary",
              fontWeight: "bold",
              borderBottom: "1px solid",
              borderColor: "divider",
              borderRadius: 0,
            },
            "& .MuiDataGrid-row": {
              borderBottom: "1px solid",
              borderColor: "divider",
              "&:hover": { bgcolor: "action.hover" },
            },
            "& .MuiDataGrid-cell": {
              borderBottom: "none",
              display: "flex",
              alignItems: "center",
            },
            "& .MuiDataGrid-footerContainer": {
              borderTop: "1px solid",
              borderColor: "divider",
            },
          }}
        />
      </List>
    </Box>
  );
};

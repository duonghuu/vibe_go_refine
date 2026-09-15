import { useEffect, useMemo, useState } from "react";
import { CrudFilter, HttpError, useGo } from "@refinedev/core";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { List, useDataGrid } from "@refinedev/mui";
import {
  Alert,
  Box,
  Button,
  Chip,
  InputAdornment,
  MenuItem,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import SettingsSuggestOutlinedIcon from "@mui/icons-material/SettingsSuggestOutlined";
import { useDebounce } from "../../hooks/use-debounce";
import { ImageContentStatusChip } from "./components/image-content-status-chip";
import { ImageContentTypeDeleteButton } from "./components/image-content-type-delete-dialog";
import {
  IImageContentType,
  IMAGE_CONTENT_FIELD_KEYS,
  fieldLabel,
} from "./image-content-types";

const formatDate = (value: string) =>
  value ? new Date(value).toLocaleDateString("vi-VN") : "-";
export const ImageContentTypeList = () => {
  const go = useGo();
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");
  const debouncedSearch = useDebounce(search, 300);
  const { dataGridProps, setFilters, tableQuery } = useDataGrid<
    IImageContentType,
    HttpError
  >({
    resource: "image-content-types",
    dataProviderName: "imageContent",
    syncWithLocation: true,
    pagination: { mode: "server", pageSize: 10 },
    sorters: { initial: [{ field: "sortOrder", order: "asc" }] },
  });
  useEffect(() => {
    const filters: CrudFilter[] = [];
    if (debouncedSearch)
      filters.push({
        field: "q",
        operator: "contains",
        value: debouncedSearch,
      });
    if (status)
      filters.push({ field: "status", operator: "eq", value: status });
    setFilters(filters);
  }, [debouncedSearch, setFilters, status]);
  const columns = useMemo<GridColDef<IImageContentType>[]>(
    () => [
      {
        field: "code",
        headerName: "Code",
        minWidth: 140,
        flex: 0.7,
        renderCell: ({ row }) => (
          <Chip
            size="small"
            label={row.code}
            variant="outlined"
            sx={{ fontFamily: "monospace" }}
          />
        ),
      },
      { field: "name", headerName: "Tên loại", minWidth: 180, flex: 1 },
      {
        field: "fieldConfig",
        headerName: "Metadata",
        minWidth: 260,
        flex: 1.4,
        sortable: false,
        renderCell: ({ row }) => (
          <Stack
            direction="row"
            gap={0.5}
            flexWrap="wrap"
            sx={{
              height: "100%",
              alignItems: "center",
              alignContent: "center",
            }}
          >
            {IMAGE_CONTENT_FIELD_KEYS.filter(
              (key) => row.fieldConfig[key].enabled,
            ).map((key) => (
              <Tooltip
                key={key}
                title={row.fieldConfig[key].required ? "Bắt buộc" : "Tùy chọn"}
              >
                <Chip
                  size="small"
                  label={`${fieldLabel[key]} · ${row.fieldConfig[key].required ? "bắt buộc" : "tùy chọn"}`}
                />
              </Tooltip>
            ))}
          </Stack>
        ),
      },
      {
        field: "itemCount",
        headerName: "Số lượng",
        width: 120,
        renderCell: ({ row }) => (
          <Typography sx={{
              height: "100%",
              alignItems: "center",
              alignContent: "center",
            }}>
            {row.itemCount} / {row.maxItems ?? "∞"}
          </Typography>
        ),
      },
      {
        field: "status",
        headerName: "Trạng thái",
        width: 150,
        renderCell: ({ row }) => <ImageContentStatusChip status={row.status} />,
      },
      {
        field: "updatedAt",
        headerName: "Cập nhật",
        width: 120,
        renderCell: ({ row }) => formatDate(row.updatedAt),
      },
      {
        field: "actions",
        headerName: "Thao tác",
        minWidth: 220,
        sortable: false,
        align: "right",
        headerAlign: "right",
        renderCell: ({ row }) => (
          <Stack
            direction="row"
            spacing={0.5}
            justifyContent="flex-end"
            width="100%"
          >
            <Button
              size="small"
              startIcon={<SettingsSuggestOutlinedIcon />}
              onClick={() =>
                go({
                  to: `/image-contents?typeCode=${encodeURIComponent(row.code)}`,
                })
              }
            >
              Nội dung
            </Button>
            <Button
              size="small"
              onClick={() => go({ to: `/image-content-types/edit/${row.id}` })}
            >
              Sửa
            </Button>
            <ImageContentTypeDeleteButton id={row.id} />
          </Stack>
        ),
      },
    ],
    [go],
  );
  return (
    <Box>
      <Stack
        direction={{ xs: "column", sm: "row" }}
        justifyContent="space-between"
        alignItems={{ sm: "center" }}
        mb={3}
        gap={1}
      >
        <Box>
          <Typography variant="h4" fontWeight={700}>
            Loại nội dung hình ảnh
          </Typography>
          <Typography color="text.secondary">
            Cấu hình field hiển thị và giới hạn cho từng loại.
          </Typography>
        </Box>
        <Button
          variant="contained"
          onClick={() => go({ to: "/image-content-types/create" })}
        >
          Thêm loại nội dung
        </Button>
      </Stack>
      {tableQuery.isError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          Không thể tải danh sách.{" "}
          <Button color="inherit" onClick={() => void tableQuery.refetch()}>
            Thử lại
          </Button>
        </Alert>
      )}
      <List
        title=""
        wrapperProps={{
          sx: {
            p: 0,
            bgcolor: "background.paper",
            border: 1,
            borderColor: "divider",
            borderRadius: 2,
            overflow: "hidden",
          },
        }}
        headerProps={{ sx: { p: 2, borderBottom: 1, borderColor: "divider" } }}
        headerButtons={() => (
          <Stack
            direction={{ xs: "column", sm: "row" }}
            spacing={1}
            width="100%"
          >
            <TextField
              size="small"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder="Tìm code hoặc tên..."
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
              sx={{ minWidth: 260 }}
            />
            <TextField
              select
              size="small"
              value={status}
              onChange={(event) => setStatus(event.target.value)}
              sx={{ minWidth: 150 }}
            >
              <MenuItem value="">Tất cả trạng thái</MenuItem>
              <MenuItem value="ACTIVE">Đang hoạt động</MenuItem>
              <MenuItem value="INACTIVE">Tạm ẩn</MenuItem>
            </TextField>
          </Stack>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          autoHeight
          density="standard"
          disableRowSelectionOnClick
          sx={{
            border: 0,
            "& .MuiDataGrid-columnHeaders": { bgcolor: "action.hover" },
          }}
        />
      </List>
    </Box>
  );
};

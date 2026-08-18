import React, { useMemo, useState } from "react";
import { HttpError, CrudFilter } from "@refinedev/core";
import {
  List,
  useDataGrid,
  EditButton,
  DeleteButton,
} from "@refinedev/mui";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { useSearchParams } from "react-router";
import {
  Box,
  TextField,
  InputAdornment,
  Typography,
  Stack,
  Select,
  MenuItem,
  Avatar,
  Chip,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import FilterListIcon from "@mui/icons-material/FilterList";
import debounce from "lodash/debounce";
import { BACKEND_URL } from "../../providers/constants";

export interface IPostCategory {
  id: number;
  name: string;
  slug: string;
  parentId: number | null;
  description: string;
  imageUrl: string;
  sortOrder: number;
  status: "ACTIVE" | "INACTIVE";
  typeCode: string;
  createdAt: string;
}

export const PostCategoryList: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const typeCode = searchParams.get("typeCode") || searchParams.get("type_code") || "NEWS";

  const permanentFilters = useMemo<CrudFilter[]>(() => {
    return [
      {
        field: "typeCode",
        operator: "eq",
        value: typeCode,
      },
    ];
  }, [typeCode]);

  const { dataGridProps, setFilters } = useDataGrid<IPostCategory, HttpError>({
    resource: "post-categories",
    filters: {
      permanent: permanentFilters,
    },
    syncWithLocation: true,
  });

  const handleSearch = useMemo(
    () =>
      debounce((e: React.ChangeEvent<HTMLInputElement>) => {
        setFilters([
          {
            field: "q",
            operator: "eq",
            value: e.target.value,
          },
        ]);
      }, 300),
    [setFilters]
  );

  const columns = useMemo<GridColDef<IPostCategory>[]>(
    () => [
      {
        field: "imageUrl",
        headerName: "Hình ảnh",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Avatar
              src={row.imageUrl?.startsWith("/") ? `${BACKEND_URL}${row.imageUrl}` : row.imageUrl}
              variant="rounded"
              sx={{ width: 60, height: 60, my: 1, borderRadius: "8px" }}
            />
          );
        },
      },
      {
        field: "name",
        headerName: "Tên danh mục",
        flex: 1,
        minWidth: 200,
        renderCell: function render({ row }) {
          return (
            <Box sx={{ pl: row.parentId ? 3 : 0, display: "flex", alignItems: "center", height: "100%" }}>
              {row.parentId && <Typography sx={{ color: "text.secondary", mr: 1 }}>—</Typography>}
              <Typography fontWeight="600" color="text.primary">
                {row.name}
              </Typography>
            </Box>
          );
        },
      },
      {
        field: "slug",
        headerName: "Đường dẫn (Slug)",
        flex: 1,
        minWidth: 150,
      },
      {
        field: "typeCode",
        headerName: "Loại bài viết",
        width: 150,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="text.primary" sx={{ display: "flex", alignItems: "center", height: "100%" }}>
              {row.typeCode}
            </Typography>
          );
        },
      },
      {
        field: "sortOrder",
        headerName: "Thứ tự",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="text.primary" sx={{ display: "flex", alignItems: "center", height: "100%" }}>
              {row.sortOrder}
            </Typography>
          );
        },
      },
      {
        field: "status",
        headerName: "Trạng thái",
        width: 150,
        renderCell: function render({ row }) {
          const isActive = row.status === "ACTIVE";
          return (
            <Box sx={{ display: "flex", alignItems: "center", height: "100%" }}>
              <Chip
                label={isActive ? "Hoạt động" : "Ẩn"}
                size="small"
                sx={{
                  backgroundColor: isActive ? "#e8f5e9" : "#f5f5f5",
                  color: isActive ? "#2e7d32" : "#757575",
                  fontWeight: "bold",
                }}
              />
            </Box>
          );
        },
      },
      {
        field: "actions",
        headerName: "Hành động",
        sortable: false,
        width: 120,
        align: "right",
        headerAlign: "right",
        renderCell: function render({ row }) {
          return (
            <Stack direction="row" spacing={1} justifyContent="flex-end" alignItems="center" sx={{ height: "100%" }}>
              <EditButton
                hideText
                recordItemId={row.id}
                size="small"
                sx={{
                  bgcolor: "#f3f4f6",
                  color: "#6b7280",
                  "&:hover": { bgcolor: "#e5e7eb", color: "#374151" },
                }}
              />
              <DeleteButton
                hideText
                recordItemId={row.id}
                size="small"
                sx={{
                  bgcolor: "#f3f4f6",
                  color: "#ef4444",
                  "&:hover": { bgcolor: "#fef2f2", color: "#dc2626" },
                }}
              />
            </Stack>
          );
        },
      },
    ],
    []
  );

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Danh sách danh mục bài viết
        </Typography>
      </Box>

      <List
        title=""
        wrapperProps={{
          sx: {
            bgcolor: "#fff",
            borderRadius: "14px",
            border: "1px solid #D5D5D5",
            boxShadow: "none",
            p: 0,
            overflow: "hidden",
          },
        }}
        headerProps={{
          sx: {
            p: 2,
            px: 3,
            borderBottom: "1px solid #D5D5D5",
          },
        }}
        createButtonProps={{
          onClick: (e) => {
            e.preventDefault();
            window.location.href = `/post-categories/create?typeCode=${typeCode}`;
          },
        }}
        headerButtons={(props) => (
          <Stack direction={{ xs: "column", md: "row" }} spacing={2} alignItems="center">
            <TextField
              placeholder="Tìm kiếm danh mục..."
              variant="outlined"
              size="small"
              onChange={handleSearch}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon sx={{ color: "text.secondary" }} />
                  </InputAdornment>
                ),
                sx: { borderRadius: "50px", backgroundColor: "#F5F6FA", width: { xs: "100%", md: "250px" } },
              }}
            />

            <Select
              size="small"
              displayEmpty
              defaultValue=""
              onChange={(e) => {
                setFilters([
                  {
                    field: "status",
                    operator: "eq",
                    value: e.target.value,
                  },
                ]);
              }}
              sx={{ borderRadius: "8px", minWidth: "150px", backgroundColor: "#F5F6FA" }}
              IconComponent={FilterListIcon}
            >
              <MenuItem value="">Tất cả trạng thái</MenuItem>
              <MenuItem value="ACTIVE">Hoạt động</MenuItem>
              <MenuItem value="INACTIVE">Ẩn</MenuItem>
            </Select>
            {props.defaultButtons}
          </Stack>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          autoHeight
          rowHeight={80}
          density="standard"
          disableColumnMenu
          checkboxSelection
          sx={{
            border: "none",
            "& .MuiDataGrid-columnHeaders": {
              bgcolor: "#F5F6FA",
              color: "#202224",
              fontWeight: "bold",
              borderBottom: "1px solid #D5D5D5",
              borderTop: "none",
              borderRadius: 0,
            },
            "& .MuiDataGrid-row": {
              borderBottom: "1px solid #f3f4f6",
              "&:hover": {
                bgcolor: "rgba(249, 250, 251, 0.5)",
              },
            },
            "& .MuiDataGrid-cell": {
              borderBottom: "none",
              display: "flex",
              alignItems: "center",
            },
            "& .MuiDataGrid-footerContainer": {
              borderTop: "1px solid #D5D5D5",
            },
          }}
        />
      </List>
    </Box>
  );
};

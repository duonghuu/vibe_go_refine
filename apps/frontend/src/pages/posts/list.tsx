import React, { useMemo, useState, useEffect } from "react";
import { HttpError, CrudFilter } from "@refinedev/core";
import {
  List,
  useDataGrid,
  EditButton,
  ShowButton,
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
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import { IPostResponse } from "./create";

export const PostList: React.FC = () => {
  const [searchParams] = useSearchParams();
  const typeCode = searchParams.get("type_code") || "NEWS";

  const [searchTerm, setSearchTerm] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");

  // Debounce logic (300ms)
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(searchTerm);
    }, 300);
    return () => clearTimeout(handler);
  }, [searchTerm]);

  const permanentFilters = useMemo<CrudFilter[]>(() => {
    return [
      {
        field: "typeCode",
        operator: "eq",
        value: typeCode,
      },
    ];
  }, [typeCode]);

  const { dataGridProps, setFilters } = useDataGrid<IPostResponse, HttpError>({
    resource: "posts",
    filters: {
      permanent: permanentFilters,
    },
    syncWithLocation: true,
  });

  // Effect to update search filter when debounced value changes
  useEffect(() => {
    setFilters([
      {
        field: "title",
        operator: "contains",
        value: debouncedSearch || undefined,
      },
    ]);
  }, [debouncedSearch, setFilters]);

  const pageTitle = useMemo(() => {
    switch (typeCode) {
      case "NEWS":
        return "Danh sách Tin tức";
      case "PAGE":
        return "Danh sách Trang";
      default:
        return `Danh sách Bài viết (${typeCode})`;
    }
  }, [typeCode]);

  const columns = useMemo<GridColDef<IPostResponse>[]>(
    () => [
      {
        field: "id",
        headerName: "ID",
        width: 80,
      },
      {
        field: "title",
        headerName: "Tiêu đề",
        flex: 1,
        minWidth: 250,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" fontWeight="600" color="text.primary">
              {row.title}
            </Typography>
          );
        },
      },
      {
        field: "slug",
        headerName: "Đường dẫn",
        flex: 1,
        minWidth: 200,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" color="text.secondary">
              {row.slug}
            </Typography>
          );
        },
      },
      {
        field: "authorId",
        headerName: "Tác giả",
        width: 100,
        align: "center",
        headerAlign: "center",
      },
      {
        field: "createdAt",
        headerName: "Ngày tạo",
        width: 180,
        renderCell: function render({ value }) {
          if (!value) return "-";
          const date = new Date(value);
          return (
            <Typography variant="body2" color="text.secondary">
              {date.toLocaleDateString("vi-VN", {
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })}
            </Typography>
          );
        },
      },
      {
        field: "actions",
        headerName: "Hành động",
        sortable: false,
        width: 150,
        align: "center",
        headerAlign: "center",
        renderCell: function render({ row }) {
          return (
            <Stack direction="row" spacing={1} justifyContent="center">
              <ShowButton 
                hideText 
                recordItemId={row.id} 
                size="small"
                sx={{
                  bgcolor: "#f3f4f6",
                  color: "#6b7280",
                  "&:hover": { bgcolor: "#e5e7eb", color: "#374151" }
                }} 
              />
              <EditButton 
                hideText 
                recordItemId={row.id} 
                size="small"
                sx={{
                  bgcolor: "#f3f4f6",
                  color: "#6b7280",
                  "&:hover": { bgcolor: "#e5e7eb", color: "#374151" }
                }} 
              />
              <DeleteButton 
                hideText 
                recordItemId={row.id} 
                size="small"
                sx={{
                  bgcolor: "#f3f4f6",
                  color: "#ef4444",
                  "&:hover": { bgcolor: "#fef2f2", color: "#dc2626" }
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
          {pageTitle}
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
            window.location.href = `/posts/create?type_code=${typeCode}`;
          }
        }}
        headerButtons={(props) => (
          <Stack direction="row" spacing={2} alignItems="center">
            <TextField
              placeholder="Tìm kiếm tiêu đề..."
              variant="outlined"
              size="small"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon sx={{ color: 'text.secondary' }} />
                  </InputAdornment>
                ),
                sx: { borderRadius: '50px', backgroundColor: '#F5F6FA', width: { xs: '100%', md: '300px' } }
              }}
            />
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
            }
          }}
        />
      </List>
    </Box>
  );
};

import React, { useMemo } from "react";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { DeleteButton, EditButton, List, useDataGrid } from "@refinedev/mui";
import { HttpError } from "@refinedev/core";
import debounce from "lodash/debounce";
import { Box, InputAdornment, MenuItem, Select, Stack, TextField, Typography } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import FilterListIcon from "@mui/icons-material/FilterList";
import { PostTypeStatusBadge } from "./components/PostTypeStatusBadge";

export interface IPostType {
  id: number;
  code: string;
  name: string;
  status: "ACTIVE" | "INACTIVE";
  sort_order: number;
}

export const PostTypeList = () => {
  const { dataGridProps, setFilters } = useDataGrid<IPostType, HttpError>({
    resource: "post-types",
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

  const columns = useMemo<GridColDef[]>(
    () => [
      {
        field: "id",
        headerName: "ID",
        width: 80,
      },
      {
        field: "code",
        headerName: "Code",
        minWidth: 150,
        flex: 1,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="primary" sx={{ 
              display: 'inline-flex', 
              alignItems: 'center' 
            }}>
              {row.code}
            </Typography>
          );
        }
      },
      {
        field: "name",
        headerName: "Tên loại",
        flex: 1,
        minWidth: 200,
      },
      {
        field: "sort_order",
        headerName: "Thứ tự",
        width: 120,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="text.primary" sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              {row.sort_order}
            </Typography>
          );
        },
      },
      {
        field: "status",
        headerName: "Trạng thái",
        width: 150,
        renderCell: function render({ row }) {
          return (
            <Box sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              <PostTypeStatusBadge status={row.status} />
            </Box>
          );
        },
      },
      {
        field: "actions",
        headerName: "Thao tác",
        align: "right",
        headerAlign: "right",
        minWidth: 120,
        sortable: false,
        renderCell: function render({ row }) {
          return (
            <Stack direction="row" spacing={1} alignItems="center" justifyContent="flex-end" sx={{ height: '100%' }}>
              <EditButton
                hideText
                recordItemId={row.id}
                sx={{
                  backgroundColor: '#f3f4f6',
                  color: '#6b7280',
                  '&:hover': { backgroundColor: '#e5e7eb', color: '#374151' }
                }}
              />
              <DeleteButton
                hideText
                recordItemId={row.id}
                sx={{
                  backgroundColor: '#f3f4f6',
                  color: '#ef4444',
                  '&:hover': { backgroundColor: '#fef2f2', color: '#dc2626' }
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
          Quản lý loại bài viết
        </Typography>
      </Box>
      <List
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
        title=""
        headerButtons={(props) => (
          <Stack direction="row" spacing={2}>
            <TextField
              placeholder="Tìm Code hoặc Tên loại..."
              variant="outlined"
              size="small"
              onChange={handleSearch}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon sx={{ color: 'text.secondary' }} />
                  </InputAdornment>
                ),
                sx: { borderRadius: '50px', backgroundColor: '#F5F6FA', width: { xs: '100%', md: '300px' } }
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
              sx={{ borderRadius: '8px', minWidth: '150px', backgroundColor: '#F5F6FA' }}
              IconComponent={FilterListIcon}
            >
              <MenuItem value="">Tất cả trạng thái</MenuItem>
              <MenuItem value="ACTIVE">Hoạt động</MenuItem>
              <MenuItem value="INACTIVE">Tạm ẩn</MenuItem>
            </Select>
            {props.defaultButtons}
          </Stack>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          density="standard"
          sx={{
            border: 'none',
            '& .MuiDataGrid-columnHeaders': {
              backgroundColor: '#F5F6FA',
              borderBottom: '1px solid #D5D5D5',
              color: 'text.primary',
              fontWeight: 'bold',
            },
            '& .MuiDataGrid-cell': {
              borderBottom: '1px solid #f3f4f6',
            },
            '& .MuiDataGrid-row:hover': {
              backgroundColor: '#f9fafb',
            }
          }}
        />
      </List>
    </Box>
  );
};

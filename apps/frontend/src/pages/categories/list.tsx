import React, { useMemo } from "react";
import {
  DataGrid,
  GridColDef,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import {
  DeleteButton,
  EditButton,
  List,
  useDataGrid,
} from "@refinedev/mui";
import { HttpError } from "@refinedev/core";
import debounce from "lodash/debounce";
import {
  Box,
  Card,
  CardContent,
  Grid,
  InputAdornment,
  MenuItem,
  Select,
  TextField,
  Typography,
  Chip,
  IconButton,
  Avatar,
  Stack,
  Button
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import FilterListIcon from "@mui/icons-material/FilterList";
import ReorderIcon from "@mui/icons-material/Reorder";
import { BACKEND_URL } from "../../providers/constants";

export interface ICategory {
  id: number;
  name: string;
  slug: string;
  parentId: number | null;
  imageUrl: string;
  sortOrder: number;
  status: "ACTIVE" | "HIDDEN";
  productCount: number;
  createdAt: string;
}

export const CategoryList = () => {
  const { dataGridProps, setFilters } = useDataGrid<ICategory, HttpError>({
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
        field: "imageUrl",
        headerName: "Image",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Avatar
              src={row.imageUrl?.startsWith("/") ? `${BACKEND_URL}${row.imageUrl}` : row.imageUrl}
              variant="rounded"
              sx={{ width: 60, height: 60, my: 1, borderRadius: '8px' }}
            />
          );
        },
      },
      {
        field: "name",
        headerName: "Category Name",
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
        headerName: "Slug",
        flex: 1,
        minWidth: 150,
      },
      {
        field: "productCount",
        headerName: "Products",
        width: 120,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="text.primary" sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              {row.productCount}
            </Typography>
          );
        },
      },
      {
        field: "sortOrder",
        headerName: "Order",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Typography fontWeight="600" color="text.primary" sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              {row.sortOrder}
            </Typography>
          );
        },
      },
      {
        field: "status",
        headerName: "Status",
        width: 150,
        renderCell: function render({ row }) {
          const isActive = row.status === "ACTIVE";
          return (
            <Box sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
              <Chip
                label={row.status === "ACTIVE" ? "Active" : "Hidden"}
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
        headerName: "Action",
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
          Category List
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
              placeholder="Search category name"
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
              <MenuItem value="">All Status</MenuItem>
              <MenuItem value="ACTIVE">Active</MenuItem>
              <MenuItem value="HIDDEN">Hidden</MenuItem>
            </Select>
            {props.defaultButtons}
          </Stack>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          rowHeight={80}
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

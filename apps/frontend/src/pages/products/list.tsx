import React, { useMemo, useState } from "react";
import { HttpError } from "@refinedev/core";
import {
  List,
  useDataGrid,
  EditButton,
  DeleteButton,
} from "@refinedev/mui";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import {
  Box,
  Typography,
  InputBase,
  Avatar,
  Stack,
  InputAdornment,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import debounce from "@mui/material/utils/debounce";

export interface ICategory {
  id: number;
  name: string;
  parentId?: number;
}

export interface IProduct {
  id: number;
  sku: string;
  name: string;
  categoryId: number;
  category: ICategory;
  price: number;
  salePrice?: number;
  stock: number;
  soldCount: number;
  status: string;
  image: string;
  createdAt: string;
  updatedAt: string;
}

export const ProductList: React.FC = () => {
  const [searchValue, setSearchValue] = useState("");

  const { dataGridProps, search, filters } = useDataGrid<IProduct, HttpError>({
    resource: "products",
    onSearch: (params: any) => {
      return [
        {
          field: "q",
          operator: "eq",
          value: params.q,
        },
      ];
    },
  });

  // Debounce search function sử dụng debounce của MUI
  const handleSearch = useMemo(
    () =>
      debounce((value: string) => {
        search([{ field: "q", operator: "eq", value }]);
      }, 300),
    [search]
  );

  const onSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchValue(e.target.value);
    handleSearch(e.target.value);
  };

  const columns = useMemo<GridColDef<IProduct>[]>(
    () => [
      {
        field: "image",
        headerName: "Image",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Avatar
              variant="rounded"
              src={row.image}
              alt={row.name}
              sx={{ width: 60, height: 60, bgcolor: "#e2e8f0" }}
            />
          );
        },
      },
      {
        field: "sku",
        headerName: "SKU",
        width: 120,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" color="text.secondary">
              {row.sku}
            </Typography>
          );
        },
      },
      {
        field: "name",
        headerName: "Product Name",
        flex: 1,
        minWidth: 200,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" fontWeight="600" color="text.primary">
              {row.name}
            </Typography>
          );
        },
      },
      {
        field: "category",
        headerName: "Category",
        width: 150,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" fontWeight="600" color="text.primary">
              {row.category?.name || "N/A"}
            </Typography>
          );
        },
      },
      {
        field: "price",
        headerName: "Price",
        width: 120,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" fontWeight="600" color="text.primary">
              ${Number(row.price).toFixed(2)}
            </Typography>
          );
        },
      },
      {
        field: "stock",
        headerName: "Stock",
        width: 100,
        renderCell: function render({ row }) {
          return (
            <Typography variant="body2" fontWeight="600" color="text.primary">
              {row.stock}
            </Typography>
          );
        },
      },
      {
        field: "status",
        headerName: "Status",
        width: 130,
        renderCell: function render({ row }) {
          const isActive = row.status === "ACTIVE";
          return (
            <Box
              sx={{
                px: 1.5,
                py: 0.5,
                borderRadius: "9999px",
                bgcolor: isActive ? "#dcfce7" : "#fee2e2",
                color: isActive ? "#166534" : "#991b1b",
                fontSize: "12px",
                fontWeight: "600",
                display: "inline-block"
              }}
            >
              {row.status}
            </Box>
          );
        },
      },
      {
        field: "actions",
        headerName: "Action",
        width: 120,
        align: "center",
        headerAlign: "center",
        sortable: false,
        renderCell: function render({ row }) {
          return (
            <Stack direction="row" spacing={1} justifyContent="center">
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
          Product Stock
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
        headerButtons={
          <Box sx={{ position: "relative", width: "260px" }}>
            <InputBase
              placeholder="Search product name"
              value={searchValue}
              onChange={onSearchChange}
              startAdornment={
                <InputAdornment position="start">
                  <SearchIcon sx={{ color: "#9ca3af", fontSize: 20 }} />
                </InputAdornment>
              }
              sx={{
                bgcolor: "#F5F6FA",
                borderRadius: "9999px",
                px: 2,
                py: 0.5,
                width: "100%",
                fontSize: "14px",
                border: "1px solid #D5D5D5",
                "&:focus-within": {
                  borderColor: "primary.main",
                  boxShadow: "0 0 0 2px rgba(72, 128, 255, 0.2)",
                },
              }}
            />
          </Box>
        }
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

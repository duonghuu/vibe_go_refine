import React, { useMemo, useState } from "react";
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

// Dữ liệu giả (Mock Data) tham chiếu từ .docs/mock-data/products.json
const mockProducts: IProduct[] = [
  {
    id: 1,
    name: "Apple Watch Series 4",
    category: { id: 1, name: "Digital Product" },
    price: 690.00,
    stock: 63,
    colors: ["#000000", "#D1D5DB", "#FBCFE8"],
    image: "https://images.unsplash.com/photo-1546868871-7041f2a55e12?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80",
    status: "Active"
  },
  {
    id: 2,
    name: "Microsoft Headsquare",
    category: { id: 1, name: "Digital Product" },
    price: 190.00,
    stock: 13,
    colors: ["#000000", "#3B82F6", "#EF4444"],
    image: "https://images.unsplash.com/photo-1550009158-9ebf69173e03?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80",
    status: "Active"
  },
  {
    id: 3,
    name: "Women's Dress",
    category: { id: 2, name: "Fashion" },
    price: 640.00,
    stock: 635,
    colors: ["#A855F7", "#FBBF24", "#000000"],
    image: "https://images.unsplash.com/photo-1595777457583-95e059d581b8?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80",
    status: "Active"
  },
  {
    id: 4,
    name: "Samsung A50",
    category: { id: 3, name: "Mobile" },
    price: 400.00,
    stock: 67,
    colors: ["#000000", "#1D4ED8", "#F3F4F6"],
    image: "https://images.unsplash.com/photo-1598327105666-5b89351cb315?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80",
    status: "Active"
  },
  {
    id: 5,
    name: "Camera",
    category: { id: 4, name: "Electronic" },
    price: 420.00,
    stock: 52,
    colors: ["#000000", "#9CA3AF", "#DC2626"],
    image: "https://images.unsplash.com/photo-1516035069371-29a1b244cc32?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80",
    status: "Active"
  }
];

export interface ICategory {
  id: number;
  name: string;
}

export interface IProduct {
  id: number;
  name: string;
  category: ICategory;
  price: number;
  stock: number;
  colors: string[];
  image: string;
  status: string;
}

export const ProductList: React.FC = () => {
  const [searchValue, setSearchValue] = useState("");

  const { dataGridProps, search, filters } = useDataGrid<IProduct>({
    onSearch: (params: any) => {
      return [
        {
          field: "name",
          operator: "contains",
          value: params.name,
        },
      ];
    },
  });

  // Debounce search function sử dụng debounce của MUI
  const handleSearch = useMemo(
    () =>
      debounce((value: string) => {
        search([{ field: "name", operator: "contains", value }]);
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
        headerName: "Piece",
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
        field: "colors",
        headerName: "Available Color",
        width: 150,
        renderCell: function render({ row }) {
          return (
            <Stack direction="row" spacing={-1}>
              {row.colors?.map((color, idx) => (
                <Box
                  key={idx}
                  sx={{
                    width: 20,
                    height: 20,
                    borderRadius: "50%",
                    bgcolor: color,
                    border: "2px solid #fff",
                    boxShadow: "0 0 2px rgba(0,0,0,0.1)",
                  }}
                />
              ))}
            </Stack>
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
          rows={mockProducts} // Sử dụng mock data thay cho dataGridProps.rows để render đúng UI mockup
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

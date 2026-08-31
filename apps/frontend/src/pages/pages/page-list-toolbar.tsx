import React from "react";
import {
  FormControl,
  InputAdornment,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
} from "@mui/material";
import FilterListIcon from "@mui/icons-material/FilterList";
import SearchIcon from "@mui/icons-material/Search";
import { PageStatusFilter } from "./page-types";

interface PageListToolbarProps {
  searchValue: string;
  status: PageStatusFilter;
  onSearchChange: (value: string) => void;
  onStatusChange: (value: PageStatusFilter) => void;
  disabled?: boolean;
  children?: React.ReactNode;
}

export const PageListToolbar: React.FC<PageListToolbarProps> = ({
  searchValue,
  status,
  onSearchChange,
  onStatusChange,
  disabled = false,
  children,
}) => (
  <Stack
    direction={{ xs: "column", sm: "row" }}
    spacing={2}
    alignItems={{ xs: "stretch", sm: "center" }}
    width="100%"
  >
    <TextField
      fullWidth
      placeholder="Tìm theo tiêu đề"
      variant="outlined"
      size="small"
      value={searchValue}
      onChange={(event) => onSearchChange(event.target.value)}
      disabled={disabled}
      inputProps={{ "aria-label": "Tìm kiếm trang" }}
      InputProps={{
        startAdornment: (
          <InputAdornment position="start">
            <SearchIcon color="action" />
          </InputAdornment>
        ),
      }}
      sx={{
        maxWidth: { sm: 340 },
        "& .MuiOutlinedInput-root": {
          borderRadius: 6,
          bgcolor: "action.hover",
        },
      }}
    />

    <FormControl size="small" sx={{ minWidth: { sm: 190 } }}>
      <InputLabel id="page-status-filter-label">Trạng thái</InputLabel>
      <Select
        labelId="page-status-filter-label"
        value={status}
        label="Trạng thái"
        onChange={(event) => onStatusChange(event.target.value as PageStatusFilter)}
        disabled={disabled}
        startAdornment={<FilterListIcon color="action" sx={{ mr: 1 }} />}
        sx={{ borderRadius: 2, bgcolor: "action.hover" }}
      >
        <MenuItem value="">Tất cả trạng thái</MenuItem>
        <MenuItem value="DRAFT">Bản nháp</MenuItem>
        <MenuItem value="PUBLISHED">Đã xuất bản</MenuItem>
      </Select>
    </FormControl>

    {children}
  </Stack>
);

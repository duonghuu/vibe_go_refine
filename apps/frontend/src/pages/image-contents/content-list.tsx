import { useEffect, useMemo, useState } from "react";
import { CrudFilter, HttpError, useGo, useList } from "@refinedev/core";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { List } from "@refinedev/mui";
import { Alert, Avatar, Box, Button, Chip, InputAdornment, MenuItem, Stack, TextField, Tooltip, Typography } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import SortIcon from "@mui/icons-material/Sort";
import AddIcon from "@mui/icons-material/Add";
import { useDataGrid } from "@refinedev/mui";
import { useDebounce } from "../../hooks/use-debounce";
import { ImageContentDeleteButton } from "./components/image-content-delete-dialog";
import { ImageContentReorderDialog } from "./components/image-content-reorder-dialog";
import { ImageContentStatusChip } from "./components/image-content-status-chip";
import { IImageContent, IImageContentType, mediaPreviewUrl } from "./image-content-types";

const date = (value: string) => value ? new Date(value).toLocaleDateString("vi-VN") : "-";
export const ImageContentList = () => {
  const go = useGo(); const params = new URLSearchParams(window.location.search); const queryType = params.get("typeCode") ?? "";
  const [search, setSearch] = useState(""); const [status, setStatus] = useState(""); const [typeCode, setTypeCode] = useState(queryType); const [reorderOpen, setReorderOpen] = useState(false); const debouncedSearch = useDebounce(search, 300);
  const types = useList<IImageContentType>({ resource: "image-content-types", dataProviderName: "imageContent", pagination: { currentPage: 1, pageSize: 100 } });
  const selectedType = types.result.data.find((type) => type.code === typeCode);
  const capacityReached = Boolean(selectedType && selectedType.maxItems !== null && selectedType.itemCount >= selectedType.maxItems);
  const { dataGridProps, setFilters, tableQuery } = useDataGrid<IImageContent, HttpError>({ resource: "image-contents", dataProviderName: "imageContent", syncWithLocation: true, pagination: { mode: "server", pageSize: 10 }, sorters: { initial: [{ field: "sortOrder", order: "asc" }] } });
  useEffect(() => { const filters: CrudFilter[] = []; if (debouncedSearch) filters.push({ field: "q", operator: "contains", value: debouncedSearch }); if (status) filters.push({ field: "status", operator: "eq", value: status }); if (typeCode) filters.push({ field: "typeCode", operator: "eq", value: typeCode }); setFilters(filters); }, [debouncedSearch, setFilters, status, typeCode]);
  const columns = useMemo<GridColDef<IImageContent>[]>(() => [
    { field: "media", headerName: "Ảnh", width: 92, sortable: false, renderCell: ({ row }) => <Avatar variant="rounded" src={mediaPreviewUrl(row.media)} alt={row.name || `${row.typeName} #${row.id}`} sx={{ width: 64, height: 44 }} /> },
    { field: "name", headerName: "Tên", minWidth: 180, flex: 1, renderCell: ({ row }) => <Typography fontWeight={600} noWrap title={row.name ?? undefined}>{row.name || `${row.typeName} #${row.id}`}</Typography> },
    { field: "typeCode", headerName: "Type", minWidth: 150, flex: 0.7, renderCell: ({ row }) => <Stack><Typography variant="body2">{row.typeName}</Typography><Typography variant="caption" color="text.secondary">{row.typeCode}</Typography></Stack> },
    { field: "description", headerName: "Mô tả", minWidth: 180, flex: 1, sortable: false, renderCell: ({ row }) => <Tooltip title={row.description ?? ""}><Typography noWrap color="text.secondary">{row.description || "-"}</Typography></Tooltip> },
    { field: "status", headerName: "Trạng thái", width: 145, renderCell: ({ row }) => <ImageContentStatusChip status={row.status} /> },
    { field: "sortOrder", headerName: "Vị trí", width: 80 }, { field: "updatedAt", headerName: "Cập nhật", width: 110, renderCell: ({ row }) => date(row.updatedAt) },
    { field: "actions", headerName: "Thao tác", width: 130, sortable: false, align: "right", headerAlign: "right", renderCell: ({ row }) => <Stack direction="row" justifyContent="flex-end" width="100%"><Button size="small" onClick={() => go({ to: `/image-contents/edit/${row.id}` })}>Sửa</Button><ImageContentDeleteButton id={row.id} /></Stack> },
  ], [go]);
  const capacity = selectedType ? `${selectedType.itemCount} / ${selectedType.maxItems ?? "Không giới hạn"}` : "Chọn Type để xem giới hạn";
  return <Box><Stack direction={{ xs: "column", sm: "row" }} justifyContent="space-between" alignItems={{ sm: "center" }} mb={3} gap={1}><Box><Typography variant="h4" fontWeight={700}>Nội dung hình ảnh</Typography><Typography color="text.secondary">Quản lý ảnh và metadata theo Type.</Typography></Box><Stack direction="row" spacing={1}><Button variant="outlined" startIcon={<SortIcon />} disabled={!selectedType || selectedType.itemCount < 2} onClick={() => setReorderOpen(true)}>Sắp xếp</Button><Button variant="contained" startIcon={<AddIcon />} disabled={capacityReached} onClick={() => go({ to: `/image-contents/create${typeCode ? `?typeCode=${encodeURIComponent(typeCode)}` : ""}` })}>Thêm nội dung</Button></Stack></Stack>{capacityReached && <Alert severity="warning" sx={{ mb: 2 }}>Type đã đạt giới hạn nội dung. Hãy tăng giới hạn hoặc xóa một item trước khi thêm mới.</Alert>}{tableQuery.isError && <Alert severity="error" sx={{ mb: 2 }}>Không thể tải danh sách. <Button color="inherit" onClick={() => void tableQuery.refetch()}>Thử lại</Button></Alert>}<List title="" wrapperProps={{ sx: { p: 0, bgcolor: "background.paper", border: 1, borderColor: "divider", borderRadius: 2, overflow: "hidden" } }} headerProps={{ sx: { p: 2, borderBottom: 1, borderColor: "divider" } }} headerButtons={() => <Stack direction={{ xs: "column", md: "row" }} spacing={1} width="100%"><TextField size="small" value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Tìm metadata..." InputProps={{ startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment> }} sx={{ minWidth: 240 }} /><TextField select size="small" value={typeCode} onChange={(event) => setTypeCode(event.target.value)} label="Type" sx={{ minWidth: 180 }}><MenuItem value="">Tất cả Type</MenuItem>{types.result.data.map((type) => <MenuItem key={type.code} value={type.code}>{type.name} ({type.code})</MenuItem>)}</TextField><TextField select size="small" value={status} onChange={(event) => setStatus(event.target.value)} label="Trạng thái" sx={{ minWidth: 150 }}><MenuItem value="">Tất cả</MenuItem><MenuItem value="ACTIVE">Đang hoạt động</MenuItem><MenuItem value="INACTIVE">Tạm ẩn</MenuItem></TextField><Chip label={`Sức chứa: ${capacity}`} sx={{ alignSelf: "center" }} /></Stack>}><DataGrid {...dataGridProps} columns={columns} autoHeight density="standard" disableRowSelectionOnClick sx={{ border: 0, "& .MuiDataGrid-columnHeaders": { bgcolor: "action.hover" } }} /></List>{typeCode && <ImageContentReorderDialog open={reorderOpen} typeCode={typeCode} onClose={() => setReorderOpen(false)} onSaved={() => { void tableQuery.refetch(); void types.query.refetch(); }} />}</Box>;
};

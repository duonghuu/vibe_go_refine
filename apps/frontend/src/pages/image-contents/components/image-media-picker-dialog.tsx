import { useEffect, useMemo, useState } from "react";
import { useList } from "@refinedev/core";
import { Avatar, Box, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, List, ListItemButton, ListItemText, Pagination, Stack, TextField, Typography } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import { useDebounce } from "../../../hooks/use-debounce";
import { imageContentMediaFromUpload } from "../../../providers/image-content-data-provider";
import { IImageContentMedia, mediaPreviewUrl } from "../image-content-types";

interface Props { open: boolean; selected: IImageContentMedia | null; onClose: () => void; onSelect: (media: IImageContentMedia) => void; }

export const ImageMediaPickerDialog = ({ open, selected, onClose, onSelect }: Props) => {
  const [search, setSearch] = useState(""); const [page, setPage] = useState(1); const [uploading, setUploading] = useState(false); const debouncedSearch = useDebounce(search, 300);
  const list = useList<IImageContentMedia>({ resource: "media", dataProviderName: "imageContent", pagination: { currentPage: page, pageSize: 8 }, filters: [{ field: "q", operator: "contains", value: debouncedSearch }], queryOptions: { enabled: open } });
  const totalPages = Math.max(1, Math.ceil((list.result.total ?? 0) / 8));
  useEffect(() => { setPage(1); }, [debouncedSearch]);
  const upload = async (file: File) => { setUploading(true); try { const media = await imageContentMediaFromUpload(file); onSelect(media); onClose(); } finally { setUploading(false); } };
  const rows = useMemo(() => list.result.data, [list.result.data]);
  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm"><DialogTitle>Chọn hình ảnh</DialogTitle><DialogContent dividers><Stack spacing={2}><Stack direction="row" spacing={1}><TextField size="small" fullWidth placeholder="Tìm tên file..." value={search} onChange={(event) => setSearch(event.target.value)} InputProps={{ startAdornment: <SearchIcon color="action" sx={{ mr: 1 }} /> }} /><Button component="label" variant="outlined" disabled={uploading}>{uploading ? <CircularProgress size={18} /> : "Tải ảnh"}<input hidden type="file" accept="image/*" onChange={(event) => { const file = event.target.files?.[0]; if (file) void upload(file); event.target.value = ""; }} /></Button></Stack>{list.query.isLoading ? <Box textAlign="center" py={5}><CircularProgress /></Box> : list.query.isError ? <Typography color="error" py={3}>Không thể tải thư viện ảnh.</Typography> : rows.length === 0 ? <Typography color="text.secondary" textAlign="center" py={5}>Không có ảnh phù hợp.</Typography> : <List disablePadding>{rows.map((media) => <ListItemButton key={media.id} selected={selected?.id === media.id} onClick={() => { onSelect(media); onClose(); }}><Avatar variant="rounded" src={mediaPreviewUrl(media)} sx={{ width: 56, height: 42, mr: 1.5 }} /><ListItemText primary={media.fileName} secondary={media.mimeType} /></ListItemButton>)}</List>}<Stack direction="row" justifyContent="flex-end"><Pagination size="small" count={totalPages} page={Math.min(page, totalPages)} onChange={(_, value) => setPage(value)} /></Stack></Stack></DialogContent><DialogActions><Button onClick={onClose}>Hủy</Button></DialogActions></Dialog>;
};


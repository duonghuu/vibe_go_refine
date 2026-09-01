import React, { useEffect, useMemo, useState } from "react";
import { HttpError, useList } from "@refinedev/core";
import { Avatar, Box, Button, Checkbox, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, List, ListItem, ListItemAvatar, ListItemButton, ListItemText, Pagination, Radio, Stack, TextField, Typography } from "@mui/material";
import SearchOutlinedIcon from "@mui/icons-material/SearchOutlined";
import { useDebounce } from "../../../hooks/use-debounce";
import { ISectionPickerOption, PageSectionItemType } from "./page-section-types";

interface Props { open: boolean; title: string; itemType: PageSectionItemType; multiple?: boolean; selected: ISectionPickerOption[]; onClose: () => void; onConfirm: (items: ISectionPickerOption[]) => void; }

export const SectionItemPickerDialog: React.FC<Props> = ({ open, title, itemType, multiple = true, selected, onClose, onConfirm }) => {
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [selection, setSelection] = useState<Map<number, ISectionPickerOption>>(new Map());
  const debouncedSearch = useDebounce(search, 300);
  const list = useList<ISectionPickerOption, HttpError>({
    resource: "section-picker-options", dataProviderName: "pageSections",
    pagination: { currentPage: page, pageSize: 5 }, filters: [{ field: "search", operator: "contains", value: debouncedSearch }],
    meta: { itemType }, queryOptions: { enabled: open },
  });
  useEffect(() => { if (open) setSelection(new Map(selected.map((item) => [item.id, item]))); }, [open, selected]);
  useEffect(() => setPage(1), [debouncedSearch, itemType]);
  const totalPages = Math.max(1, Math.ceil((list.result.total ?? 0) / 5));
  const toggle = (option: ISectionPickerOption) => setSelection((current) => {
    if (!multiple) return new Map([[option.id, option]]);
    const next = new Map(current); if (next.has(option.id)) next.delete(option.id); else next.set(option.id, option); return next;
  });
  const orderedSelection = useMemo(() => Array.from(selection.values()), [selection]);
  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
    <DialogTitle>{title}</DialogTitle>
    <DialogContent dividers>
      <TextField fullWidth size="small" value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Tìm theo tên hoặc slug..." InputProps={{ startAdornment: <SearchOutlinedIcon color="action" sx={{ mr: 1 }} /> }} />
      {list.query.isLoading ? <Box display="grid" sx={{ placeItems: "center", py: 6 }}><CircularProgress size={28} /></Box> : list.query.isError ? <Typography color="error" py={3}>Không thể tải danh sách. Vui lòng thử lại.</Typography> : list.result.data.length === 0 ? <Typography color="text.secondary" textAlign="center" py={6}>Không có dữ liệu phù hợp.</Typography> : <List disablePadding sx={{ mt: 1 }}>{list.result.data.map((option) => <ListItem key={option.id} disablePadding secondaryAction={multiple ? <Checkbox checked={selection.has(option.id)} onChange={() => toggle(option)} /> : <Radio checked={selection.has(option.id)} onChange={() => toggle(option)} />}><ListItemButton onClick={() => toggle(option)}><ListItemAvatar><Avatar src={option.imageUrl} variant="rounded">{option.label.slice(0, 1)}</Avatar></ListItemAvatar><ListItemText primary={option.label} secondary={option.secondary} /></ListItemButton></ListItem>)}</List>}
      <Stack direction="row" justifyContent="space-between" alignItems="center" mt={2}><Typography variant="caption" color="text.secondary">Đã chọn {selection.size}</Typography><Pagination size="small" count={totalPages} page={Math.min(page, totalPages)} onChange={(_, value) => setPage(value)} /></Stack>
    </DialogContent>
    <DialogActions><Button type="button" color="secondary" onClick={onClose}>Hủy</Button><Button type="button" variant="contained" disabled={selection.size === 0} onClick={() => { onConfirm(orderedSelection); onClose(); }}>Xác nhận ({selection.size})</Button></DialogActions>
  </Dialog>;
};

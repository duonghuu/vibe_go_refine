import { useEffect, useState } from "react";
import { useList, useNotification } from "@refinedev/core";
import {
  Avatar,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Stack,
  Typography,
} from "@mui/material";
import ArrowDownwardIcon from "@mui/icons-material/ArrowDownward";
import ArrowUpwardIcon from "@mui/icons-material/ArrowUpward";
import { customRequest } from "../../../providers/data";
import { API_URL } from "../../../providers/constants";
import { IImageContent, mediaPreviewUrl } from "../image-content-types";

interface Props {
  open: boolean;
  typeCode: string;
  onClose: () => void;
  onSaved: () => void;
}
export const ImageContentReorderDialog = ({
  open,
  typeCode,
  onClose,
  onSaved,
}: Props) => {
  const [items, setItems] = useState<IImageContent[]>([]);
  const [saving, setSaving] = useState(false);
  const { open: notify } = useNotification();
  const list = useList<IImageContent>({
    resource: "image-contents",
    dataProviderName: "imageContent",
    pagination: { currentPage: 1, pageSize: 100 },
    filters: [{ field: "typeCode", operator: "eq", value: typeCode }],
    sorters: [{ field: "sortOrder", order: "asc" }],
    queryOptions: { enabled: open && Boolean(typeCode) },
  });
  useEffect(() => {
    if (open) setItems(list.result.data);
  }, [open, list.result.data]);
  const move = (index: number, direction: -1 | 1) => {
    const target = index + direction;
    if (target < 0 || target >= items.length) return;
    const next = [...items];
    [next[index], next[target]] = [next[target], next[index]];
    setItems(next);
  };
  const save = async () => {
    setSaving(true);
    try {
      const response = await customRequest({
        url: `${API_URL}/admin/image-contents/order`,
        method: "PUT",
        payload: {
          typeCode,
          items: items.map((item, index) => ({
            id: item.id,
            sortOrder: index,
          })),
        },
      });
      if (!response.ok) throw new Error("Không thể lưu thứ tự.");
      notify?.({ type: "success", message: "Đã lưu thứ tự nội dung." });
      onSaved();
      onClose();
    } catch (error) {
      notify?.({
        type: "error",
        message:
          error instanceof Error ? error.message : "Không thể lưu thứ tự.",
      });
    } finally {
      setSaving(false);
    }
  };
  return (
    <Dialog
      open={open}
      onClose={saving ? undefined : onClose}
      fullWidth
      maxWidth="md"
    >
      <DialogTitle>Sắp xếp nội dung · {typeCode}</DialogTitle>
      <DialogContent dividers>
        {list.query.isLoading ? (
          <Typography py={4} textAlign="center">
            Đang tải thứ tự...
          </Typography>
        ) : list.query.isError ? (
          <Typography color="error" py={4}>
            Không thể tải dữ liệu.
          </Typography>
        ) : (
          <List>
            {items.map((item, index) => (
              <ListItem
                key={item.id}
                secondaryAction={
                  <Stack direction="row">
                    <IconButton
                      aria-label="Đưa lên"
                      disabled={index === 0}
                      onClick={() => move(index, -1)}
                    >
                      <ArrowUpwardIcon />
                    </IconButton>
                    <IconButton
                      aria-label="Đưa xuống"
                      disabled={index === items.length - 1}
                      onClick={() => move(index, 1)}
                    >
                      <ArrowDownwardIcon />
                    </IconButton>
                  </Stack>
                }
              >
                <Avatar
                  variant="rounded"
                  src={mediaPreviewUrl(item.media)}
                  sx={{ mr: 1.5 }}
                />
                <ListItemText
                  primary={item.name || `${item.typeName} #${item.id}`}
                  secondary={`${item.status === "ACTIVE" ? "Đang hoạt động" : "Tạm ẩn"} · vị trí ${index + 1}`}
                />
              </ListItem>
            ))}
          </List>
        )}
        {!list.query.isLoading && !list.query.isError && items.length === 0 && (
          <Box py={4}>
            <Typography color="text.secondary" textAlign="center">
              Type này chưa có nội dung.
            </Typography>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={saving}>
          Hủy
        </Button>
        <Button
          variant="contained"
          onClick={() => void save()}
          disabled={saving || items.length < 2}
        >
          Lưu thứ tự
        </Button>
      </DialogActions>
    </Dialog>
  );
};

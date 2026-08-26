import React, { useRef, useState } from "react";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  IconButton,
  LinearProgress,
  Skeleton,
  Stack,
  Typography,
} from "@mui/material";
import AddPhotoAlternateOutlinedIcon from "@mui/icons-material/AddPhotoAlternateOutlined";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import DragIndicatorIcon from "@mui/icons-material/DragIndicator";
import ReplayIcon from "@mui/icons-material/Replay";
import { IUploadedMedia } from "./post-media-types";

interface MediaUploaderProps {
  multiple?: boolean;
  disabled?: boolean;
  onFiles: (files: FileList) => void;
  label: string;
}

const MediaUploader: React.FC<MediaUploaderProps> = ({ multiple = false, disabled = false, onFiles, label }) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setIsDragOver(false);
    if (!disabled && event.dataTransfer.files.length > 0) onFiles(event.dataTransfer.files);
  };

  return (
    <Box
      onDragOver={(event) => { event.preventDefault(); setIsDragOver(true); }}
      onDragLeave={() => setIsDragOver(false)}
      onDrop={handleDrop}
      sx={{
        border: 1,
        borderStyle: "dashed",
        borderColor: isDragOver ? "primary.main" : "divider",
        bgcolor: isDragOver ? "action.hover" : "background.paper",
        borderRadius: 1,
        p: 2,
        textAlign: "center",
        transition: "border-color 150ms ease, background-color 150ms ease",
      }}
    >
      <input
        ref={inputRef}
        hidden
        type="file"
        accept="image/jpeg,image/png,image/webp"
        multiple={multiple}
        onChange={(event) => {
          if (event.target.files && event.target.files.length > 0) onFiles(event.target.files);
          event.target.value = "";
        }}
      />
      <AddPhotoAlternateOutlinedIcon color="primary" sx={{ fontSize: 34, mb: 0.5 }} />
      <Typography variant="body2" color="text.secondary" mb={1}>
        Kéo thả ảnh vào đây hoặc chọn từ máy tính
      </Typography>
      <Button
        variant="outlined"
        size="small"
        disabled={disabled}
        onClick={() => inputRef.current?.click()}
        sx={{ textTransform: "none" }}
      >
        {label}
      </Button>
      <Typography variant="caption" display="block" color="text.disabled" mt={1}>
        JPG, PNG, WEBP · tối đa 5MB mỗi ảnh
      </Typography>
    </Box>
  );
};

interface MediaPreviewProps {
  media: IUploadedMedia;
  onRemove: () => void;
  compact?: boolean;
}

const MediaPreview: React.FC<MediaPreviewProps> = ({ media, onRemove, compact = false }) => (
  <Box sx={{ position: "relative", border: 1, borderColor: "divider", borderRadius: 1, overflow: "hidden", bgcolor: "background.default" }}>
    {media.originalUrl ? (
      <Box
        component="img"
        src={media.thumbnailUrl || media.mediumUrl || media.originalUrl}
        alt={media.fileName}
        sx={{ width: "100%", height: compact ? 120 : 180, objectFit: "cover", display: "block" }}
      />
    ) : (
      <Box sx={{ height: compact ? 120 : 180, display: "grid", placeItems: "center" }}>
        <Typography variant="caption" color="text.secondary">Ảnh chưa có preview</Typography>
      </Box>
    )}
    <IconButton
      size="small"
      aria-label={`Xóa ${media.fileName}`}
      onClick={onRemove}
      sx={{ position: "absolute", top: 6, right: 6, bgcolor: "background.paper", "&:hover": { bgcolor: "background.paper" } }}
    >
      <DeleteOutlineIcon fontSize="small" color="error" />
    </IconButton>
    <Box sx={{ px: 1, py: 0.75 }}>
      <Typography variant="caption" noWrap display="block" title={media.fileName}>{media.fileName}</Typography>
    </Box>
  </Box>
);

interface PostMediaCardProps {
  thumbnail: IUploadedMedia | null;
  gallery: IUploadedMedia[];
  isUploading: boolean;
  uploadingCount: number;
  isLoading: boolean;
  error: string | null;
  onUploadThumbnail: (file: File) => void;
  onUploadGallery: (files: FileList) => void;
  onRemoveThumbnail: () => void;
  onRemoveGalleryItem: (mediaId: number) => void;
  onReorderGallery: (fromIndex: number, toIndex: number) => void;
  onRetry: () => void;
}

export const PostMediaCard: React.FC<PostMediaCardProps> = ({
  thumbnail, gallery, isUploading, uploadingCount, isLoading, error,
  onUploadThumbnail, onUploadGallery, onRemoveThumbnail, onRemoveGalleryItem,
  onReorderGallery, onRetry,
}) => {
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);

  return (
    <Card sx={{ borderRadius: 1.75, border: 1, borderColor: "divider", boxShadow: "none" }}>
      <CardContent sx={{ p: 3 }}>
        <Typography variant="h6" fontWeight={600} mb={3}>Hình ảnh bài viết</Typography>
        {isLoading ? (
          <Stack spacing={1.5}>
            <Skeleton variant="rectangular" height={92} />
            <Skeleton variant="text" width="55%" />
            <Typography variant="body2" color="text.secondary">Đang tải hình ảnh...</Typography>
          </Stack>
        ) : (
          <Stack spacing={3}>
            {error && <Alert severity="error" action={<Button color="inherit" size="small" startIcon={<ReplayIcon />} onClick={onRetry}>Thử lại</Button>}>{error}</Alert>}
            <Box>
              <Typography variant="subtitle2" mb={1}>Ảnh đại diện</Typography>
              {thumbnail ? (
                <MediaPreview media={thumbnail} onRemove={onRemoveThumbnail} />
              ) : (
                <MediaUploader label="Chọn ảnh đại diện" disabled={isUploading} onFiles={(files) => { const file = files.item(0); if (file) onUploadThumbnail(file); }} />
              )}
            </Box>
            <Box>
              <Typography variant="subtitle2" mb={1}>Thư viện ảnh</Typography>
              <MediaUploader multiple label="Thêm ảnh vào thư viện" disabled={isUploading} onFiles={onUploadGallery} />
              {isUploading && (
                <Stack direction="row" spacing={1} alignItems="center" mt={1}>
                  <LinearProgress sx={{ flex: 1 }} />
                  <Typography variant="caption" color="text.secondary">Đang tải {uploadingCount} ảnh</Typography>
                </Stack>
              )}
              {gallery.length === 0 ? (
                <Typography variant="body2" color="text.secondary" textAlign="center" mt={2}>Chưa có ảnh trong thư viện</Typography>
              ) : (
                <Box sx={{ display: "grid", gridTemplateColumns: "repeat(2, minmax(0, 1fr))", gap: 1.5, mt: 2 }}>
                  {gallery.map((media, index) => (
                    <Box
                      key={media.id}
                      draggable
                      onDragStart={() => setDraggedIndex(index)}
                      onDragOver={(event) => event.preventDefault()}
                      onDrop={() => {
                        if (draggedIndex !== null) onReorderGallery(draggedIndex, index);
                        setDraggedIndex(null);
                      }}
                      onDragEnd={() => setDraggedIndex(null)}
                      sx={{ position: "relative", opacity: draggedIndex === index ? 0.5 : 1, cursor: "grab" }}
                    >
                      <MediaPreview compact media={media} onRemove={() => onRemoveGalleryItem(media.id)} />
                      <DragIndicatorIcon fontSize="small" color="action" sx={{ position: "absolute", bottom: 28, left: 5, bgcolor: "background.paper", borderRadius: 0.5 }} />
                    </Box>
                  ))}
                </Box>
              )}
            </Box>
          </Stack>
        )}
      </CardContent>
    </Card>
  );
};

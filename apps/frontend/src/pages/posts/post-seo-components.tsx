import React from "react";
import { UseFormReturn } from "react-hook-form";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  FormControl,
  FormHelperText,
  InputLabel,
  MenuItem,
  Select,
  Skeleton,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import RefreshIcon from "@mui/icons-material/Refresh";
import {
  IPostSeoFormValues,
  IPostSeoPreview,
  SeoOgType,
  SeoRobots,
  SeoTwitterCard,
} from "./post-seo-types";

interface PostSeoCardProps {
  form: UseFormReturn<IPostSeoFormValues>;
  preview: IPostSeoPreview;
  isLoading: boolean;
  isSaving: boolean;
  error: Error | null;
  onRetry: () => Promise<void>;
}

const isValidUrl = (value: string): boolean => {
  if (!value.trim()) return true;
  try {
    const url = new URL(value);
    return Boolean(url.protocol && url.host);
  } catch {
    return false;
  }
};

const fieldError = (
  message: string | undefined,
  fallback: string,
): string | undefined => message ?? fallback;

const PreviewLine = ({
  label,
  value,
  fallback,
}: {
  label: string;
  value: string;
  fallback: boolean;
}) => (
  <Box>
    <Typography variant="caption" color="text.secondary">
      {label}
    </Typography>
    <Typography
      variant="body2"
      color={value ? "text.primary" : "text.disabled"}
      sx={{ overflowWrap: "anywhere" }}
    >
      {value || "Chưa có dữ liệu"}
      {fallback ? " (fallback)" : ""}
    </Typography>
  </Box>
);

export const PostSeoCard: React.FC<PostSeoCardProps> = ({
  form,
  preview,
  isLoading,
  isSaving,
  error,
  onRetry,
}) => {
  const {
    register,
    formState: { errors },
    watch,
  } = form;
  const description = watch("metaDescription");
  const ogDescription = watch("ogDescription");
  const twitterDescription = watch("twitterDescription");

  if (isLoading) {
    return (
      <Card sx={{ borderRadius: 1.75, border: 1, borderColor: "divider", boxShadow: "none" }}>
        <CardContent sx={{ p: 3 }}>
          <Skeleton variant="text" width={150} height={36} sx={{ mb: 2 }} />
          <Skeleton variant="rounded" height={56} sx={{ mb: 2 }} />
          <Skeleton variant="rounded" height={110} sx={{ mb: 2 }} />
          <Skeleton variant="rounded" height={56} />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card sx={{ borderRadius: 1.75, border: 1, borderColor: "divider", boxShadow: "none" }}>
      <CardContent sx={{ p: 3 }}>
        <Stack direction="row" justifyContent="space-between" alignItems="center" mb={1}>
          <Typography variant="h6" fontWeight={600} color="text.primary">
            SEO Meta
          </Typography>
          {isSaving && (
            <Typography variant="caption" color="text.secondary">
              Đang lưu...
            </Typography>
          )}
        </Stack>

        <Typography variant="body2" color="text.secondary" mb={2}>
          Tùy chỉnh thông tin hiển thị trên công cụ tìm kiếm và mạng xã hội.
        </Typography>

        {error && (
          <Alert
            severity="error"
            sx={{ mb: 2 }}
            action={
              <Button
                color="inherit"
                size="small"
                startIcon={<RefreshIcon />}
                onClick={() => { void onRetry(); }}
              >
                Thử lại
              </Button>
            }
          >
            {error.message}
          </Alert>
        )}

        <Accordion defaultExpanded disableGutters elevation={0}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography fontWeight={600}>Cơ bản</Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ px: 0 }}>
            <TextField
              {...register("metaTitle", {
                maxLength: { value: 255, message: "SEO Title tối đa 255 ký tự" },
              })}
              label="SEO Title"
              fullWidth
              error={Boolean(errors.metaTitle)}
              helperText={fieldError(errors.metaTitle?.message, "Để trống để dùng tiêu đề bài viết")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("metaDescription", {
                maxLength: { value: 500, message: "SEO Description tối đa 500 ký tự" },
              })}
              label="SEO Description"
              fullWidth
              multiline
              minRows={3}
              error={Boolean(errors.metaDescription)}
              helperText={
                errors.metaDescription?.message ??
                `${description.length}/500 ký tự · Khuyến nghị khoảng 150–160 ký tự`
              }
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("metaKeywords", {
                maxLength: { value: 500, message: "Meta Keywords tối đa 500 ký tự" },
              })}
              label="Meta Keywords"
              fullWidth
              error={Boolean(errors.metaKeywords)}
              helperText={fieldError(errors.metaKeywords?.message, "Phân tách các từ khóa bằng dấu phẩy")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <FormControl fullWidth error={Boolean(errors.robots)}>
              <InputLabel id="post-seo-robots-label">Robots</InputLabel>
              <Select<SeoRobots>
                {...register("robots")}
                labelId="post-seo-robots-label"
                label="Robots"
                defaultValue="index,follow"
                sx={{ borderRadius: 1 }}
              >
                <MenuItem value="index,follow">index,follow</MenuItem>
                <MenuItem value="noindex,follow">noindex,follow</MenuItem>
                <MenuItem value="noindex,nofollow">noindex,nofollow</MenuItem>
              </Select>
              <FormHelperText>{errors.robots?.message ?? "Mặc định: index,follow"}</FormHelperText>
            </FormControl>
          </AccordionDetails>
        </Accordion>

        <Accordion disableGutters elevation={0}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography fontWeight={600}>Canonical</Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ px: 0 }}>
            <TextField
              {...register("canonicalUrl", {
                validate: (value) => isValidUrl(value) || "Canonical URL không hợp lệ",
                maxLength: { value: 500, message: "Canonical URL tối đa 500 ký tự" },
              })}
              label="Canonical URL"
              fullWidth
              error={Boolean(errors.canonicalUrl)}
              helperText={fieldError(errors.canonicalUrl?.message, "Để trống để tự động dùng URL từ slug")}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
          </AccordionDetails>
        </Accordion>

        <Accordion disableGutters elevation={0}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography fontWeight={600}>Open Graph</Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ px: 0 }}>
            <TextField
              {...register("ogTitle", {
                maxLength: { value: 255, message: "OG Title tối đa 255 ký tự" },
              })}
              label="OG Title"
              fullWidth
              error={Boolean(errors.ogTitle)}
              helperText={fieldError(errors.ogTitle?.message, "Để trống để dùng SEO Title")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("ogDescription", {
                maxLength: { value: 500, message: "OG Description tối đa 500 ký tự" },
              })}
              label="OG Description"
              fullWidth
              multiline
              minRows={3}
              error={Boolean(errors.ogDescription)}
              helperText={errors.ogDescription?.message ?? `${ogDescription.length}/500 ký tự`}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("ogImage", {
                validate: (value) => isValidUrl(value) || "OG Image phải là URL hợp lệ",
                maxLength: { value: 500, message: "OG Image tối đa 500 ký tự" },
              })}
              label="OG Image URL"
              fullWidth
              error={Boolean(errors.ogImage)}
              helperText={fieldError(errors.ogImage?.message, "Để trống để dùng thumbnail bài viết")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <FormControl fullWidth error={Boolean(errors.ogType)}>
              <InputLabel id="post-seo-og-type-label">OG Type</InputLabel>
              <Select<SeoOgType | "">
                {...register("ogType")}
                labelId="post-seo-og-type-label"
                label="OG Type"
                defaultValue=""
                sx={{ borderRadius: 1 }}
              >
                <MenuItem value="">Mặc định</MenuItem>
                <MenuItem value="article">article</MenuItem>
                <MenuItem value="website">website</MenuItem>
                <MenuItem value="product">product</MenuItem>
              </Select>
            </FormControl>
          </AccordionDetails>
        </Accordion>

        <Accordion disableGutters elevation={0}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography fontWeight={600}>Twitter Card</Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ px: 0 }}>
            <TextField
              {...register("twitterTitle", {
                maxLength: { value: 255, message: "Twitter Title tối đa 255 ký tự" },
              })}
              label="Twitter Title"
              fullWidth
              error={Boolean(errors.twitterTitle)}
              helperText={fieldError(errors.twitterTitle?.message, "Để trống để dùng OG Title")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("twitterDescription", {
                maxLength: { value: 500, message: "Twitter Description tối đa 500 ký tự" },
              })}
              label="Twitter Description"
              fullWidth
              multiline
              minRows={3}
              error={Boolean(errors.twitterDescription)}
              helperText={errors.twitterDescription?.message ?? `${twitterDescription.length}/500 ký tự`}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <TextField
              {...register("twitterImage", {
                validate: (value) => isValidUrl(value) || "Twitter Image phải là URL hợp lệ",
                maxLength: { value: 500, message: "Twitter Image tối đa 500 ký tự" },
              })}
              label="Twitter Image URL"
              fullWidth
              error={Boolean(errors.twitterImage)}
              helperText={fieldError(errors.twitterImage?.message, "Để trống để dùng OG Image")}
              sx={{ mb: 2 }}
              InputProps={{ sx: { borderRadius: 1 } }}
            />
            <FormControl fullWidth error={Boolean(errors.twitterCard)}>
              <InputLabel id="post-seo-twitter-card-label">Twitter Card</InputLabel>
              <Select<SeoTwitterCard | "">
                {...register("twitterCard")}
                labelId="post-seo-twitter-card-label"
                label="Twitter Card"
                defaultValue=""
                sx={{ borderRadius: 1 }}
              >
                <MenuItem value="">Mặc định</MenuItem>
                <MenuItem value="summary">summary</MenuItem>
                <MenuItem value="summary_large_image">summary_large_image</MenuItem>
              </Select>
            </FormControl>
          </AccordionDetails>
        </Accordion>

        <Accordion disableGutters elevation={0}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography fontWeight={600}>Structured Data</Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ px: 0 }}>
            <TextField
              {...register("schemaJsonText", {
                validate: (value) => {
                  if (!value.trim()) return true;
                  try {
                    JSON.parse(value);
                    return true;
                  } catch {
                    return "Schema JSON-LD không hợp lệ";
                  }
                },
              })}
              label="Schema JSON-LD"
              fullWidth
              multiline
              minRows={8}
              error={Boolean(errors.schemaJsonText)}
              helperText={errors.schemaJsonText?.message ?? "Nhập JSON-LD hợp lệ nếu cần khai báo structured data"}
              InputProps={{ sx: { borderRadius: 1, fontFamily: "monospace" } }}
            />
          </AccordionDetails>
        </Accordion>

        <Box sx={{ mt: 2, p: 2, bgcolor: "background.default", border: 1, borderColor: "divider", borderRadius: 1.5 }}>
          <Typography variant="subtitle2" fontWeight={600} mb={1}>
            SEO Preview
          </Typography>
          <Stack spacing={1.5}>
            <PreviewLine label="Title" value={preview.title} fallback={preview.isTitleFallback} />
            <PreviewLine label="Description" value={preview.description} fallback={preview.isDescriptionFallback} />
            <PreviewLine label="Canonical" value={preview.canonicalUrl} fallback={preview.isCanonicalFallback} />
            {preview.ogImage && (
              <Box>
                <Typography variant="caption" color="text.secondary">OG Image</Typography>
                <Box
                  component="img"
                  src={preview.ogImage}
                  alt="OG preview"
                  sx={{ display: "block", mt: 0.5, width: "100%", maxHeight: 160, objectFit: "cover", borderRadius: 1 }}
                />
                {preview.isOgImageFallback && (
                  <Typography variant="caption" color="text.secondary">Fallback từ thumbnail bài viết</Typography>
                )}
              </Box>
            )}
          </Stack>
        </Box>
      </CardContent>
    </Card>
  );
};


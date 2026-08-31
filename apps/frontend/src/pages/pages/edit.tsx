import React, { useEffect, useState } from "react";
import { Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { HttpError, useGo, useNotification } from "@refinedev/core";
import { useParams } from "react-router";
import { Alert, Box, Button, Card, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, FormControl, FormHelperText, Grid, InputLabel, MenuItem, Select, Skeleton, Stack, Typography } from "@mui/material";
import { IPageResponse, IPageUpdatePayload, PageStatus } from "./page-form-types";
import { PageBasicInfoCard } from "./page-basic-info-card";
import { PagePostCreateNotice } from "./page-post-create-notice";
import { PagePublicUrlPreview } from "./page-public-url-preview";
import { PageEditActions } from "./page-edit-actions";

const getErrorMessage = (error: unknown, fallback: string): string =>
  error instanceof Error && error.message ? error.message : fallback;

export const PageEdit: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const go = useGo();
  const { open } = useNotification();
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  const pageId = id && Number.isInteger(Number(id)) && Number(id) > 0 ? Number(id) : undefined;

  const form = useForm<IPageResponse, HttpError, IPageUpdatePayload, object, IPageResponse, IPageResponse>({
    mode: "onBlur",
    refineCoreProps: {
      action: "edit",
      resource: "pages",
      redirect: false,
      id: pageId,
      meta: { method: "put" },
      onMutationSuccess: () => {
        open?.({ type: "success", message: "Cập nhật trang thành công", description: "Các thay đổi đã được lưu." });
        go({ to: "/pages" });
      },
      onMutationError: (error) => {
        if (error.statusCode === 409) {
          form.setError("slug", { type: "server", message: error.message || "Slug đã tồn tại" });
        }
        open?.({ type: "error", message: "Không thể cập nhật trang", description: getErrorMessage(error, "Vui lòng kiểm tra dữ liệu và thử lại.") });
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const { register, handleSubmit, watch, setValue, formState: { errors, isDirty }, refineCore: { formLoading, query: pageQuery, onFinish } } = form;
  const page = pageQuery?.data?.data;
  const slug = watch("slug") ?? "";
  const status = watch("status") ?? "DRAFT";

  useEffect(() => {
    if (!pageId) {
      open?.({ type: "error", message: "Không tìm thấy trang", description: "Đường dẫn trang không hợp lệ." });
      go({ to: "/pages" });
    }
  }, [go, open, pageId]);

  const handleCancel = () => {
    if (isDirty) setCancelDialogOpen(true);
    else go({ to: "/pages" });
  };

  if (!pageId) return <Box p={3}>Đang kiểm tra thông tin trang...</Box>;

  if (pageQuery?.isLoading) {
    return <Box><Skeleton variant="text" width={280} height={56} sx={{ mb: 2 }} /><Grid container spacing={3}><Grid item xs={12} md={8}><Card sx={{ border: 1, borderColor: "divider", boxShadow: "none" }}><CardContent sx={{ p: 3 }}><Skeleton variant="text" width={180} height={36} /><Skeleton variant="rounded" height={56} sx={{ my: 2 }} /><Skeleton variant="rounded" height={56} sx={{ mb: 2 }} /><Skeleton variant="rounded" height={360} /></CardContent></Card></Grid><Grid item xs={12} md={4}><Skeleton variant="rounded" height={180} /><Skeleton variant="rounded" height={180} sx={{ mt: 3 }} /></Grid></Grid></Box>;
  }

  if (pageQuery?.isError || !page) {
    return <Stack spacing={2}><Alert severity="error">{getErrorMessage(pageQuery?.error, "Không tìm thấy trang hoặc dữ liệu không khả dụng.")}</Alert><Button variant="outlined" onClick={() => go({ to: "/pages" })} sx={{ alignSelf: "flex-start" }}>Quay về danh sách</Button></Stack>;
  }

  return (
    <Box>
      <Typography variant="h4" fontWeight={700} letterSpacing="-0.02em" mb={0.5}>Chỉnh sửa trang</Typography>
      <Typography variant="body2" color="text.secondary" mb={3}>Trang chủ / Quản trị nội dung / Trang / Chỉnh sửa</Typography>
      <Edit title="" isLoading={formLoading} wrapperProps={{ sx: { bgcolor: "transparent", boxShadow: "none", p: 0 } }} headerProps={{ sx: { display: "none" } }} footerButtons={<PageEditActions disabled={formLoading} onCancel={handleCancel} onSubmit={() => { void handleSubmit((values) => onFinish(values))(); }} />}>
        <Box component="form" autoComplete="off">
          <Grid container spacing={3}>
            <Grid item xs={12} md={8}><PageBasicInfoCard register={register} errors={errors} watch={watch} onSlugChange={(value) => setValue("slug", value, { shouldDirty: true, shouldValidate: true })} /><Box mt={3}><PagePublicUrlPreview slug={slug} /></Box></Grid>
            <Grid item xs={12} md={4}>
              <Card sx={{ border: 1, borderColor: "divider", borderRadius: 2, boxShadow: "none", mb: 3 }}><CardContent sx={{ p: { xs: 2, md: 3 } }}><Typography variant="h6" fontWeight={700} mb={1.5}>Trạng thái</Typography><FormControl fullWidth error={Boolean(errors.status)}><InputLabel id="page-status-label">Trạng thái</InputLabel><Select<IPageUpdatePayload["status"]> {...form.register("status", { required: "Trạng thái là bắt buộc" })} labelId="page-status-label" label="Trạng thái" value={status} onChange={(event) => setValue("status", event.target.value as PageStatus, { shouldDirty: true, shouldValidate: true })}><MenuItem value="DRAFT">Bản nháp</MenuItem><MenuItem value="PUBLISHED">Đã xuất bản</MenuItem></Select><FormHelperText>{errors.status?.message ?? (status === "PUBLISHED" ? "Trang đang ở trạng thái công khai." : "Trang đang ở trạng thái bản nháp.")}</FormHelperText></FormControl></CardContent></Card>
              <Stack spacing={3}><PagePostCreateNotice icon="media" title="Hình ảnh" description="Media Page sẽ được quản lý theo ID trang hiện tại." /><PagePostCreateNotice icon="seo" title="SEO" description="SEO Meta sẽ được cấu hình theo ID trang hiện tại." /></Stack>
            </Grid>
          </Grid>
        </Box>
      </Edit>
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Hủy thay đổi?</DialogTitle>
        <DialogContent><DialogContentText>Bạn có thay đổi chưa lưu. Nếu rời trang, các thay đổi này sẽ bị mất.</DialogContentText></DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)} color="secondary">Tiếp tục chỉnh sửa</Button>
          <Button onClick={() => go({ to: "/pages" })} variant="contained">Rời trang</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

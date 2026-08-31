import React, { useEffect, useState } from "react";
import { Create } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { HttpError, useGo, useNotification } from "@refinedev/core";
import { Box, Button, Card, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Grid, Stack, Typography } from "@mui/material";
import { IPageCreatePayload, IPageResponse, PageStatus } from "./page-form-types";
import { PageBasicInfoCard } from "./page-basic-info-card";
import { PagePostCreateNotice } from "./page-post-create-notice";
import { generateSlug } from "../../utils/generate-slug";

export const PageCreate: React.FC = () => {
  const go = useGo();
  const { open } = useNotification();
  const [submitStatus, setSubmitStatus] = useState<PageStatus | null>(null);
  const [isSlugManuallyEdited, setIsSlugManuallyEdited] = useState(false);
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);

  const form = useForm<IPageResponse, HttpError, IPageCreatePayload, object, IPageResponse, IPageResponse>({
    mode: "onBlur",
    refineCoreProps: {
      action: "create",
      resource: "pages",
      redirect: false,
      onMutationSuccess: ({ data }) => {
        open?.({
          type: "success",
          message: submitStatus === "PUBLISHED" ? "Xuất bản trang thành công" : "Lưu nháp thành công",
          description: "Đang chuyển tới màn hình chỉnh sửa trang.",
        });
        go({ to: `/pages/edit/${data.id}` });
      },
      onMutationError: (error) => {
        setSubmitStatus(null);
        if (error.statusCode === 409) {
          form.setError("slug", { type: "server", message: error.message || "Slug đã tồn tại" });
        }
        open?.({ type: "error", message: "Không thể tạo trang", description: error.message });
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const { register, handleSubmit, setValue, watch, formState: { errors, isDirty }, refineCore: { formLoading, onFinish } } = form;
  const title = watch("title");

  useEffect(() => {
    if (!isSlugManuallyEdited) setValue("slug", generateSlug(title ?? ""), { shouldValidate: true });
  }, [title, isSlugManuallyEdited, setValue]);

  const submitPage = (status: PageStatus) => handleSubmit((values) => {
    setSubmitStatus(status);
    void onFinish({ ...values, status });
  })();

  const handleSlugChange = (value: string) => {
    setIsSlugManuallyEdited(value.length > 0);
    setValue("slug", value, { shouldDirty: true, shouldValidate: true });
  };

  const handleCancel = () => {
    if (isDirty) setCancelDialogOpen(true);
    else go({ to: "/pages" });
  };

  const isSaving = formLoading;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" fontWeight={700} letterSpacing="-0.02em">Tạo trang mới</Typography>
          <Typography variant="body2" color="text.secondary" mt={0.5}>Trang chủ / Quản trị nội dung / Trang / Tạo mới</Typography>
        </Box>
      </Box>
      <Create
        title=""
        isLoading={formLoading}
        wrapperProps={{ sx: { bgcolor: "transparent", boxShadow: "none", p: 0 } }}
        headerProps={{ sx: { display: "none" } }}
        footerButtons={(
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: "100%" }}>
            <Button variant="outlined" color="secondary" disabled={isSaving} onClick={handleCancel} sx={{ textTransform: "none", fontWeight: 700 }}>Hủy</Button>
            <Button variant="outlined" color="primary" disabled={isSaving} onClick={() => submitPage("DRAFT")} sx={{ textTransform: "none", fontWeight: 700 }}>Lưu nháp</Button>
            <Button variant="contained" color="success" disabled={isSaving} onClick={() => submitPage("PUBLISHED")} sx={{ textTransform: "none", fontWeight: 700, boxShadow: "none" }}>Xuất bản</Button>
          </Stack>
        )}
      >
        <Box component="form" autoComplete="off">
          <Grid container spacing={3}>
            <Grid item xs={12} md={8}>
              <PageBasicInfoCard register={register} errors={errors} watch={watch} onSlugChange={handleSlugChange} />
            </Grid>
            <Grid item xs={12} md={4}>
              <Card sx={{ border: "1px solid", borderColor: "divider", borderRadius: 2, boxShadow: "none", mb: 3 }}>
                <CardContent sx={{ p: { xs: 2, md: 3 } }}>
                  <Typography variant="h6" fontWeight={700} mb={1.5}>Trạng thái</Typography>
                  <Typography variant="body2" color="text.secondary">{submitStatus === "PUBLISHED" ? "Trang sẽ được xuất bản ngay sau khi lưu." : "Trang được lưu dưới dạng bản nháp để tiếp tục hoàn thiện."}</Typography>
                </CardContent>
              </Card>
              <Stack spacing={3}>
                <PagePostCreateNotice icon="media" title="Hình ảnh" description="Lưu trang trước để quản lý thumbnail và gallery." />
                <PagePostCreateNotice icon="seo" title="SEO" description="Lưu trang trước để cấu hình thông tin SEO." />
              </Stack>
            </Grid>
          </Grid>
        </Box>
      </Create>
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Hủy tạo trang?</DialogTitle>
        <DialogContent><DialogContentText>Dữ liệu chưa lưu sẽ bị mất nếu rời trang.</DialogContentText></DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)} color="secondary">Tiếp tục chỉnh sửa</Button>
          <Button onClick={() => go({ to: "/pages" })} variant="contained">Rời trang</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

import React, { useEffect, useRef, useState } from "react";
import { Create } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { useSearchParams, useNavigate } from "react-router";
import {
  Box,
  Card,
  CardContent,
  TextField,
  Typography,
  Button,
  Stack,
  Grid,
  Autocomplete,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";
import { HttpError, useNotification, useSelect } from "@refinedev/core";
import { Controller } from "react-hook-form";
import { PostMediaCard } from "./post-media-components";
import { usePostMedia } from "./use-post-media";
import { PostSeoCard } from "./post-seo-components";
import { UsePostSeoResult, usePostSeo } from "./use-post-seo";

export interface IPostResponse {
  id: number;
  typeCode: string;
  title: string;
  slug: string;
  content: string;
  authorId: number;
  categoryId?: number;
  createdAt: string;
  thumbnailUrl?: string | null;
}

export interface ICreatePostRequest {
  typeCode: string;
  title: string;
  slug: string;
  content: string;
  categoryId?: number;
}

interface ICreatedPostPayload {
  id?: number;
  data?: { id?: number };
}

interface IPostCategoryOption {
  id: number;
  name: string;
  parentId: number | null;
  typeCode: string;
}

interface IHierarchicalCategoryOption extends IPostCategoryOption {
  level: number;
}

const generateSlug = (text: string) => {
  return text
    .toString()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .trim()
    .replace(/\s+/g, "-")
    .replace(/[^\w-]+/g, "")
    .replace(/--+/g, "-");
};

export const PostCreate: React.FC = () => {
  const [searchParams] = useSearchParams();
  const typeCode = searchParams.get("type_code");
  const navigate = useNavigate();
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  
  const { open } = useNotification();
  const postMedia = usePostMedia();
  const postSeoRef = useRef<UsePostSeoResult | null>(null);

  useEffect(() => {
    if (!typeCode) {
      open?.({
        type: "error",
        message: "Lỗi truy cập",
        description: "Thiếu thông tin mã loại bài viết (type_code). Đang chuyển hướng...",
      });
      setTimeout(() => {
        navigate("/posts");
      }, 2000);
    }
  }, [typeCode, navigate, open]);

  const {
    refineCore: { formLoading, onFinish },
    register,
    control,
    handleSubmit,
    setValue,
    formState: { errors, isDirty },
    watch,
  } = useForm<IPostResponse, HttpError, ICreatePostRequest>({
    refineCoreProps: {
      action: "create",
      resource: "posts",
      redirect: false,
      onMutationSuccess: async (response: unknown) => {
        const createdResponse = response as ICreatedPostPayload | undefined;
        const createdId = createdResponse?.data?.id ?? createdResponse?.id;
        if (!createdId) {
          open?.({ type: "error", message: "Không lấy được mã bài viết", description: "Bài viết đã tạo nhưng chưa thể đồng bộ hình ảnh." });
          return;
        }
        try {
          await postMedia.sync(createdId);
          if (postSeoRef.current?.isDirty) {
            await postSeoRef.current.save(createdId);
          }
          open?.({ type: "success", message: "Tạo bài viết thành công", description: "Bài viết và hình ảnh đã được lưu." });
          navigate(`/posts?type_code=${typeCode}`);
        } catch (error) {
          open?.({ type: "error", message: "Không thể đồng bộ hình ảnh", description: error instanceof Error ? error.message : "Vui lòng mở trang chỉnh sửa để thử lại." });
        }
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const titleValue = watch("title");
  const postSeo = usePostSeo(undefined, {
    title: titleValue ?? "",
    slug: watch("slug") ?? "",
    content: watch("content") ?? "",
    thumbnailUrl: postMedia.thumbnail?.thumbnailUrl ?? postMedia.thumbnail?.originalUrl,
  });
  postSeoRef.current = postSeo;

  const { query: categoryQuery } = useSelect<IPostCategoryOption, HttpError>({
    resource: "post-categories",
    optionLabel: "name",
    optionValue: "id",
    filters: [
      {
        field: "typeCode",
        operator: "eq",
        value: typeCode || "",
      },
    ],
    queryOptions: {
      enabled: !!typeCode,
    },
    pagination: {
      mode: "off",
    },
  });

  const categoryOptions = React.useMemo<IHierarchicalCategoryOption[]>(() => {
    const categories = categoryQuery.data?.data ?? [];
    const childrenByParent = new Map<number | null, IPostCategoryOption[]>();

    categories.forEach((category) => {
      const children = childrenByParent.get(category.parentId) ?? [];
      children.push(category);
      childrenByParent.set(category.parentId, children);
    });

    const flatten = (
      parentId: number | null,
      level: number,
      visited: Set<number>,
    ): IHierarchicalCategoryOption[] => {
      return (childrenByParent.get(parentId) ?? []).flatMap((category) => {
        if (visited.has(category.id)) {
          return [];
        }

        const nextVisited = new Set(visited).add(category.id);
        return [
          { ...category, level },
          ...flatten(category.id, level + 1, nextVisited),
        ];
      });
    };

    const roots = flatten(null, 0, new Set<number>());
    const rootIds = new Set(roots.map((category) => category.id));
    const orphanedCategories = categories.filter(
      (category) => !rootIds.has(category.id),
    );

    return [...roots, ...orphanedCategories.map((category) => ({ ...category, level: 0 }))];
  }, [categoryQuery.data?.data]);

  useEffect(() => {
    if (titleValue) {
      setValue("slug", generateSlug(titleValue), { shouldValidate: true });
    } else {
      setValue("slug", "");
    }
  }, [titleValue, setValue]);

  const onCustomSubmit = async (data: ICreatePostRequest) => {
    if (!typeCode) return;
    if (postMedia.isUploading || postMedia.uploadingCount > 0) return;
    const seoValid = await postSeo.form.trigger();
    if (!seoValid) return;
    
    onFinish({
      ...data,
      typeCode: typeCode,
    });
  };

  const handleCancel = () => {
    if (isDirty || postMedia.isDirty || postSeo.isDirty) {
      setCancelDialogOpen(true);
      return;
    }
    navigate(`/posts?type_code=${typeCode}`);
  };

  if (!typeCode) {
    return <Box p={3}>Đang kiểm tra thông tin mã loại bài viết...</Box>;
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Thêm bài viết mới
        </Typography>
      </Box>

      <Create
        isLoading={formLoading}
        title=""
        wrapperProps={{
          sx: {
            bgcolor: "transparent",
            boxShadow: "none",
            p: 0,
          },
        }}
        headerProps={{
          sx: { display: "none" },
        }}
        footerButtons={
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: '100%' }}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={handleCancel}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", color: 'text.secondary', borderColor: '#D5D5D5' }}
            >
              Hủy
            </Button>
            <Button
              variant="contained"
              color="primary"
              disabled={formLoading || postMedia.isUploading || postMedia.isSyncing || postSeo.isSaving}
              onClick={handleSubmit(onCustomSubmit)}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'primary.main', boxShadow: 'none', '&:hover': { bgcolor: 'primary.dark', boxShadow: 'none' } }}
            >
              Lưu bài viết
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
          <Grid container spacing={3}>
            {/* Cột chính */}
            <Grid item xs={12} md={8}>
              <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                    Thông tin cơ bản
                  </Typography>
                  <TextField
                    {...register("title", { required: "Tiêu đề là bắt buộc" })}
                    error={!!errors.title}
                    helperText={errors.title?.message as string}
                    label="Tiêu đề (*)"
                    margin="normal"
                    variant="outlined"
                    fullWidth
                    sx={{ mb: 3 }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />
                  <TextField
                    {...register("slug", {
                      required: "Đường dẫn (slug) là bắt buộc",
                      pattern: {
                        value: /^[a-z0-9-]+$/,
                        message: "Slug chỉ chứa chữ thường, số và dấu gạch ngang",
                      },
                    })}
                    error={!!errors.slug}
                    helperText={errors.slug?.message as string}
                    label="Đường dẫn (Slug) (*)"
                    margin="normal"
                    variant="outlined"
                    fullWidth
                    sx={{ mb: 3 }}
                    InputLabelProps={{ shrink: !!watch("slug") || undefined }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />
                  <TextField
                    {...register("content", { required: "Nội dung là bắt buộc" })}
                    error={!!errors.content}
                    helperText={errors.content?.message as string}
                    label="Nội dung (*)"
                    margin="normal"
                    variant="outlined"
                    fullWidth
                    multiline
                    rows={15}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />
                </CardContent>
              </Card>
            </Grid>

            {/* Cột phụ */}
            <Grid item xs={12} md={4}>
              <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                    Phân loại
                  </Typography>
                  <Controller
                    control={control}
                    name="categoryId"
                    render={({ field }) => (
                      <Autocomplete
                        options={categoryOptions}
                        loading={categoryQuery.isLoading}
                        value={categoryOptions.find((option) => option.id === field.value) ?? null}
                        onChange={(_, value) => field.onChange(value?.id ?? null)}
                        getOptionLabel={(option) => option.name}
                        isOptionEqualToValue={(option, value) =>
                          option.id === value.id
                        }
                        renderOption={(props, option) => (
                          <li {...props} key={option.id}>
                            <Box
                              component="span"
                              sx={{
                                pl: option.level * 2,
                                color: option.level > 0 ? "text.secondary" : "text.primary",
                              }}
                            >
                              {option.level > 0 ? "└─ " : ""}
                              {option.name}
                            </Box>
                          </li>
                        )}
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            label="Danh mục bài viết"
                            margin="normal"
                            variant="outlined"
                            error={!!errors.categoryId}
                            helperText={errors.categoryId?.message as string}
                            InputProps={{
                              ...params.InputProps,
                              sx: { borderRadius: "8px" },
                            }}
                          />
                        )}
                      />
                    )}
                  />
                </CardContent>
              </Card>
              <Box mt={3}>
                <PostMediaCard
                  thumbnail={postMedia.thumbnail}
                  gallery={postMedia.gallery}
                  isUploading={postMedia.isUploading}
                  uploadingCount={postMedia.uploadingCount}
                  isLoading={postMedia.isLoading}
                  error={postMedia.error}
                  onUploadThumbnail={(file) => { void postMedia.uploadThumbnail(file); }}
                  onUploadGallery={(files) => { void postMedia.uploadGallery(files); }}
                  onRemoveThumbnail={() => { void postMedia.removeThumbnail(); }}
                  onRemoveGalleryItem={(mediaId) => { void postMedia.removeGalleryItem(mediaId); }}
                  onReorderGallery={postMedia.reorderGallery}
                  onRetry={() => { void postMedia.load(); }}
                />
              </Box>
              <Box mt={3}>
                <PostSeoCard
                  form={postSeo.form}
                  preview={postSeo.preview}
                  isLoading={postSeo.isLoading}
                  isSaving={postSeo.isSaving}
                  error={postSeo.error}
                  onRetry={postSeo.retry}
                />
              </Box>
            </Grid>
          </Grid>
        </Box>
      </Create>
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Hủy tạo bài viết?</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Bạn có dữ liệu, hình ảnh hoặc SEO chưa lưu. Nếu rời trang, các thay đổi này sẽ bị mất.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)} color="secondary">Tiếp tục chỉnh sửa</Button>
          <Button onClick={() => navigate(`/posts?type_code=${typeCode}`)} color="primary" variant="contained">Rời trang</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

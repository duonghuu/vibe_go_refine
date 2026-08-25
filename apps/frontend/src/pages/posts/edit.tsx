import React, { useEffect, useMemo, useState } from "react";
import { Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { HttpError, useNotification, useSelect } from "@refinedev/core";
import { Controller } from "react-hook-form";
import { useNavigate, useParams, useSearchParams } from "react-router";
import {
  Alert,
  Autocomplete,
  Box,
  Button,
  Card,
  CardContent,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Grid2,
  Skeleton,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { IPostResponse } from "./create";

interface IUpdatePostRequest {
  title: string;
  slug: string;
  content: string;
  categoryId?: number | null;
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

const flattenCategories = (
  categories: IPostCategoryOption[],
): IHierarchicalCategoryOption[] => {
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

      return [
        { ...category, level },
        ...flatten(category.id, level + 1, new Set(visited).add(category.id)),
      ];
    });
  };

  const roots = flatten(null, 0, new Set<number>());
  const rootIds = new Set(roots.map((category) => category.id));
  const orphaned = categories.filter((category) => !rootIds.has(category.id));

  return [
    ...roots,
    ...orphaned.map((category) => ({ ...category, level: 0 })),
  ];
};

const getErrorMessage = (error: unknown, fallback: string): string => {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
};

export const PostEdit: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { open } = useNotification();
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);

  const {
    refineCore: { formLoading, onFinish, query: postQuery },
    register,
    control,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm<IPostResponse, HttpError, IUpdatePostRequest>({
    refineCoreProps: {
      action: "edit",
      resource: "posts",
      redirect: false,
      meta: {
        method: "put",
      },
      onMutationSuccess: () => {
        const resolvedTypeCode = postQuery?.data?.data?.typeCode ?? searchParams.get("type_code");
        open?.({
          type: "success",
          message: "Cập nhật bài viết thành công",
          description: "Các thay đổi đã được lưu.",
        });
        navigate(resolvedTypeCode ? `/posts?type_code=${resolvedTypeCode}` : "/posts");
      },
      onMutationError: (error) => {
        open?.({
          type: "error",
          message: "Không thể cập nhật bài viết",
          description: getErrorMessage(error, "Vui lòng kiểm tra dữ liệu và thử lại."),
        });
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const postData = postQuery?.data?.data;
  const typeCode = postData?.typeCode ?? searchParams.get("type_code") ?? "";
  const isValidId = Boolean(id && Number.isInteger(Number(id)) && Number(id) > 0);

  const { query: categoryQuery } = useSelect<IPostCategoryOption, HttpError>({
    resource: "post-categories",
    optionLabel: "name",
    optionValue: "id",
    filters: [
      {
        field: "typeCode",
        operator: "eq",
        value: typeCode,
      },
    ],
    pagination: { mode: "off" },
    queryOptions: { enabled: Boolean(typeCode) },
  });

  const categoryOptions = useMemo(
    () => flattenCategories(categoryQuery.data?.data ?? []),
    [categoryQuery.data?.data],
  );

  useEffect(() => {
    if (postQuery?.isError) {
      open?.({
        type: "error",
        message: "Không thể tải bài viết",
        description: getErrorMessage(postQuery.error, "Đã xảy ra lỗi khi tải dữ liệu."),
      });
    }
  }, [open, postQuery?.error, postQuery?.isError]);

  useEffect(() => {
    if (categoryQuery.isError) {
      open?.({
        type: "error",
        message: "Không thể tải danh mục bài viết",
        description: getErrorMessage(categoryQuery.error, "Danh mục tạm thời không khả dụng."),
      });
    }
  }, [categoryQuery.error, categoryQuery.isError, open]);

  useEffect(() => {
    if (!isValidId) {
      open?.({
        type: "error",
        message: "Không tìm thấy bài viết",
        description: "Đường dẫn bài viết không hợp lệ.",
      });
      navigate("/posts");
    }
  }, [isValidId, navigate, open]);

  const navigateBack = () => {
    navigate(typeCode ? `/posts?type_code=${typeCode}` : "/posts");
  };

  const handleCancel = () => {
    if (isDirty) {
      setCancelDialogOpen(true);
      return;
    }

    navigateBack();
  };

  if (!isValidId) {
    return <Box p={3}>Đang kiểm tra thông tin bài viết...</Box>;
  }

  if (postQuery?.isLoading) {
    return (
      <Box>
        <Skeleton variant="text" width={280} height={56} sx={{ mb: 2 }} />
        <Grid2 container spacing={3}>
          <Grid2 size={{ xs: 12, md: 8 }}>
            <Card sx={{ borderRadius: 1.75, border: 1, borderColor: "divider", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Skeleton variant="text" width={180} height={36} sx={{ mb: 2 }} />
                <Skeleton variant="rounded" height={56} sx={{ mb: 2 }} />
                <Skeleton variant="rounded" height={56} sx={{ mb: 2 }} />
                <Skeleton variant="rounded" height={360} />
              </CardContent>
            </Card>
          </Grid2>
          <Grid2 size={{ xs: 12, md: 4 }}>
            <Card sx={{ borderRadius: 1.75, border: 1, borderColor: "divider", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Skeleton variant="text" width={130} height={36} sx={{ mb: 2 }} />
                <Skeleton variant="rounded" height={56} sx={{ mb: 2 }} />
                <Skeleton variant="rounded" height={56} />
              </CardContent>
            </Card>
          </Grid2>
        </Grid2>
      </Box>
    );
  }

  if (postQuery?.isError || !postData) {
    return (
      <Box>
        <Alert severity="error" action={<Button onClick={navigateBack}>Quay về danh sách</Button>}>
          {getErrorMessage(postQuery?.error, "Không tìm thấy bài viết hoặc dữ liệu không khả dụng.")}
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      <Typography
        variant="h4"
        fontWeight={700}
        color="text.primary"
        letterSpacing="-0.02em"
        mb={3}
      >
        Chỉnh sửa bài viết
      </Typography>

      <Edit
        isLoading={formLoading || postQuery?.isLoading}
        title=""
        wrapperProps={{
          sx: {
            bgcolor: "transparent",
            boxShadow: "none",
            p: 0,
          },
        }}
        headerProps={{ sx: { display: "none" } }}
        footerButtons={
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: "100%" }}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={handleCancel}
              sx={{
                borderRadius: 1,
                textTransform: "none",
                fontWeight: 600,
                color: "text.secondary",
                borderColor: "divider",
              }}
            >
              Hủy
            </Button>
            <Button
              type="submit"
              variant="contained"
              color="success"
              disabled={formLoading || postQuery?.isLoading}
              onClick={handleSubmit((data) => onFinish(data))}
              sx={{
                borderRadius: 1,
                textTransform: "none",
                fontWeight: 600,
                bgcolor: "success.main",
                boxShadow: "none",
                "&:hover": { bgcolor: "success.dark", boxShadow: "none" },
              }}
            >
              Lưu thay đổi
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
          <Grid2 container spacing={3}>
            <Grid2 size={{ xs: 12, md: 8 }}>
              <Card
                sx={{
                  borderRadius: 1.75,
                  border: 1,
                  borderColor: "divider",
                  boxShadow: "none",
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight={600} mb={3} color="text.primary">
                    Thông tin cơ bản
                  </Typography>

                  <TextField
                    {...register("title", {
                      required: "Tiêu đề là bắt buộc",
                      maxLength: { value: 255, message: "Tiêu đề tối đa 255 ký tự" },
                    })}
                    error={Boolean(errors.title)}
                    helperText={errors.title?.message}
                    label="Tiêu đề (*)"
                    fullWidth
                    sx={{ mb: 3 }}
                    InputProps={{ sx: { borderRadius: 1 } }}
                  />

                  <TextField
                    {...register("slug", {
                      required: "Đường dẫn (slug) là bắt buộc",
                      maxLength: { value: 255, message: "Slug tối đa 255 ký tự" },
                      pattern: {
                        value: /^[a-z0-9-]+$/,
                        message: "Slug chỉ chứa chữ thường, số và dấu gạch ngang",
                      },
                    })}
                    error={Boolean(errors.slug)}
                    helperText={errors.slug?.message}
                    label="Đường dẫn (Slug) (*)"
                    fullWidth
                    sx={{ mb: 3 }}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: 1 } }}
                  />

                  <TextField
                    {...register("content", {
                      required: "Nội dung là bắt buộc",
                      validate: (value) => value.trim().length > 0 || "Nội dung là bắt buộc",
                    })}
                    error={Boolean(errors.content)}
                    helperText={errors.content?.message}
                    label="Nội dung (*)"
                    fullWidth
                    multiline
                    minRows={15}
                    InputProps={{ sx: { borderRadius: 1 } }}
                  />
                </CardContent>
              </Card>
            </Grid2>

            <Grid2 size={{ xs: 12, md: 4 }}>
              <Card
                sx={{
                  borderRadius: 1.75,
                  border: 1,
                  borderColor: "divider",
                  boxShadow: "none",
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight={600} mb={3} color="text.primary">
                    Phân loại
                  </Typography>

                  <TextField
                    value={typeCode}
                    label="Loại bài viết (Post Type)"
                    fullWidth
                    disabled
                    sx={{ mb: 3 }}
                    InputProps={{ sx: { borderRadius: 1 } }}
                  />

                  <Controller
                    control={control}
                    name="categoryId"
                    render={({ field }) => (
                      <Autocomplete<IHierarchicalCategoryOption>
                        options={categoryOptions}
                        loading={categoryQuery.isLoading}
                        value={categoryOptions.find((option) => option.id === field.value) ?? null}
                        onChange={(_, value) => field.onChange(value?.id ?? null)}
                        getOptionLabel={(option) => option.name}
                        isOptionEqualToValue={(option, value) => option.id === value.id}
                        noOptionsText={
                          categoryQuery.isError
                            ? "Không thể tải danh mục"
                            : "Chưa có danh mục phù hợp"
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
                            error={Boolean(errors.categoryId)}
                            helperText={errors.categoryId?.message}
                            InputProps={{
                              ...params.InputProps,
                              sx: { borderRadius: 1 },
                            }}
                          />
                        )}
                      />
                    )}
                  />
                </CardContent>
              </Card>
            </Grid2>
          </Grid2>
        </Box>
      </Edit>

      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Hủy thay đổi?</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Bạn có thay đổi chưa được lưu. Nếu rời trang, các thay đổi này sẽ bị mất.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)} color="secondary">
            Tiếp tục chỉnh sửa
          </Button>
          <Button onClick={navigateBack} color="primary" variant="contained">
            Rời trang
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

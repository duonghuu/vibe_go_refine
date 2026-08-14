import React, { useEffect } from "react";
import { Create } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { useSearchParams } from "react-router";
import {
  Box,
  Card,
  CardContent,
  TextField,
  Typography,
  Button,
  Stack,
} from "@mui/material";
import { HttpError, useNavigation, useNotification } from "@refinedev/core";

export interface IPostResponse {
  id: number;
  typeId: number;
  title: string;
  slug: string;
  content: string;
  createdAt: string;
}

export interface ICreatePostRequest {
  typeId: number;
  title: string;
  slug: string;
  content: string;
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
  const typeIdParam = searchParams.get("type_id");
  const typeId = typeIdParam ? parseInt(typeIdParam, 10) : null;
  
  const { goBack } = useNavigation();
  const { open } = useNotification();

  useEffect(() => {
    if (!typeId || isNaN(typeId)) {
      open?.({
        type: "error",
        message: "Lỗi truy cập",
        description: "Thiếu thông tin loại bài viết (type_id). Đang chuyển hướng...",
      });
      setTimeout(() => {
        goBack();
      }, 2000);
    }
  }, [typeId, goBack, open]);

  const {
    refineCore: { formLoading, onFinish },
    register,
    handleSubmit,
    setValue,
    formState: { errors },
    watch,
  } = useForm<IPostResponse, HttpError, ICreatePostRequest>({
    refineCoreProps: {
      action: "create",
      resource: "posts",
      redirect: false,
      onMutationSuccess: () => {
        goBack(); // Quay lại trang danh sách bài viết theo type
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const titleValue = watch("title");

  useEffect(() => {
    if (titleValue) {
      setValue("slug", generateSlug(titleValue), { shouldValidate: true });
    } else {
      setValue("slug", "");
    }
  }, [titleValue, setValue]);

  const onCustomSubmit = (data: ICreatePostRequest) => {
    if (!typeId || isNaN(typeId)) return;
    
    onFinish({
      ...data,
      typeId: typeId,
    });
  };

  if (!typeId || isNaN(typeId)) {
    return <Box p={3}>Đang kiểm tra thông tin loại bài viết...</Box>;
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
              onClick={() => goBack()}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", color: 'text.secondary', borderColor: '#D5D5D5' }}
            >
              Hủy
            </Button>
            <Button
              variant="contained"
              color="primary"
              disabled={formLoading}
              onClick={handleSubmit(onCustomSubmit)}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'primary.main', boxShadow: 'none', '&:hover': { bgcolor: 'primary.dark', boxShadow: 'none' } }}
            >
              Lưu bài viết
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
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
        </Box>
      </Create>
    </Box>
  );
};

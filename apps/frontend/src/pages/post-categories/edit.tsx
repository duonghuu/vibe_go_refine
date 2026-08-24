import React, { useEffect, useState, useMemo } from "react";
import {
  Box,
  TextField,
  Typography,
  Grid2,
  Card,
  CardContent,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  RadioGroup,
  FormControlLabel,
  Radio,
  Button,
  Stack,
  IconButton,
} from "@mui/material";
import { Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import DeleteIcon from "@mui/icons-material/Delete";
import { HttpError, useSelect, CrudFilter } from "@refinedev/core";
import { API_URL, BACKEND_URL } from "../../providers/constants";

interface IPostCategoryResponse {
  id: number;
  name: string;
  slug: string;
  parentId: number | null;
  description: string;
  imageUrl: string;
  sortOrder: number;
  status: string;
  typeCode: string;
  createdAt: string;
}

interface IEditPostCategoryRequest {
  name: string;
  slug: string;
  parentId: number | null;
  description: string;
  imageUrl: string;
  sortOrder: number;
  status: string;
}

export const PostCategoryEdit = () => {
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const {
    saveButtonProps,
    refineCore: { formLoading, query: categoryQuery },
    register,
    control,
    setValue,
    watch,
    formState: { errors },
  } = useForm<IPostCategoryResponse, HttpError, IEditPostCategoryRequest>({
    refineCoreProps: {
      action: "edit",
      resource: "post-categories",
      redirect: "list",
      meta: {
        method: "put",
      },
    },
  });

  const categoryData = categoryQuery?.data?.data;
  const typeCode = categoryData?.typeCode;
  const currentId = categoryData?.id;
  const initialImageUrl = categoryData?.imageUrl;

  useEffect(() => {
    if (initialImageUrl) {
      setPreviewImage(initialImageUrl);
    }
  }, [initialImageUrl]);

  const permanentFilters = useMemo<CrudFilter[]>(() => {
    return typeCode
      ? [
          {
            field: "typeCode",
            operator: "eq",
            value: typeCode,
          },
        ]
      : [];
  }, [typeCode]);

  const { options, query: parentQuery } = useSelect<IPostCategoryResponse, HttpError>({
    resource: "post-categories/tree",
    optionLabel: "name",
    optionValue: "id",
    filters: permanentFilters,
    queryOptions: {
      enabled: !!typeCode,
    },
  });

  // Filter out the current category and its children to prevent circular reference
  const filteredOptions = useMemo(() => {
    const rawData = (parentQuery?.data?.data ?? []) as unknown as IPostCategoryResponse[];
    const getDescendants = (parentId: number): number[] => {
      const children = rawData.filter((item) => item.parentId === parentId);
      return children.flatMap((child) => [child.id, ...getDescendants(child.id)]);
    };

    const invalidIds = currentId ? [currentId, ...getDescendants(currentId)] : [];
    const loadedOptions = rawData.length > 0
      ? rawData.map((item) => ({ label: item.name, value: item.id }))
      : options;

    return loadedOptions.filter((option) => !invalidIds.includes(Number(option.value)));
  }, [currentId, options, parentQuery?.data?.data]);

  const handleImageChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const formData = new FormData();
      formData.append("file", file);

      setIsUploading(true);
      try {
        const response = await fetch(`${API_URL}/media/upload`, {
          method: "POST",
          body: formData,
        });
        const data = await response.json();
        const url = data?.data?.originalUrl || data?.originalUrl;
        if (url) {
          setPreviewImage(url);
          setValue("imageUrl", url, { shouldValidate: true, shouldDirty: true });
        }
      } catch (error) {
        console.error("Upload failed", error);
      } finally {
        setIsUploading(false);
      }
    }
  };

  const handleRemoveImage = () => {
    setPreviewImage(null);
    setValue("imageUrl", "", { shouldValidate: true, shouldDirty: true });
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Chỉnh sửa danh mục bài viết
        </Typography>
      </Box>

      <Edit
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
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: "100%" }}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={() => window.history.back()}
              sx={{
                borderRadius: "8px",
                textTransform: "none",
                fontWeight: "600",
                color: "text.secondary",
                borderColor: "#D5D5D5",
              }}
            >
              Hủy
            </Button>
            <Button
              {...saveButtonProps}
              variant="contained"
              color="success"
              sx={{
                borderRadius: "8px",
                textTransform: "none",
                fontWeight: "600",
                bgcolor: "success.main",
                boxShadow: "none",
                "&:hover": { bgcolor: "success.dark", boxShadow: "none" },
              }}
            >
              Lưu
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
          <Grid2 container spacing={3}>
            {/* Cột chính: Thông tin cơ bản */}
            <Grid2 size={{ xs: 12, md: 8 }}>
              <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none", height: "100%" }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                    Thông tin cơ bản
                  </Typography>

                  <TextField
                    value={typeCode || ""}
                    margin="normal"
                    fullWidth
                    label="Loại bài viết (Post Type)"
                    disabled
                    InputProps={{ sx: { borderRadius: "8px" } }}
                    sx={{ mb: 3 }}
                  />

                  <TextField
                    {...register("name", {
                      required: "Tên danh mục là bắt buộc",
                    })}
                    error={!!(errors as any)?.name}
                    helperText={(errors as any)?.name?.message}
                    margin="normal"
                    fullWidth
                    label="Tên danh mục (*)"
                    name="name"
                    sx={{ mb: 3 }}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />

                  <TextField
                    {...register("slug", {
                      required: "Slug là bắt buộc",
                      pattern: {
                        value: /^[a-z0-9-]+$/,
                        message: "Slug chỉ chứa chữ thường, số và dấu gạch ngang",
                      },
                    })}
                    error={!!(errors as any)?.slug}
                    helperText={(errors as any)?.slug?.message}
                    margin="normal"
                    fullWidth
                    label="Slug (*)"
                    name="slug"
                    sx={{ mb: 3 }}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />

                  <Controller
                    control={control}
                    name="parentId"
                    render={({ field }) => {
                      const currentValue = field.value ?? "";
                      return (
                        <FormControl fullWidth sx={{ mb: 3 }}>
                          <InputLabel id="parent-category-label" shrink>Danh mục cha</InputLabel>
                          <Select
                            {...field}
                            labelId="parent-category-label"
                            label="Danh mục cha"
                            displayEmpty
                            value={currentValue}
                            onChange={(e) => field.onChange(e.target.value === "" ? null : e.target.value)}
                            sx={{ borderRadius: "8px" }}
                            notched
                          >
                            {parentQuery.isLoading ? (
                              <MenuItem disabled value="">
                                <em>Đang tải...</em>
                              </MenuItem>
                            ) : (
                              [
                                <MenuItem key="empty" value="">
                                  <em>Không có (Danh mục gốc)</em>
                                </MenuItem>,
                                ...filteredOptions.map((option) => (
                                  <MenuItem key={option.value} value={option.value}>
                                    {option.label}
                                  </MenuItem>
                                ))
                              ]
                            )}
                          </Select>
                        </FormControl>
                      )
                    }}
                  />

                  <TextField
                    {...register("description")}
                    margin="normal"
                    fullWidth
                    multiline
                    rows={4}
                    label="Mô tả"
                    name="description"
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />
                </CardContent>
              </Card>
            </Grid2>

            {/* Cột phụ: Media & Settings */}
            <Grid2 size={{ xs: 12, md: 4 }}>
              <Stack spacing={3}>
                {/* Hình ảnh */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Hình ảnh
                    </Typography>

                    <Box
                      sx={{
                        border: "2px dashed #D5D5D5",
                        borderRadius: "8px",
                        p: 3,
                        textAlign: "center",
                        bgcolor: "#FAFAFA",
                        position: "relative",
                        minHeight: "200px",
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "center",
                        justifyContent: "center",
                      }}
                    >
                      {previewImage ? (
                        <Box position="relative" display="inline-block" width="100%">
                          <img
                            src={previewImage?.startsWith("/") ? `${BACKEND_URL}${previewImage}` : previewImage}
                            alt="Preview"
                            style={{ width: "100%", height: "200px", objectFit: "cover", borderRadius: "8px" }}
                          />
                          <IconButton
                            size="small"
                            onClick={handleRemoveImage}
                            sx={{
                              position: "absolute",
                              top: 8,
                              right: 8,
                              bgcolor: "rgba(255, 255, 255, 0.9)",
                              "&:hover": { bgcolor: "white" },
                            }}
                          >
                            <DeleteIcon color="error" />
                          </IconButton>
                        </Box>
                      ) : (
                        <>
                          <CloudUploadIcon sx={{ fontSize: 48, color: "text.secondary", mb: 1 }} />
                          <Typography variant="body2" color="text.secondary" mb={2}>
                            Kéo thả hoặc click để chọn ảnh
                          </Typography>
                          <Button
                            variant="contained"
                            component="label"
                            color="primary"
                            sx={{ textTransform: "none", borderRadius: "8px", boxShadow: "none" }}
                            disabled={isUploading}
                          >
                            {isUploading ? "Đang tải lên..." : "Tải ảnh lên"}
                            <input
                              type="file"
                              hidden
                              accept="image/jpeg, image/png, image/webp"
                              onChange={handleImageChange}
                            />
                          </Button>
                          <Typography variant="caption" display="block" color="text.secondary" mt={2}>
                            Định dạng: JPG, PNG, WEBP (≤ 5MB)
                          </Typography>
                        </>
                      )}
                    </Box>
                    {/* Hidden field for react-hook-form */}
                    <input type="hidden" {...register("imageUrl")} />
                  </CardContent>
                </Card>

                {/* Thiết lập */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Thiết lập
                    </Typography>

                    <Controller
                      control={control}
                      name="status"
                      defaultValue="ACTIVE"
                      render={({ field }) => (
                        <FormControl component="fieldset" sx={{ mb: 3, width: "100%" }}>
                          <Typography variant="body2" fontWeight="600" color="text.secondary" mb={1}>
                            Trạng thái
                          </Typography>
                          <RadioGroup {...field} row>
                            <FormControlLabel value="ACTIVE" control={<Radio color="primary" />} label="Hoạt động" />
                            <FormControlLabel value="INACTIVE" control={<Radio color="primary" />} label="Ẩn" />
                          </RadioGroup>
                        </FormControl>
                      )}
                    />

                    <TextField
                      {...register("sortOrder", {
                        valueAsNumber: true,
                        min: { value: 0, message: "Thứ tự phải ≥ 0" },
                      })}
                      error={!!(errors as any)?.sortOrder}
                      helperText={(errors as any)?.sortOrder?.message}
                      type="number"
                      fullWidth
                      label="Thứ tự hiển thị"
                      InputLabelProps={{ shrink: true }}
                      InputProps={{ sx: { borderRadius: "8px" } }}
                    />
                  </CardContent>
                </Card>
              </Stack>
            </Grid2>
          </Grid2>
        </Box>
      </Edit>
    </Box>
  );
};

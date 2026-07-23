import React, { useState, useEffect } from "react";
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  List,
  ListItem,
  ListItemText,
} from "@mui/material";
import { Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import DeleteIcon from "@mui/icons-material/Delete";
import { HttpError, useSelect, useParsed, useCustomMutation, useNotification } from "@refinedev/core";
import { API_URL, BACKEND_URL } from "../../providers/constants";

export interface ICategoryResponse {
  id: number;
  name: string;
  slug: string;
  parentId: number | null;
  description: string;
  imageUrl: string;
  sortOrder: number;
  status: string;
  productCount: number;
  createdAt: string;
  updatedAt?: string;
  createdBy?: string;
}

export interface IUpdateCategoryRequest {
  name: string;
  slug: string;
  parentId: number | null;
  description: string;
  imageUrl: string;
  sortOrder: number;
  status: string;
}

export const CategoryEdit = () => {
  const { id } = useParsed();
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [showStatusWarning, setShowStatusWarning] = useState(false);
  const [pendingStatus, setPendingStatus] = useState<string | null>(null);

  const { open } = useNotification();
  const [isUploading, setIsUploading] = useState(false);

  const {
    saveButtonProps,
    refineCore: { formLoading, query, onFinish },
    register,
    control,
    setValue,
    watch,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm<ICategoryResponse, HttpError, IUpdateCategoryRequest>({
    refineCoreProps: {
      action: "edit",
      redirect: "list",
      meta: {
        method: "put",
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const categoryData = query?.data?.data as unknown as ICategoryResponse;
  
  // Set initial image preview
  useEffect(() => {
    if (categoryData?.imageUrl) {
      setPreviewImage(categoryData.imageUrl);
    }
  }, [categoryData]);

  const { options } = useSelect<ICategoryResponse, HttpError>({
    resource: "categories",
    optionLabel: "name",
    optionValue: "id",
  });

  // Filter out the current category and its children to prevent circular dependency
  // For now, we filter out the current category itself.
  const parentOptions = options.filter((option) => String(option.value) !== String(id));

  const currentStatus = watch("status");

  const handleStatusChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newStatus = e.target.value;
    if (categoryData?.status === "ACTIVE" && newStatus === "HIDDEN") {
      setPendingStatus(newStatus);
      setShowStatusWarning(true);
    } else {
      setValue("status", newStatus, { shouldValidate: true, shouldDirty: true });
    }
  };

  const confirmStatusChange = () => {
    if (pendingStatus) {
      setValue("status", pendingStatus, { shouldValidate: true, shouldDirty: true });
    }
    setShowStatusWarning(false);
    setPendingStatus(null);
  };

  const cancelStatusChange = () => {
    setValue("status", categoryData?.status || "ACTIVE");
    setShowStatusWarning(false);
    setPendingStatus(null);
  };

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
          Chỉnh sửa danh mục
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
          <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: '100%' }}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={() => {
                if (isDirty) {
                  if (window.confirm("Bạn có thay đổi chưa lưu. Bạn có chắc chắn muốn hủy?")) {
                    window.history.back();
                  }
                } else {
                  window.history.back();
                }
              }}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", color: 'text.secondary', borderColor: '#D5D5D5' }}
            >
              Hủy
            </Button>
            <Button
              disabled={formLoading}
              onClick={handleSubmit(
                (data) => onFinish(data),
                () => {
                  open?.({
                    type: "error",
                    message: "Lưu thất bại",
                    description: "Vui lòng kiểm tra lại các trường thông tin không hợp lệ.",
                    key: "validation-error",
                  });
                }
              )}
              variant="contained"
              color="success"
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'success.main', boxShadow: 'none', '&:hover': { bgcolor: 'success.dark', boxShadow: 'none' } }}
            >
              Lưu thay đổi
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
          <Grid2 container spacing={3}>
            {/* Cột chính: Thông tin cơ bản */}
            <Grid2 size={{ xs: 12, md: 8 }}>
              <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none", height: "100%", mb: 3 }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                    Thông tin cơ bản
                  </Typography>

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
                    label="Slug"
                    name="slug"
                    sx={{ mb: 3 }}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />

                  <Controller
                    control={control}
                    name="parentId"
                    render={({ field }) => (
                      <FormControl fullWidth sx={{ mb: 3 }}>
                        <InputLabel id="parent-category-label" shrink>Danh mục cha</InputLabel>
                        <Select
                          {...field}
                          labelId="parent-category-label"
                          label="Danh mục cha"
                          displayEmpty
                          value={field.value || ""}
                          sx={{ borderRadius: "8px" }}
                        >
                          <MenuItem value="">
                            <em>Không có (Danh mục gốc)</em>
                          </MenuItem>
                          {parentOptions.map((option) => (
                            <MenuItem key={option.value} value={option.value}>
                              {option.label}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    )}
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
                        justifyContent: "center"
                      }}
                    >
                      {previewImage ? (
                        <Box position="relative" display="inline-block" width="100%">
                          <img src={previewImage?.startsWith("/") ? `${BACKEND_URL}${previewImage}` : previewImage} alt="Preview" style={{ width: "100%", height: "200px", objectFit: "cover", borderRadius: "8px" }} />
                          <IconButton
                            size="small"
                            onClick={handleRemoveImage}
                            sx={{
                              position: "absolute",
                              top: 8,
                              right: 8,
                              bgcolor: "rgba(255, 255, 255, 0.9)",
                              '&:hover': { bgcolor: "white" }
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
                          <Button variant="contained" component="label" color="primary" sx={{ textTransform: "none", borderRadius: "8px", boxShadow: "none" }} disabled={isUploading}>
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
                          <RadioGroup 
                            row 
                            name="status"
                            value={field.value || "ACTIVE"}
                            onChange={handleStatusChange}
                          >
                            <FormControlLabel value="ACTIVE" control={<Radio color="primary" />} label="Hoạt động" />
                            <FormControlLabel value="HIDDEN" control={<Radio color="primary" />} label="Ẩn" />
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

                {/* Thông tin hệ thống (Chỉ đọc) */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #E0E0E0", boxShadow: "none", bgcolor: "#FAFAFA" }}>
                  <CardContent sx={{ p: 2, pb: "16px !important" }}>
                    <Typography variant="subtitle2" fontWeight="600" mb={1} color="text.secondary">
                      Thông tin hệ thống
                    </Typography>
                    <List dense sx={{ p: 0 }}>
                      <ListItem sx={{ px: 0, py: 0.5 }}>
                        <ListItemText 
                          primary={<Typography variant="caption" color="text.secondary">ID</Typography>} 
                          secondary={<Typography variant="body2" fontWeight="500">{categoryData?.id || "-"}</Typography>} 
                        />
                      </ListItem>
                      <ListItem sx={{ px: 0, py: 0.5 }}>
                        <ListItemText 
                          primary={<Typography variant="caption" color="text.secondary">Ngày tạo</Typography>} 
                          secondary={<Typography variant="body2" fontWeight="500">{categoryData?.createdAt || "-"}</Typography>} 
                        />
                      </ListItem>
                      <ListItem sx={{ px: 0, py: 0.5 }}>
                        <ListItemText 
                          primary={<Typography variant="caption" color="text.secondary">Ngày cập nhật</Typography>} 
                          secondary={<Typography variant="body2" fontWeight="500">{categoryData?.updatedAt || "-"}</Typography>} 
                        />
                      </ListItem>
                      <ListItem sx={{ px: 0, py: 0.5 }}>
                        <ListItemText 
                          primary={<Typography variant="caption" color="text.secondary">Người tạo</Typography>} 
                          secondary={<Typography variant="body2" fontWeight="500">{categoryData?.createdBy || "-"}</Typography>} 
                        />
                      </ListItem>
                    </List>
                  </CardContent>
                </Card>
              </Stack>
            </Grid2>
          </Grid2>
        </Box>
      </Edit>

      {/* Dialog Cảnh báo chuyển trạng thái */}
      <Dialog
        open={showStatusWarning}
        onClose={cancelStatusChange}
      >
        <DialogTitle>Xác nhận ẩn danh mục</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Việc ẩn danh mục này có thể ảnh hưởng đến sản phẩm và danh mục con bên trong. Bạn có chắc chắn muốn tiếp tục?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={cancelStatusChange} color="secondary">
            Hủy
          </Button>
          <Button onClick={confirmStatusChange} color="primary" variant="contained" autoFocus>
            Đồng ý
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

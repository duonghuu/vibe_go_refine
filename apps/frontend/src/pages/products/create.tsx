import React, { useState, useRef } from "react";
import { Create, useAutocomplete } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import {
  Box,
  Card,
  CardContent,
  Grid2,
  TextField,
  Typography,
  RadioGroup,
  FormControlLabel,
  Radio,
  FormControl,
  Autocomplete,
  Button,
  Stack,
  IconButton,
} from "@mui/material";
import { HttpError, useNavigation } from "@refinedev/core";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import DeleteIcon from "@mui/icons-material/Delete";
import { API_URL, BACKEND_URL } from "../../providers/constants";

export interface ICategoryResponse {
  id: number;
  name: string;
}

export interface IProductResponse {
  id: number;
  name: string;
  description: string;
  categoryId: number;
  price: number;
  salePrice?: number;
  sku: string;
  stock: number;
  soldCount: number;
  image: string;
  status: string;
  createdAt: string;
}

export interface ICreateProductRequest {
  name: string;
  description: string;
  categoryId: number;
  price: number;
  salePrice?: number;
  sku: string;
  stock: number;
  image: string;
  status: string;
}

export const ProductCreate: React.FC = () => {
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const isSaveAndNewRef = useRef(false);
  const { list } = useNavigation();

  const {
    saveButtonProps,
    refineCore: { formLoading, onFinish },
    register,
    control,
    setValue,
    reset,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<IProductResponse, HttpError, ICreateProductRequest>({
    refineCoreProps: {
      action: "create",
      resource: "products",
      redirect: false,
      onMutationSuccess: () => {
        if (isSaveAndNewRef.current) {
          reset();
          setPreviewImage(null);
          isSaveAndNewRef.current = false;
        } else {
          list("products");
        }
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const { autocompleteProps } = useAutocomplete<ICategoryResponse, HttpError>({
    resource: "categories",
  });

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
          setValue("image", url, { shouldValidate: true });
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
    setValue("image", "", { shouldValidate: true });
  };

  const onCustomSubmit = (data: ICreateProductRequest, customStatus: string) => {
    onFinish({
      ...data,
      status: customStatus,
    });
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Thêm sản phẩm mới
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
              onClick={() => window.history.back()}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", color: 'text.secondary', borderColor: '#D5D5D5' }}
            >
              Hủy
            </Button>
            <Button
              variant="contained"
              disabled={formLoading}
              onClick={handleSubmit((data) => {
                isSaveAndNewRef.current = false;
                onCustomSubmit(data, "HIDDEN");
              })}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'primary.main', boxShadow: 'none' }}
            >
              Lưu & Ẩn
            </Button>
            <Button
              variant="contained"
              color="success"
              onClick={handleSubmit((data) => {
                isSaveAndNewRef.current = false;
                onCustomSubmit(data, "ACTIVE");
              })}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'success.main', boxShadow: 'none', '&:hover': { bgcolor: 'success.dark', boxShadow: 'none' } }}
            >
              Lưu & Xuất bản
            </Button>
          </Stack>
        }
      >
        <Box component="form" autoComplete="off">
          <Grid2 container spacing={3}>
            {/* Main Column */}
            <Grid2 size={{ xs: 12, md: 8 }}>
              <Stack spacing={3}>
                {/* Basic Info Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                      Thông tin cơ bản
                    </Typography>
                    <TextField
                      {...register("name", { required: "Tên sản phẩm là bắt buộc" })}
                      error={!!errors.name}
                      helperText={errors.name?.message as string}
                      label="Tên sản phẩm (*)"
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      sx={{ mb: 3 }}
                      InputProps={{ sx: { borderRadius: "8px" } }}
                    />
                    <TextField
                      {...register("description")}
                      label="Mô tả chi tiết"
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      multiline
                      rows={5}
                      InputProps={{ sx: { borderRadius: "8px" } }}
                    />
                  </CardContent>
                </Card>

                {/* Pricing Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                      Giá cả
                    </Typography>
                    <Grid2 container spacing={2}>
                      <Grid2 size={{ xs: 12, sm: 6 }}>
                        <TextField
                          {...register("price", {
                            required: "Giá gốc là bắt buộc",
                            valueAsNumber: true,
                            min: { value: 1, message: "Giá gốc phải lớn hơn 0" },
                          })}
                          error={!!errors.price}
                          helperText={errors.price?.message as string}
                          label="Giá bán gốc (*)"
                          type="number"
                          margin="normal"
                          variant="outlined"
                          fullWidth
                          InputProps={{ sx: { borderRadius: "8px" } }}
                        />
                      </Grid2>
                      <Grid2 size={{ xs: 12, sm: 6 }}>
                        <TextField
                          {...register("salePrice", {
                            valueAsNumber: true,
                            validate: (value) => {
                              if (!value) return true;
                              const price = watch("price");
                              if (price && Number(value) >= Number(price)) {
                                return "Giá khuyến mãi phải nhỏ hơn Giá gốc";
                              }
                              return true;
                            }
                          })}
                          error={!!errors.salePrice}
                          helperText={errors.salePrice?.message as string}
                          label="Giá khuyến mãi"
                          type="number"
                          margin="normal"
                          variant="outlined"
                          fullWidth
                          InputProps={{ sx: { borderRadius: "8px" } }}
                        />
                      </Grid2>
                    </Grid2>
                  </CardContent>
                </Card>

                {/* Inventory Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                      Quản lý kho
                    </Typography>
                    <Grid2 container spacing={2}>
                      <Grid2 size={{ xs: 12, sm: 6 }}>
                        <TextField
                          {...register("sku")}
                          label="Mã sản phẩm (SKU)"
                          margin="normal"
                          variant="outlined"
                          fullWidth
                          InputProps={{ sx: { borderRadius: "8px" } }}
                        />
                      </Grid2>
                      <Grid2 size={{ xs: 12, sm: 6 }}>
                        <TextField
                          {...register("stock", {
                            valueAsNumber: true,
                            min: { value: 0, message: "Tồn kho không được âm" }
                          })}
                          error={!!errors.stock}
                          helperText={errors.stock?.message as string}
                          label="Số lượng tồn kho"
                          type="number"
                          defaultValue={0}
                          margin="normal"
                          variant="outlined"
                          fullWidth
                          InputProps={{ sx: { borderRadius: "8px" } }}
                        />
                      </Grid2>
                    </Grid2>
                  </CardContent>
                </Card>
              </Stack>
            </Grid2>

            {/* Side Column */}
            <Grid2 size={{ xs: 12, md: 4 }}>
              <Stack spacing={3}>
                {/* Media Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Hình ảnh đại diện
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
                    <input type="hidden" {...register("image")} />
                  </CardContent>
                </Card>

                {/* Organization Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                      Phân loại
                    </Typography>
                    <Controller
                      control={control}
                      name="categoryId"
                      rules={{ required: "Danh mục là bắt buộc" }}
                      render={({ field }) => (
                        <Autocomplete
                          {...autocompleteProps}
                          value={
                            autocompleteProps?.options?.find(
                              (item) => item.id === field.value
                            ) || null
                          }
                          getOptionLabel={(item) => item.name}
                          isOptionEqualToValue={(option, value) =>
                            value === undefined || option?.id?.toString() === value?.id?.toString()
                          }
                          onChange={(_, value) => {
                            field.onChange(value?.id ?? "");
                          }}
                          renderInput={(params) => (
                            <TextField
                              {...params}
                              label="Danh mục (*)"
                              variant="outlined"
                              error={!!errors.categoryId}
                              helperText={errors.categoryId?.message as string}
                              InputProps={{
                                ...params.InputProps,
                                sx: { borderRadius: "8px" }
                              }}
                            />
                          )}
                        />
                      )}
                    />
                  </CardContent>
                </Card>

                {/* Status Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Trạng thái
                    </Typography>
                    <Controller
                      control={control}
                      name="status"
                      defaultValue="ACTIVE"
                      render={({ field }) => (
                        <FormControl component="fieldset" sx={{ mb: 3, width: "100%" }}>
                          <RadioGroup {...field} row>
                            <FormControlLabel value="ACTIVE" control={<Radio color="primary" />} label="Hoạt động" />
                            <FormControlLabel value="HIDDEN" control={<Radio color="primary" />} label="Ẩn" />
                          </RadioGroup>
                        </FormControl>
                      )}
                    />
                  </CardContent>
                </Card>
              </Stack>
            </Grid2>
          </Grid2>
        </Box>
      </Create>
    </Box>
  );
};

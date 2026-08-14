import { Card, CardContent, FormControl, FormControlLabel, FormHelperText, Grid, Radio, RadioGroup, TextField, Typography } from "@mui/material";
import { UseFormReturnType } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import { HttpError } from "@refinedev/core";

export interface IPostType {
  id: number;
  code: string;
  name: string;
  status: string;
  sort_order: number;
  createdAt: string;
}

export interface IPostTypeForm {
  code: string;
  name: string;
  status: "ACTIVE" | "INACTIVE";
  sort_order: number;
}

interface PostTypeFormProps {
  form: UseFormReturnType<IPostType, HttpError, IPostTypeForm>;
  isEdit?: boolean;
}

export const PostTypeForm = ({ form, isEdit = false }: PostTypeFormProps) => {
  const {
    register,
    control,
    formState: { errors },
  } = form;

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} md={8}>
        <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
          <CardContent sx={{ p: 3 }}>
            <Typography variant="h6" fontWeight="bold" mb={3}>Thông tin cơ bản</Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <TextField
                  {...register("code", {
                    required: "Mã loại bài viết là bắt buộc",
                    pattern: {
                      value: /^[A-Z0-9_]+$/,
                      message: "Chỉ cho phép chữ hoa, số và dấu gạch dưới",
                    },
                    maxLength: {
                      value: 50,
                      message: "Tối đa 50 ký tự",
                    }
                  })}
                  error={!!(errors as any)?.code}
                  helperText={(errors as any)?.code?.message}
                  margin="normal"
                  fullWidth
                  InputLabelProps={{ shrink: true }}
                  type="text"
                  label="Mã loại bài viết (Code) *"
                  name="code"
                  disabled={isEdit} // Disable when editing
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <TextField
                  {...register("name", {
                    required: "Tên loại bài viết là bắt buộc",
                    maxLength: {
                      value: 100,
                      message: "Tối đa 100 ký tự",
                    }
                  })}
                  error={!!(errors as any)?.name}
                  helperText={(errors as any)?.name?.message}
                  margin="normal"
                  fullWidth
                  InputLabelProps={{ shrink: true }}
                  type="text"
                  label="Tên loại (Name) *"
                  name="name"
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <TextField
                  {...register("sort_order", {
                    required: "Thứ tự hiển thị là bắt buộc",
                    valueAsNumber: true,
                    min: {
                      value: 1,
                      message: "Thứ tự phải lớn hơn 0"
                    }
                  })}
                  error={!!(errors as any)?.sort_order}
                  helperText={(errors as any)?.sort_order?.message}
                  margin="normal"
                  fullWidth
                  InputLabelProps={{ shrink: true }}
                  type="number"
                  label="Thứ tự (Sort Order) *"
                  name="sort_order"
                />
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </Grid>
      
      <Grid item xs={12} md={4}>
        <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none", height: "100%" }}>
          <CardContent sx={{ p: 3 }}>
            <Typography variant="h6" fontWeight="bold" mb={3}>Trạng thái</Typography>
            <FormControl error={!!(errors as any)?.status}>
              <Controller
                control={control}
                name="status"
                defaultValue="ACTIVE"
                rules={{ required: "Vui lòng chọn trạng thái" }}
                render={({ field }) => (
                  <RadioGroup {...field}>
                    <FormControlLabel value="ACTIVE" control={<Radio />} label="Hoạt động (ACTIVE)" />
                    <FormControlLabel value="INACTIVE" control={<Radio />} label="Tạm ẩn (INACTIVE)" />
                  </RadioGroup>
                )}
              />
              {(errors as any)?.status && (
                <FormHelperText>{(errors as any)?.status?.message}</FormHelperText>
              )}
            </FormControl>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
};

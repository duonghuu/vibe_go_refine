import React from "react";
import { Edit } from "@refinedev/mui";
import { 
  Box, 
  Card, 
  CardContent, 
  Grid2, 
  Typography, 
  TextField, 
  MenuItem, 
  FormControl, 
  InputLabel, 
  Select, 
  Chip,
  Alert
} from "@mui/material";
import { useForm } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import { HttpError, useNavigation } from "@refinedev/core";
import { IUser } from "./list";

export const UserEdit: React.FC = () => {
  const { list } = useNavigation();

  const {
    saveButtonProps,
    register,
    control,
    formState: { errors },
    refineCore: { query },
  } = useForm<IUser, HttpError, any>({
    refineCoreProps: {
      action: "edit",
      resource: "users",
      redirect: false,
      onMutationSuccess: () => {
        list("users");
      },
      meta: {
        method: "put",
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const userData = query?.data?.data;

  return (
    <Edit 
      saveButtonProps={saveButtonProps}
      headerProps={{
        sx: {
          bgcolor: 'transparent',
          px: 0,
          pt: 0,
          mb: 3
        }
      }}
      wrapperProps={{
        sx: {
          bgcolor: 'transparent',
          boxShadow: 'none',
          p: 0,
        }
      }}
    >
      <Grid2 container spacing={3}>
        {/* Main Column */}
        <Grid2 size={{ xs: 12, md: 8 }}>
          <Box display="flex" flexDirection="column" gap={3}>
            {/* Card 1: Basic Information */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Thông tin cơ bản
                </Typography>
                
                <Box display="flex" flexDirection="column" gap={3}>
                  <TextField
                    {...register("name", { required: "Tên là bắt buộc" })}
                    label="Họ và tên (*)"
                    error={!!errors.name}
                    helperText={errors.name?.message as string}
                    fullWidth
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />

                  <TextField
                    {...register("email")}
                    label="Email"
                    fullWidth
                    disabled={true}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{ sx: { borderRadius: "8px" } }}
                  />
                </Box>
              </CardContent>
            </Card>

            {/* Card 2: Permission */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Phân quyền
                </Typography>
                
                <Controller
                  control={control}
                  name="role"
                  render={({ field }) => (
                    <Box display="flex" flexDirection="column" gap={2}>
                      <FormControl fullWidth>
                        <InputLabel shrink>Vai trò</InputLabel>
                        <Select {...field} label="Vai trò" sx={{ borderRadius: "8px" }} value={field.value ?? ""} displayEmpty>
                          <MenuItem value="ADMIN">ADMIN</MenuItem>
                          <MenuItem value="STAFF">STAFF</MenuItem>
                          <MenuItem value="CUSTOMER">CUSTOMER</MenuItem>
                        </Select>
                      </FormControl>
                      {field.value === "ADMIN" && (
                        <Alert severity="warning" sx={{ borderRadius: "8px" }}>
                          Bạn đang cấp quyền Quản trị viên (ADMIN) cho tài khoản này. Người dùng sẽ có toàn quyền truy cập hệ thống.
                        </Alert>
                      )}
                    </Box>
                  )}
                />
              </CardContent>
            </Card>

            {/* Card 3: Status */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Trạng thái tài khoản
                </Typography>
                
                <Controller
                  control={control}
                  name="status"
                  render={({ field }) => (
                    <Box display="flex" flexDirection="column" gap={2}>
                      <FormControl fullWidth>
                        <InputLabel shrink>Trạng thái</InputLabel>
                        <Select {...field} label="Trạng thái" sx={{ borderRadius: "8px" }} value={field.value ?? ""}  displayEmpty>
                          <MenuItem value="ACTIVE">Hoạt động (ACTIVE)</MenuItem>
                          <MenuItem value="INACTIVE">Khóa (INACTIVE)</MenuItem>
                        </Select>
                      </FormControl>
                      {field.value === "INACTIVE" && (
                        <Alert severity="error" sx={{ borderRadius: "8px" }}>
                          Tài khoản bị khóa sẽ không thể đăng nhập vào hệ thống.
                        </Alert>
                      )}
                    </Box>
                  )}
                />
              </CardContent>
            </Card>
          </Box>
        </Grid2>

        {/* Side Column */}
        <Grid2 size={{ xs: 12, md: 4 }}>
          <Box display="flex" flexDirection="column" gap={3}>
            {/* Card 4: Current Status */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Trạng thái hiện tại
                </Typography>
                <Box display="flex" gap={2}>
                  {userData?.role && (
                    <Chip
                      label={userData.role}
                      color={
                        userData.role === "ADMIN" ? "info" :
                        userData.role === "STAFF" ? "secondary" : "default"
                      }
                      sx={{ fontWeight: "bold" }}
                    />
                  )}
                  {userData?.status && (
                    <Chip
                      label={userData.status}
                      color={userData.status === "ACTIVE" ? "success" : "error"}
                      variant="outlined"
                      sx={{ fontWeight: "bold" }}
                    />
                  )}
                </Box>
              </CardContent>
            </Card>

            {/* Card 5: System Information */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Thông tin hệ thống
                </Typography>
                <Box display="flex" flexDirection="column" gap={2}>
                  <Box display="flex" justifyContent="space-between">
                    <Typography color="text.secondary">ID Người dùng</Typography>
                    <Typography fontWeight="500">{userData?.id || "-"}</Typography>
                  </Box>
                  <Box display="flex" justifyContent="space-between">
                    <Typography color="text.secondary">Ngày tạo</Typography>
                    <Typography fontWeight="500">
                      {userData?.createdAt ? new Date(userData.createdAt).toLocaleDateString() : "-"}
                    </Typography>
                  </Box>
                  <Box display="flex" justifyContent="space-between">
                    <Typography color="text.secondary">Lần đăng nhập cuối</Typography>
                    <Typography fontWeight="500">
                      {userData?.lastLoginAt ? new Date(userData.lastLoginAt).toLocaleString() : "-"}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>

            {/* Card 6: Account Security */}
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent sx={{ p: 3 }}>
                <Typography variant="h6" fontWeight={600} mb={3}>
                  Bảo mật
                </Typography>
                <TextField
                  label="Mật khẩu"
                  value="********"
                  disabled
                  fullWidth
                  InputProps={{ sx: { borderRadius: "8px" } }}
                />
              </CardContent>
            </Card>
          </Box>
        </Grid2>
      </Grid2>
    </Edit>
  );
};

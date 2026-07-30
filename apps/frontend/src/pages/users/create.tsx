import React, { useState } from "react";
import { Create } from "@refinedev/mui";
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
  Button,
  Stack,
  IconButton,
  InputAdornment,
  Alert,
  LinearProgress,
} from "@mui/material";
import { HttpError, useNavigation } from "@refinedev/core";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";

export interface ICreateUserRequest {
  name: string;
  email: string;
  password?: string;
  confirmPassword?: string;
  role: string;
  status: string;
}

export const UserCreate: React.FC = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { list } = useNavigation();

  const {
    saveButtonProps,
    refineCore: { formLoading, onFinish },
    register,
    control,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<any, HttpError, ICreateUserRequest>({
    refineCoreProps: {
      action: "create",
      resource: "users",
      redirect: false,
      onMutationSuccess: () => {
        list("users");
      },
    },
    warnWhenUnsavedChanges: true,
  });

  const passwordValue = watch("password");
  const roleValue = watch("role");
  const statusValue = watch("status");

  const getPasswordStrength = (pwd?: string) => {
    let strength = 0;
    if (!pwd) return strength;
    if (pwd.length >= 8) strength += 25;
    if (/[A-Z]/.test(pwd)) strength += 25;
    if (/[a-z]/.test(pwd)) strength += 25;
    if (/[0-9]/.test(pwd) || /[^A-Za-z0-9]/.test(pwd)) strength += 25;
    return strength;
  };

  const getStrengthColor = (strength: number) => {
    if (strength === 0) return "inherit";
    if (strength <= 50) return "error";
    if (strength <= 75) return "warning";
    return "success";
  };

  const passwordStrength = getPasswordStrength(passwordValue);

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Thêm người dùng mới
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
              Hủy bỏ
            </Button>
            <Button
              variant="contained"
              disabled={formLoading}
              onClick={handleSubmit((data) => {
                onFinish(data);
              })}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: '#ff8c42', boxShadow: 'none', '&:hover': { bgcolor: '#e67a33', boxShadow: 'none' } }}
            >
              Tạo người dùng
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
                      Thông tin người dùng
                    </Typography>
                    <TextField
                      {...register("name", { required: "Họ và tên là bắt buộc", maxLength: { value: 100, message: "Tối đa 100 ký tự" } })}
                      error={!!errors.name}
                      helperText={errors.name?.message as string}
                      label="Họ và tên (*)"
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      autoFocus
                      sx={{ mb: 3 }}
                      InputProps={{ sx: { borderRadius: "8px" } }}
                    />
                    <TextField
                      {...register("email", {
                        required: "Email là bắt buộc",
                        pattern: {
                          value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                          message: "Email không hợp lệ",
                        },
                      })}
                      error={!!errors.email}
                      helperText={errors.email?.message as string}
                      label="Email (*)"
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      InputProps={{ sx: { borderRadius: "8px" } }}
                    />
                  </CardContent>
                </Card>

                {/* Authentication Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={3} color="text.primary">
                      Thông tin đăng nhập
                    </Typography>
                    
                    <TextField
                      {...register("password", {
                        required: "Mật khẩu là bắt buộc",
                        minLength: { value: 8, message: "Tối thiểu 8 ký tự" },
                      })}
                      error={!!errors.password}
                      helperText={errors.password?.message as string}
                      label="Mật khẩu (*)"
                      type={showPassword ? "text" : "password"}
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        sx: { borderRadius: "8px" },
                        endAdornment: (
                          <InputAdornment position="end">
                            <IconButton
                              onClick={() => setShowPassword(!showPassword)}
                              edge="end"
                            >
                              {showPassword ? <VisibilityOff /> : <Visibility />}
                            </IconButton>
                          </InputAdornment>
                        )
                      }}
                    />
                    
                    {passwordValue && (
                       <Box sx={{ mt: 1, mb: 3 }}>
                        <LinearProgress 
                          variant="determinate" 
                          value={passwordStrength} 
                          color={getStrengthColor(passwordStrength)} 
                          sx={{ height: 6, borderRadius: 3, mb: 1 }}
                        />
                        <Typography variant="caption" color="text.secondary">
                          Mức độ mạnh: {passwordStrength <= 50 ? "Yếu" : passwordStrength <= 75 ? "Trung bình" : "Mạnh"}
                        </Typography>
                      </Box>
                    )}

                    <TextField
                      {...register("confirmPassword", {
                        required: "Xác nhận mật khẩu là bắt buộc",
                        validate: (value) => value === passwordValue || "Mật khẩu không khớp",
                      })}
                      error={!!errors.confirmPassword}
                      helperText={errors.confirmPassword?.message as string}
                      label="Xác nhận mật khẩu (*)"
                      type={showConfirmPassword ? "text" : "password"}
                      margin="normal"
                      variant="outlined"
                      fullWidth
                      sx={passwordValue ? { mb: 0 } : { mb: 3, mt: 3 }}
                      InputProps={{
                        sx: { borderRadius: "8px" },
                        endAdornment: (
                          <InputAdornment position="end">
                            <IconButton
                              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                              edge="end"
                            >
                              {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                            </IconButton>
                          </InputAdornment>
                        )
                      }}
                    />
                  </CardContent>
                </Card>

                {/* Permission Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Phân quyền
                    </Typography>
                    
                    <Controller
                      control={control}
                      name="role"
                      defaultValue="CUSTOMER"
                      rules={{ required: "Vai trò là bắt buộc" }}
                      render={({ field }) => (
                        <FormControl component="fieldset" sx={{ mb: 2, width: "100%" }}>
                          <RadioGroup {...field} row>
                            <FormControlLabel value="ADMIN" control={<Radio color="primary" />} label="ADMIN" sx={{ color: '#1976d2' }} />
                            <FormControlLabel value="STAFF" control={<Radio color="primary" />} label="STAFF" sx={{ color: '#9c27b0' }} />
                            <FormControlLabel value="CUSTOMER" control={<Radio color="primary" />} label="CUSTOMER" sx={{ color: 'text.secondary' }} />
                          </RadioGroup>
                        </FormControl>
                      )}
                    />
                    {errors.role && (
                      <Typography variant="caption" color="error" display="block">
                        {errors.role.message as string}
                      </Typography>
                    )}

                    {roleValue === "ADMIN" && (
                      <Alert severity="warning" sx={{ borderRadius: '8px' }}>
                        Người dùng này sẽ có toàn quyền quản trị hệ thống.
                      </Alert>
                    )}
                  </CardContent>
                </Card>
              </Stack>
            </Grid2>

            {/* Side Column */}
            <Grid2 size={{ xs: 12, md: 4 }}>
              <Stack spacing={3}>
                {/* Status Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Trạng thái tài khoản
                    </Typography>
                    <Controller
                      control={control}
                      name="status"
                      defaultValue="ACTIVE"
                      render={({ field }) => (
                        <FormControl component="fieldset" sx={{ mb: 2, width: "100%" }}>
                          <RadioGroup {...field} row>
                            <FormControlLabel value="ACTIVE" control={<Radio color="success" />} label="ACTIVE" sx={{ color: 'success.main' }} />
                            <FormControlLabel value="INACTIVE" control={<Radio color="error" />} label="INACTIVE" sx={{ color: 'error.main' }} />
                          </RadioGroup>
                        </FormControl>
                      )}
                    />

                    {statusValue === "INACTIVE" && (
                       <Alert severity="error" sx={{ borderRadius: '8px' }}>
                         Người dùng sẽ chưa thể đăng nhập cho đến khi được kích hoạt.
                       </Alert>
                    )}
                  </CardContent>
                </Card>

                {/* System Info Card */}
                <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
                  <CardContent sx={{ p: 3 }}>
                    <Typography variant="h6" fontWeight="600" mb={2} color="text.primary">
                      Thông tin hệ thống
                    </Typography>
                    <Stack spacing={2}>
                      <TextField
                        disabled
                        label="User ID"
                        defaultValue="Tự động tạo"
                        variant="filled"
                        size="small"
                        InputProps={{ disableUnderline: true, sx: { borderRadius: "6px" } }}
                      />
                      <TextField
                        disabled
                        label="Created At"
                        defaultValue="Tự động tạo"
                        variant="filled"
                        size="small"
                        InputProps={{ disableUnderline: true, sx: { borderRadius: "6px" } }}
                      />
                      <TextField
                        disabled
                        label="Updated At"
                        defaultValue="Tự động tạo"
                        variant="filled"
                        size="small"
                        InputProps={{ disableUnderline: true, sx: { borderRadius: "6px" } }}
                      />
                       <TextField
                        disabled
                        label="Last Login"
                        defaultValue="Tự động tạo"
                        variant="filled"
                        size="small"
                        InputProps={{ disableUnderline: true, sx: { borderRadius: "6px" } }}
                      />
                    </Stack>
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

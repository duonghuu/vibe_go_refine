import React, { useState, useMemo } from "react";
import {
  useDataGrid,
  List,
  EditButton,
  ShowButton,
  CreateButton,
} from "@refinedev/mui";
import {
  DataGrid,
  GridColDef,
  GridActionsCellItem,
} from "@mui/x-data-grid";
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid2,
  Chip,
  IconButton,
  TextField,
  InputAdornment,
  MenuItem,
  Drawer,
  Button,
  Stack,
  FormControl,
  InputLabel,
  Select,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import LockIcon from "@mui/icons-material/Lock";
import LockOpenIcon from "@mui/icons-material/LockOpen";
import KeyIcon from "@mui/icons-material/Key";
import { useForm } from "@refinedev/react-hook-form";
import { Controller } from "react-hook-form";
import { useCustomMutation, useUpdate, HttpError, useGetIdentity } from "@refinedev/core";

export interface IUser {
  id: number;
  email: string;
  name: string;
  role: "ADMIN" | "STAFF" | "CUSTOMER";
  status: "ACTIVE" | "INACTIVE";
  lastLoginAt: string;
  createdAt: string;
}

export const UserList = () => {
  const { dataGridProps, setFilters, filters } = useDataGrid<IUser>({
    resource: "users",
    syncWithLocation: true,
    pagination: {
      mode: "server",
    },
    sorters: {
      initial: [
        {
          field: "id",
          order: "desc",
        },
      ],
    },
  });

  const { data: identity } = useGetIdentity<{ id: number; role: string }>();

  // Filter state
  const [searchTerm, setSearchTerm] = useState("");
  const [roleFilter, setRoleFilter] = useState<string>("ALL");
  const [statusFilter, setStatusFilter] = useState<string>("ALL");

  // Drawer state for Create / Edit
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  // Dialogs
  const [statusDialog, setStatusDialog] = useState<{ open: boolean; user: IUser | null }>({ open: false, user: null });
  const [passwordDialog, setPasswordDialog] = useState<{ open: boolean; userId: number | null }>({ open: false, userId: null });
  const [newPassword, setNewPassword] = useState("");

  const { mutate: updateStatus } = useUpdate<IUser>();
  const { mutate: customMutate } = useCustomMutation<any>();

  // Form hook for drawer
  const {
    refineCore: { onFinish, formLoading },
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<IUser, HttpError, any>({
    refineCoreProps: {
      resource: "users",
      action: editingId ? "edit" : "create",
      id: editingId || undefined,
      onMutationSuccess: () => {
        setDrawerOpen(false);
        reset();
      },
    },
  });

  // Handle Search with debounce
  React.useEffect(() => {
    const handler = setTimeout(() => {
      setFilters([
        {
          field: "search", // Can map this to name/email in data provider or backend
          operator: "contains",
          value: searchTerm,
        },
        {
          field: "role",
          operator: "eq",
          value: roleFilter === "ALL" ? undefined : roleFilter,
        },
        {
          field: "status",
          operator: "eq",
          value: statusFilter === "ALL" ? undefined : statusFilter,
        },
      ]);
    }, 300);

    return () => {
      clearTimeout(handler);
    };
  }, [searchTerm, roleFilter, statusFilter, setFilters]);

  // Open Drawer for Edit
  const handleEdit = (user: IUser) => {
    setEditingId(user.id);
    reset({
      email: user.email,
      name: user.name,
      role: user.role,
      status: user.status,
    });
    setDrawerOpen(true);
  };

  const handleToggleStatus = () => {
    if (statusDialog.user) {
      updateStatus({
        resource: "users",
        id: statusDialog.user.id,
        values: {
          status: statusDialog.user.status === "ACTIVE" ? "INACTIVE" : "ACTIVE",
        },
      }, {
        onSuccess: () => setStatusDialog({ open: false, user: null }),
      });
    }
  };

  const handleResetPassword = () => {
    if (passwordDialog.userId && newPassword.length >= 8) {
      customMutate({
        url: `/users/${passwordDialog.userId}/reset-password`,
        method: "post",
        values: { password: newPassword },
      }, {
        onSuccess: () => {
          setPasswordDialog({ open: false, userId: null });
          setNewPassword("");
        }
      });
    }
  };

  // Metrics
  const totalUsers = dataGridProps.rows.length; // Actually total from pagination but we can use rows for mock. Real total: dataGridProps.rowCount
  const totalAdmin = dataGridProps.rows.filter(r => r.role === "ADMIN").length;
  const totalStaff = dataGridProps.rows.filter(r => r.role === "STAFF").length;
  const totalActive = dataGridProps.rows.filter(r => r.status === "ACTIVE").length;

  const columns = useMemo<GridColDef<IUser>[]>(
    () => [
      { field: "id", headerName: "ID", width: 70 },
      { field: "email", headerName: "Email", minWidth: 200, flex: 1 },
      { field: "name", headerName: "Tên", minWidth: 150, flex: 1 },
      {
        field: "role",
        headerName: "Quyền",
        width: 120,
        renderCell: ({ value }) => (
          <Chip
            label={value}
            size="small"
            color={
              value === "ADMIN" ? "info" :
              value === "STAFF" ? "secondary" : "default"
            }
            sx={{ fontWeight: "bold" }}
          />
        ),
      },
      {
        field: "status",
        headerName: "Trạng thái",
        width: 120,
        renderCell: ({ value }) => (
          <Chip
            label={value}
            size="small"
            color={value === "ACTIVE" ? "success" : "error"}
            variant="outlined"
            sx={{ fontWeight: "bold" }}
          />
        ),
      },
      {
        field: "lastLoginAt",
        headerName: "Last Login",
        width: 160,
        renderCell: ({ value }) => value ? new Date(value).toLocaleString() : "-",
      },
      {
        field: "createdAt",
        headerName: "Ngày tạo",
        width: 120,
        renderCell: ({ value }) => value ? new Date(value).toLocaleDateString() : "-",
      },
      {
        field: "actions",
        headerName: "Thao tác",
        type: "actions",
        width: 180,
        getActions: ({ row }) => {
          const isCurrentUser = identity?.id === row.id;
          const canEdit = identity?.role === "ADMIN" || identity?.role === "STAFF";

          return [
            <GridActionsCellItem
              key="edit"
              icon={<EditButton hideText size="small" />}
              label="Edit"
              onClick={() => handleEdit(row)}
              disabled={!canEdit}
            />,
            <GridActionsCellItem
              key="password"
              icon={<KeyIcon color="warning" />}
              label="Reset Password"
              onClick={() => setPasswordDialog({ open: true, userId: row.id })}
              disabled={identity?.role !== "ADMIN"}
            />,
            <GridActionsCellItem
              key="status"
              icon={row.status === "ACTIVE" ? <LockIcon color="error" /> : <LockOpenIcon color="success" />}
              label="Toggle Status"
              onClick={() => setStatusDialog({ open: true, user: row })}
              disabled={isCurrentUser}
            />,
          ];
        },
      },
    ],
    [identity]
  );

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" fontWeight="700" color="#202224" letterSpacing="-0.02em">
          Quản lý người dùng
        </Typography>
      </Box>

      {/* Summary Cards */}
      <Grid2 container spacing={3} mb={4}>
        {[
          { label: "Tổng số User", value: dataGridProps.rowCount || 0, color: "primary.main" },
          { label: "Tổng Admin", value: totalAdmin, color: "info.main" },
          { label: "Tổng Staff", value: totalStaff, color: "secondary.main" },
          { label: "Đang hoạt động", value: totalActive, color: "success.main" },
        ].map((card, idx) => (
          <Grid2 size={{ xs: 12, sm: 6, md: 3 }} key={idx}>
            <Card sx={{ borderRadius: "14px", border: "1px solid #D5D5D5", boxShadow: "none" }}>
              <CardContent>
                <Typography variant="body2" color="text.secondary" fontWeight={600} gutterBottom>
                  {card.label}
                </Typography>
                <Typography variant="h5" fontWeight={700} color={card.color}>
                  {card.value}
                </Typography>
              </CardContent>
            </Card>
          </Grid2>
        ))}
      </Grid2>

      <List 
        wrapperProps={{
          sx: {
            bgcolor: "#fff",
            borderRadius: "14px",
            border: "1px solid #D5D5D5",
            boxShadow: "none",
            p: 0,
            overflow: "hidden",
          },
        }}
        headerProps={{
          sx: {
            p: 2,
            px: 3,
            borderBottom: "1px solid #D5D5D5",
          },
        }}
        title=""
        headerButtons={(props) => (
          <Stack direction="row" spacing={2}>
            <TextField
              placeholder="Tìm theo tên, email..."
              variant="outlined"
              size="small"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon sx={{ color: 'text.secondary' }} />
                  </InputAdornment>
                ),
                sx: { borderRadius: '50px', backgroundColor: '#F5F6FA', width: { xs: '100%', md: '300px' } }
              }}
            />
            
            <Select
              size="small"
              displayEmpty
              value={roleFilter}
              onChange={(e) => setRoleFilter(e.target.value)}
              sx={{ borderRadius: '8px', minWidth: '150px', backgroundColor: '#F5F6FA' }}
            >
              <MenuItem value="ALL">Tất cả vai trò</MenuItem>
              <MenuItem value="ADMIN">ADMIN</MenuItem>
              <MenuItem value="STAFF">STAFF</MenuItem>
              <MenuItem value="CUSTOMER">CUSTOMER</MenuItem>
            </Select>

            <Select
              size="small"
              displayEmpty
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              sx={{ borderRadius: '8px', minWidth: '150px', backgroundColor: '#F5F6FA' }}
            >
              <MenuItem value="ALL">Tất cả trạng thái</MenuItem>
              <MenuItem value="ACTIVE">Hoạt động</MenuItem>
              <MenuItem value="INACTIVE">Đã khóa</MenuItem>
            </Select>
            <CreateButton />
          </Stack>
        )}
      >
        <DataGrid
          {...dataGridProps}
          columns={columns}
          rowHeight={80}
          density="standard"
          disableColumnMenu
          checkboxSelection
          sx={{
            border: 'none',
            '& .MuiDataGrid-columnHeaders': {
              backgroundColor: '#F5F6FA',
              borderBottom: '1px solid #D5D5D5',
              color: 'text.primary',
              fontWeight: 'bold',
            },
            '& .MuiDataGrid-cell': {
              borderBottom: '1px solid #f3f4f6',
            },
            '& .MuiDataGrid-row:hover': {
              backgroundColor: '#f9fafb',
            }
          }}
        />
      </List>

      {/* Drawer for Create / Edit */}
      <Drawer
        anchor="right"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        PaperProps={{ sx: { width: { xs: '100%', sm: 400 }, p: 3 } }}
      >
        <Typography variant="h6" fontWeight="700" mb={3}>
          Chỉnh sửa người dùng
        </Typography>

        <Box component="form" onSubmit={handleSubmit(onFinish)} display="flex" flexDirection="column" gap={3}>
          <TextField
            {...register("name", { required: "Tên là bắt buộc" })}
            label="Họ và tên (*)"
            error={!!errors.name}
            helperText={errors.name?.message as string}
            fullWidth
            InputProps={{ sx: { borderRadius: "8px" } }}
          />

          <TextField
            {...register("email", { 
              required: "Email là bắt buộc",
              pattern: { value: /^\S+@\S+$/i, message: "Email không hợp lệ" }
            })}
            label="Email (*)"
            error={!!errors.email}
            helperText={errors.email?.message as string}
            fullWidth
            disabled={true}
            InputProps={{ sx: { borderRadius: "8px" } }}
          />

          <Controller
            control={control}
            name="role"
            render={({ field }) => (
              <FormControl fullWidth>
                <InputLabel>Quyền</InputLabel>
                <Select {...field} label="Quyền" sx={{ borderRadius: "8px" }}>
                  <MenuItem value="ADMIN">ADMIN</MenuItem>
                  <MenuItem value="STAFF">STAFF</MenuItem>
                  <MenuItem value="CUSTOMER">CUSTOMER</MenuItem>
                </Select>
              </FormControl>
            )}
          />

          <Controller
            control={control}
            name="status"
            render={({ field }) => (
              <FormControl fullWidth>
                <InputLabel>Trạng thái</InputLabel>
                <Select {...field} label="Trạng thái" sx={{ borderRadius: "8px" }}>
                  <MenuItem value="ACTIVE">Hoạt động</MenuItem>
                  <MenuItem value="INACTIVE">Khóa</MenuItem>
                </Select>
              </FormControl>
            )}
          />

          <Stack direction="row" spacing={2} justifyContent="flex-end" mt={2}>
            <Button
              variant="outlined"
              color="secondary"
              onClick={() => setDrawerOpen(false)}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", borderColor: '#D5D5D5' }}
            >
              Hủy
            </Button>
            <Button
              type="submit"
              variant="contained"
              disabled={formLoading}
              sx={{ borderRadius: "8px", textTransform: "none", fontWeight: "600", bgcolor: 'primary.main', boxShadow: 'none' }}
            >
              Lưu
            </Button>
          </Stack>
        </Box>
      </Drawer>

      {/* Status Confirm Dialog */}
      <Dialog open={statusDialog.open} onClose={() => setStatusDialog({ open: false, user: null })}>
        <DialogTitle>Xác nhận thay đổi</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Bạn có chắc chắn muốn {statusDialog.user?.status === "ACTIVE" ? "khóa" : "kích hoạt"} tài khoản <b>{statusDialog.user?.email}</b>?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setStatusDialog({ open: false, user: null })} color="secondary" sx={{ textTransform: "none", fontWeight: "600" }}>Hủy</Button>
          <Button onClick={handleToggleStatus} color={statusDialog.user?.status === "ACTIVE" ? "error" : "success"} variant="contained" sx={{ textTransform: "none", fontWeight: "600", borderRadius: "8px", boxShadow: "none" }}>
            Xác nhận
          </Button>
        </DialogActions>
      </Dialog>

      {/* Reset Password Dialog */}
      <Dialog open={passwordDialog.open} onClose={() => setPasswordDialog({ open: false, userId: null })}>
        <DialogTitle>Đặt lại mật khẩu</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            Nhập mật khẩu mới cho người dùng. Mật khẩu phải có ít nhất 8 ký tự.
          </DialogContentText>
          <TextField
            fullWidth
            type="password"
            label="Mật khẩu mới"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            InputProps={{ sx: { borderRadius: "8px" } }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPasswordDialog({ open: false, userId: null })} color="secondary" sx={{ textTransform: "none", fontWeight: "600" }}>Hủy</Button>
          <Button 
            onClick={handleResetPassword} 
            color="success" 
            variant="contained" 
            disabled={newPassword.length < 8}
            sx={{ textTransform: "none", fontWeight: "600", borderRadius: "8px", boxShadow: "none" }}
          >
            Lưu mật khẩu
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

import { PropsWithChildren } from "react";
import { useGetIdentity } from "@refinedev/core";
import { Alert, CircularProgress, Stack } from "@mui/material";

interface Identity { role?: string; }
export const RoleGuard = ({ roles, children }: PropsWithChildren<{ roles: string[] }>) => {
  const { data, isLoading } = useGetIdentity<Identity>();
  if (isLoading) return <Stack alignItems="center" py={8}><CircularProgress /></Stack>;
  if (!data?.role || !roles.includes(data.role)) return <Alert severity="error">Bạn không có quyền truy cập chức năng này.</Alert>;
  return <>{children}</>;
};


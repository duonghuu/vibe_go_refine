import React from "react";
import { Button, Stack } from "@mui/material";

interface PageEditActionsProps {
  disabled: boolean;
  onCancel: () => void;
  onSubmit: () => void;
}

export const PageEditActions: React.FC<PageEditActionsProps> = ({ disabled, onCancel, onSubmit }) => (
  <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2, width: "100%" }}>
    <Button variant="outlined" color="secondary" disabled={disabled} onClick={onCancel} sx={{ textTransform: "none", fontWeight: 700 }}>
      Hủy
    </Button>
    <Button variant="contained" color="success" disabled={disabled} onClick={onSubmit} sx={{ textTransform: "none", fontWeight: 700, boxShadow: "none" }}>
      Lưu thay đổi
    </Button>
  </Stack>
);

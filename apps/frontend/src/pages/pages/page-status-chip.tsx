import React from "react";
import { Chip } from "@mui/material";
import { PageStatus } from "./page-types";

interface PageStatusChipProps {
  status: PageStatus;
}

export const PageStatusChip: React.FC<PageStatusChipProps> = ({ status }) => (
  <Chip
    label={status === "PUBLISHED" ? "Đã xuất bản" : "Bản nháp"}
    size="small"
    color={status === "PUBLISHED" ? "success" : "default"}
    variant={status === "PUBLISHED" ? "filled" : "outlined"}
    sx={{ fontWeight: 700 }}
  />
);

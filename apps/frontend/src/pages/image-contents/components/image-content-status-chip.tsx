import { Chip } from "@mui/material";
import { ImageContentStatus } from "../image-content-types";

export const ImageContentStatusChip = ({ status }: { status: ImageContentStatus }) => (
  <Chip size="small" label={status === "ACTIVE" ? "Đang hoạt động" : "Tạm ẩn"} color={status === "ACTIVE" ? "success" : "default"} />
);


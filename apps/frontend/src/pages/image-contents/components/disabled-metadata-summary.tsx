import { Alert, Stack, Typography } from "@mui/material";
import { IImageContent, IMAGE_CONTENT_FIELD_KEYS, fieldLabel } from "../image-content-types";

export const DisabledMetadataSummary = ({ item }: { item: IImageContent }) => {
  const values = IMAGE_CONTENT_FIELD_KEYS.map((key) => ({ key, value: item[key === "url" ? "url" : key] })).filter(({ value }) => value);
  if (values.length === 0) return null;
  return <Alert severity="info"><Typography variant="subtitle2" mb={0.5}>Dữ liệu đang tạm ẩn</Typography><Stack spacing={0.25}>{values.map(({ key, value }) => <Typography key={key} variant="body2"><strong>{fieldLabel[key]}:</strong> {value}</Typography>)}</Stack></Alert>;
};


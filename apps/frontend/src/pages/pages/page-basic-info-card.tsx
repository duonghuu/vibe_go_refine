import React from "react";
import { Card, CardContent, TextField, Typography } from "@mui/material";
import { FieldErrors, UseFormRegister, UseFormWatch } from "react-hook-form";
import { IPageCreatePayload } from "./page-form-types";

interface PageBasicInfoCardProps {
  register: UseFormRegister<IPageCreatePayload>;
  errors: FieldErrors<IPageCreatePayload>;
  watch: UseFormWatch<IPageCreatePayload>;
  onSlugChange: (value: string) => void;
}

export const PageBasicInfoCard: React.FC<PageBasicInfoCardProps> = ({
  register,
  errors,
  watch,
  onSlugChange,
}) => (
  <Card sx={{ border: "1px solid", borderColor: "divider", borderRadius: 2, boxShadow: "none" }}>
    <CardContent sx={{ p: { xs: 2, md: 3 } }}>
      <Typography variant="h6" fontWeight={700} mb={3}>
        Thông tin cơ bản
      </Typography>
      <TextField
        {...register("title", {
          required: "Tiêu đề là bắt buộc",
          validate: (value) => value.trim().length > 0 || "Tiêu đề không được để trống",
          maxLength: { value: 255, message: "Tiêu đề tối đa 255 ký tự" },
        })}
        label="Tiêu đề (*)"
        fullWidth
        error={Boolean(errors.title)}
        helperText={errors.title?.message}
        margin="normal"
        InputProps={{ sx: { borderRadius: 1 } }}
      />
      <TextField
        {...register("slug", {
          required: "Slug là bắt buộc",
          validate: (value) => value.trim().length > 0 || "Slug không được để trống",
          maxLength: { value: 255, message: "Slug tối đa 255 ký tự" },
          pattern: {
            value: /^[a-z0-9]+(?:-[a-z0-9]+)*$/,
            message: "Slug chỉ gồm chữ thường, số và dấu gạch ngang",
          },
          onChange: (event) => onSlugChange(event.target.value),
        })}
        label="Đường dẫn (Slug) (*)"
        fullWidth
        error={Boolean(errors.slug)}
        helperText={errors.slug?.message ?? "Đường dẫn public của trang"}
        margin="normal"
        InputLabelProps={{ shrink: !!watch("slug") || undefined }}
        value={watch("slug")}
        InputProps={{ sx: { borderRadius: 1 } }}
      />
      <TextField
        {...register("content", {
          required: "Nội dung là bắt buộc",
          validate: (value) => value.replace(/<[^>]*>/g, "").trim().length > 0 || "Nội dung không được để trống",
        })}
        label="Nội dung (*)"
        fullWidth
        multiline
        minRows={16}
        error={Boolean(errors.content)}
        helperText={errors.content?.message ?? "Soạn nội dung trang tại đây"}
        margin="normal"
        InputProps={{ sx: { borderRadius: 1, alignItems: "flex-start" } }}
      />
    </CardContent>
  </Card>
);

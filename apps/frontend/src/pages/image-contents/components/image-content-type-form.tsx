import { Controller, UseFormReturn } from "react-hook-form";
import { Card, CardContent, Checkbox, FormControlLabel, Grid2, MenuItem, Stack, Switch, TextField, Typography } from "@mui/material";
import { IImageContentTypeFormValues, IMAGE_CONTENT_FIELD_KEYS, fieldLabel } from "../image-content-types";

interface Props { form: UseFormReturn<IImageContentTypeFormValues>; isEdit?: boolean; }

export const ImageContentTypeForm = ({ form, isEdit = false }: Props) => {
  const { register, control, watch, setValue, formState: { errors } } = form;
  const limitEnabled = watch("hasItemLimit");
  return <Grid2 container spacing={3}>
    <Grid2 size={{ xs: 12, md: 6 }}>
      <Card variant="outlined"><CardContent><Typography variant="h6" mb={2}>Thông tin chung</Typography><Stack spacing={2}>
        <TextField label="Mã loại" disabled={isEdit} fullWidth {...register("code", { required: "Mã loại là bắt buộc", pattern: { value: /^[A-Z][A-Z0-9_]{0,49}$/, message: "Dùng chữ in hoa, số và dấu gạch dưới" } })} error={Boolean(errors.code)} helperText={errors.code?.message} onChange={(event) => { if (!isEdit) setValue("code", event.target.value.toUpperCase(), { shouldDirty: true, shouldValidate: true }); }} />
        <TextField label="Tên loại" fullWidth {...register("name", { required: "Tên loại là bắt buộc", maxLength: { value: 100, message: "Tối đa 100 ký tự" } })} error={Boolean(errors.name)} helperText={errors.name?.message} />
        <TextField label="Thứ tự" type="number" fullWidth {...register("sortOrder", { valueAsNumber: true, min: { value: 0, message: "Phải lớn hơn hoặc bằng 0" } })} error={Boolean(errors.sortOrder)} helperText={errors.sortOrder?.message} />
      </Stack></CardContent></Card>
    </Grid2>
    <Grid2 size={{ xs: 12, md: 6 }}>
      <Card variant="outlined"><CardContent><Typography variant="h6" mb={2}>Trạng thái và giới hạn</Typography><Stack spacing={2}>
        <TextField select label="Trạng thái" fullWidth {...register("status", { required: true })}><MenuItem value="ACTIVE">Đang hoạt động</MenuItem><MenuItem value="INACTIVE">Tạm ẩn</MenuItem></TextField>
        <FormControlLabel control={<Switch checked={limitEnabled} onChange={(event) => { setValue("hasItemLimit", event.target.checked, { shouldDirty: true }); if (!event.target.checked) setValue("maxItems", null, { shouldDirty: true }); }} />} label="Giới hạn số nội dung" />
        {limitEnabled && <TextField label="Số lượng tối đa" type="number" fullWidth {...register("maxItems", { valueAsNumber: true, required: "Nhập giới hạn", min: { value: 1, message: "Từ 1 đến 1000" }, max: { value: 1000, message: "Từ 1 đến 1000" } })} error={Boolean(errors.maxItems)} helperText={errors.maxItems?.message} />}
      </Stack></CardContent></Card>
    </Grid2>
    <Grid2 size={12}>
      <Card variant="outlined"><CardContent><Typography variant="h6" mb={1}>Cấu hình metadata</Typography><Typography variant="body2" color="text.secondary" mb={2}>Chỉ các field được bật mới xuất hiện trên form nội dung.</Typography>
        <Stack spacing={1}>{IMAGE_CONTENT_FIELD_KEYS.map((key) => <Stack key={key} direction={{ xs: "column", sm: "row" }} alignItems={{ sm: "center" }} justifyContent="space-between" sx={{ py: 1, borderBottom: 1, borderColor: "divider" }}>
          <Typography fontWeight={600}>{fieldLabel[key]}</Typography><Stack direction="row" spacing={2}><Controller control={control} name={`fieldConfig.${key}.enabled`} render={({ field }) => <FormControlLabel control={<Switch checked={field.value} onChange={(event) => { field.onChange(event.target.checked); if (!event.target.checked) setValue(`fieldConfig.${key}.required`, false, { shouldDirty: true }); }} />} label="Hiển thị" />} /><Controller control={control} name={`fieldConfig.${key}.required`} render={({ field }) => <FormControlLabel control={<Checkbox checked={field.value} disabled={!watch(`fieldConfig.${key}.enabled`)} onChange={(_, value) => field.onChange(value)} />} label="Bắt buộc" />} /></Stack>
        </Stack>)}</Stack>
      </CardContent></Card>
    </Grid2>
  </Grid2>;
};


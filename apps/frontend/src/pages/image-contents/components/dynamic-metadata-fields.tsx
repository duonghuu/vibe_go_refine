import { Control, Controller, FieldErrors } from "react-hook-form";
import { Stack, TextField } from "@mui/material";
import {
  IImageContentFieldConfig,
  IImageContentFormValues,
  ImageContentFieldKey,
  fieldLabel,
} from "../image-content-types";

interface Props {
  control: Control<IImageContentFormValues>;
  errors: FieldErrors<IImageContentFormValues>;
  config: IImageContentFieldConfig;
}
const maxLength: Record<ImageContentFieldKey, number> = {
  name: 255,
  description: 5000,
  secondaryDescription: 5000,
  url: 500,
};

export const DynamicMetadataFields = ({ control, errors, config }: Props) => (
  <Stack spacing={2}>
    {(
      [
        "name",
        "description",
        "secondaryDescription",
        "url",
      ] as ImageContentFieldKey[]
    )
      .filter((key) => config[key].enabled)
      .map((key) => (
        <Controller
          key={key}
          name={key}
          control={control}
          rules={{
            required: config[key].required
              ? `${fieldLabel[key]} là bắt buộc`
              : false,
            maxLength: {
              value: maxLength[key],
              message: `Tối đa ${maxLength[key]} ký tự`,
            },
            validate:
              key === "url"
                ? (value) => {
                    if (!value) return true;
                    const valid =
                      /^(\/(?!\/)[^\s]*|https?:\/\/[^\s]+)$/i.test(value) &&
                      !/[\u0000-\u001f]/.test(value);
                    return (
                      valid ||
                      "Liên kết phải là đường dẫn nội bộ hoặc http/https hợp lệ"
                    );
                  }
                : undefined,
          }}
          render={({ field }) => (
            <TextField
              {...field}
              fullWidth
              label={`${fieldLabel[key]}${config[key].required ? " *" : ""}`}
              multiline={
                key === "description" || key === "secondaryDescription"
              }
              minRows={
                key === "description" || key === "secondaryDescription"
                  ? 4
                  : undefined
              }
              slotProps={{
                inputLabel: { shrink: true },
              }}
              inputProps={{ maxLength: maxLength[key] }}
              error={Boolean(errors[key])}
              helperText={
                errors[key]?.message?.toString() ??
                (key === "description" || key === "secondaryDescription"
                  ? `${(field.value ?? "").length}/${maxLength[key]}`
                  : undefined)
              }
            />
          )}
        />
      ))}
  </Stack>
);

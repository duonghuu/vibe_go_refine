import { useEffect, useMemo, useState } from "react";
import { Create, Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { HttpError, useList, useNavigation } from "@refinedev/core";
import {
  Box,
  Button,
  Card,
  CardContent,
  FormControlLabel,
  MenuItem,
  Stack,
  Switch,
  TextField,
  Typography,
} from "@mui/material";
import { ImageMediaPickerDialog } from "./components/image-media-picker-dialog";
import { DisabledMetadataSummary } from "./components/disabled-metadata-summary";
import { DynamicMetadataFields } from "./components/dynamic-metadata-fields";
import {
  EMPTY_FIELD_CONFIG,
  ICreateImageContentPayload,
  IImageContent,
  IImageContentFormValues,
  IImageContentMedia,
  IImageContentType,
  IUpdateImageContentPayload,
  mediaPreviewUrl,
} from "./image-content-types";

const empty: IImageContentFormValues = {
  type: null,
  media: null,
  name: "",
  description: "",
  secondaryDescription: "",
  url: "",
  status: "ACTIVE",
};
const compact = (value: string): string | null => value.trim() || null;
const createPayload = (
  values: IImageContentFormValues,
): ICreateImageContentPayload => ({
  typeCode: values.type?.code ?? "",
  mediaId: values.media?.id ?? 0,
  name: compact(values.name),
  description: compact(values.description),
  secondaryDescription: compact(values.secondaryDescription),
  url: compact(values.url),
  status: values.status,
});
const ActionBar = ({
  loading,
  onCancel,
  onSave,
}: {
  loading: boolean;
  onCancel: () => void;
  onSave: () => void;
}) => (
  <Stack direction="row" justifyContent="flex-end" spacing={1.5} mt={3}>
    <Button variant="outlined" onClick={onCancel}>
      Hủy
    </Button>
    <Button
      variant="contained"
      color="success"
      disabled={loading}
      onClick={onSave}
    >
      Lưu nội dung
    </Button>
  </Stack>
);

const ImageContentForm = ({
  form,
  types,
  isEdit,
  item,
  onSubmit,
}: {
  form: ReturnType<
    typeof useForm<IImageContent, HttpError, IImageContentFormValues>
  >;
  types: IImageContentType[];
  isEdit?: boolean;
  item?: IImageContent;
  onSubmit: (values: IImageContentFormValues) => void;
}) => {
  const {
    control,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = form;
  const [pickerOpen, setPickerOpen] = useState(false);
  const selectedType = watch("type");
  const selectedMedia = watch("media");
  useEffect(() => {
    if (!isEdit && !selectedType && types.length === 1)
      setValue("type", types[0], { shouldDirty: false });
  }, [isEdit, selectedType, setValue, types]);
  return (
    <>
      <Card variant="outlined">
        <CardContent>
          <Typography variant="h6" mb={2}>
            Thông tin nội dung
          </Typography>
          <Stack spacing={2}>
            <TextField
              select
              fullWidth
              label="Type"
              disabled={isEdit}
              value={selectedType?.code ?? ""}
              onChange={(event) =>
                setValue(
                  "type",
                  types.find((type) => type.code === event.target.value) ??
                    null,
                  { shouldDirty: true },
                )
              }
            >
              <MenuItem value="">Chọn Type</MenuItem>
              {types.map((type) => (
                <MenuItem key={type.code} value={type.code}>
                  {type.name} ({type.code})
                </MenuItem>
              ))}
            </TextField>
            {selectedType && (
              <Typography variant="body2" color="text.secondary">
                Type{" "}
                {selectedType.status === "ACTIVE" ? "đang hoạt động" : "tạm ẩn"}{" "}
                · {selectedType.itemCount} / {selectedType.maxItems ?? "∞"}
              </Typography>
            )}
            <Box>
              <Button variant="outlined" onClick={() => setPickerOpen(true)}>
                {selectedMedia ? "Đổi ảnh" : "Chọn ảnh"}
              </Button>
              {selectedMedia && (
                <Stack direction="row" alignItems="center" spacing={1.5} mt={1}>
                  <Box
                    component="img"
                    src={mediaPreviewUrl(selectedMedia)}
                    alt={selectedMedia.fileName}
                    sx={{
                      width: 96,
                      height: 64,
                      objectFit: "cover",
                      borderRadius: 1,
                    }}
                  />
                  <Typography variant="body2">
                    {selectedMedia.fileName}
                  </Typography>
                </Stack>
              )}
            </Box>
            <DynamicMetadataFields
              control={control}
              errors={errors}
              config={
                selectedType?.fieldConfig ??
                item?.fieldConfig ??
                EMPTY_FIELD_CONFIG
              }
            />
            <Stack
              direction="row"
              alignItems="center"
              justifyContent="space-between"
            >
              <Typography>Status</Typography>
              <FormControlLabel
                control={
                  <Switch
                    checked={watch("status") === "ACTIVE"}
                    onChange={(event) =>
                      setValue(
                        "status",
                        event.target.checked ? "ACTIVE" : "INACTIVE",
                        { shouldDirty: true },
                      )
                    }
                  />
                }
                label={
                  watch("status") === "ACTIVE" ? "Đang hoạt động" : "Tạm ẩn"
                }
              />
            </Stack>
          </Stack>
        </CardContent>
      </Card>
      <ImageMediaPickerDialog
        open={pickerOpen}
        selected={selectedMedia}
        onClose={() => setPickerOpen(false)}
        onSelect={(media: IImageContentMedia) =>
          setValue("media", media, { shouldDirty: true, shouldValidate: true })
        }
      />
      <input
        type="hidden"
        {...form.register("media", {
          validate: (value) => (value?.id ? true : "Hãy chọn một hình ảnh"),
        })}
      />
      {item && (
        <Box mt={2}>
          <DisabledMetadataSummary item={item} />
        </Box>
      )}
      <ActionBar
        loading={form.formState.isSubmitting || form.refineCore.formLoading}
        onCancel={() => window.history.back()}
        onSave={() => void handleSubmit(onSubmit)()}
      />
    </>
  );
};

export const ImageContentCreate = () => {
  const { list } = useNavigation();
  const params = new URLSearchParams(window.location.search);
  const typeCode = params.get("typeCode") ?? "";
  const types = useList<IImageContentType>({
    resource: "image-content-types",
    dataProviderName: "imageContent",
    pagination: { currentPage: 1, pageSize: 100 },
  });
  const form = useForm<IImageContent, HttpError, IImageContentFormValues>({
    refineCoreProps: {
      resource: "image-contents",
      action: "create",
      redirect: false,
      onMutationSuccess: () => list("image-contents"),
    },
    defaultValues: empty,
    warnWhenUnsavedChanges: true,
  });
  const onFinishPayload = form.refineCore.onFinish as unknown as (
    values: ICreateImageContentPayload,
  ) => Promise<unknown>;
  useEffect(() => {
    if (typeCode && types.result.data.length)
      form.setValue(
        "type",
        types.result.data.find((type) => type.code === typeCode) ?? null,
        { shouldDirty: false },
      );
  }, [form, typeCode, types.result.data]);
  return (
    <Box>
      <Typography variant="h4" fontWeight={700} mb={3}>
        Thêm nội dung hình ảnh
      </Typography>
      <Create
        title=""
        isLoading={form.refineCore.formLoading}
        wrapperProps={{
          sx: { p: 0, bgcolor: "transparent", boxShadow: "none" },
        }}
        headerProps={{ sx: { display: "none" } }}
        footerButtons={null}
      >
        <ImageContentForm
          form={form}
          types={types.result.data}
          onSubmit={(values) => void onFinishPayload(createPayload(values))}
        />
      </Create>
    </Box>
  );
};

export const ImageContentEdit = () => {
  const { list } = useNavigation();
  const form = useForm<IImageContent, HttpError, IImageContentFormValues>({
    refineCoreProps: {
      resource: "image-contents",
      action: "edit",
      redirect: false,
      mutationMeta: { method: "put" },
      onMutationSuccess: () => list("image-contents"),
    },
    warnWhenUnsavedChanges: true,
  });
  const item = form.refineCore.query?.data?.data;
  const types = useList<IImageContentType>({
    resource: "image-content-types",
    dataProviderName: "imageContent",
    pagination: { currentPage: 1, pageSize: 100 },
  });
  const onFinishPayload = form.refineCore.onFinish as unknown as (
    values: IUpdateImageContentPayload,
  ) => Promise<unknown>;
  useEffect(() => {
    if (item && !form.getValues("media"))
      form.reset({
        type: types.result.data.find((type) => type.code === item.typeCode) ?? {
          id: 0,
          code: item.typeCode,
          name: item.typeName,
          status: "ACTIVE",
          sortOrder: 0,
          maxItems: null,
          fieldConfig: item.fieldConfig,
          itemCount: 0,
          createdAt: "",
          updatedAt: "",
        },
        media: item.media,
        name: item.name ?? "",
        description: item.description ?? "",
        secondaryDescription: item.secondaryDescription ?? "",
        url: item.url ?? "",
        status: item.status,
      });
  }, [form, item, types.result.data]);
  const update = (
    values: IImageContentFormValues,
  ): IUpdateImageContentPayload => {
    const result: IUpdateImageContentPayload = {
      mediaId: values.media?.id ?? 0,
      status: values.status,
    };
    const config = values.type?.fieldConfig ?? item?.fieldConfig;
    if (config?.name.enabled) result.name = compact(values.name);
    if (config?.description.enabled)
      result.description = compact(values.description);
    if (config?.secondaryDescription.enabled)
      result.secondaryDescription = compact(values.secondaryDescription);
    if (config?.url.enabled) result.url = compact(values.url);
    return result;
  };
  return (
    <Box>
      <Typography variant="h4" fontWeight={700} mb={3}>
        Chỉnh sửa nội dung hình ảnh
      </Typography>
      <Edit
        title=""
        isLoading={form.refineCore.formLoading}
        wrapperProps={{
          sx: { p: 0, bgcolor: "transparent", boxShadow: "none" },
        }}
        headerProps={{ sx: { display: "none" } }}
        footerButtons={null}
      >
        <ImageContentForm
          form={form}
          types={types.result.data}
          isEdit
          item={item}
          onSubmit={(values) => void onFinishPayload(update(values))}
        />
      </Edit>
    </Box>
  );
};

import { useEffect } from "react";
import { Create, Edit } from "@refinedev/mui";
import { useForm } from "@refinedev/react-hook-form";
import { HttpError, useNavigation } from "@refinedev/core";
import { Box, Button, Stack, Typography } from "@mui/material";
import { ImageContentTypeForm } from "./components/image-content-type-form";
import { EMPTY_FIELD_CONFIG, ICreateImageContentTypePayload, IImageContentType, IImageContentTypeFormValues, IUpdateImageContentTypePayload } from "./image-content-types";

const defaults: IImageContentTypeFormValues = { code: "", name: "", status: "ACTIVE", sortOrder: 0, hasItemLimit: false, maxItems: null, fieldConfig: EMPTY_FIELD_CONFIG };
const payload = (values: IImageContentTypeFormValues): ICreateImageContentTypePayload => ({ code: values.code.trim().toUpperCase(), name: values.name.trim(), status: values.status, sortOrder: Number(values.sortOrder) || 0, maxItems: values.hasItemLimit && values.maxItems ? Number(values.maxItems) : null, fieldConfig: values.fieldConfig });
const Actions = ({ saveButtonProps, onCancel, loading, edit }: { saveButtonProps: Record<string, unknown>; onCancel: () => void; loading: boolean; edit?: boolean }) => <Stack direction="row" justifyContent="flex-end" spacing={1.5} mt={3}><Button variant="outlined" onClick={onCancel}>Hủy</Button><Button variant="contained" color={edit ? "success" : "primary"} {...saveButtonProps} disabled={loading}>Lưu</Button></Stack>;

export const ImageContentTypeCreate = () => {
  const { list } = useNavigation();
  const form = useForm<IImageContentType, HttpError, IImageContentTypeFormValues>({ refineCoreProps: { resource: "image-content-types", action: "create", redirect: false, onMutationSuccess: () => list("image-content-types") }, defaultValues: defaults, warnWhenUnsavedChanges: true });
  const { refineCore: { formLoading } } = form;
  const onFinishPayload = form.refineCore.onFinish as unknown as (values: ICreateImageContentTypePayload) => Promise<unknown>;
  return <Box><Typography variant="h4" fontWeight={700} mb={3}>Thêm loại nội dung hình ảnh</Typography><Create title="" isLoading={formLoading} wrapperProps={{ sx: { p: 0, bgcolor: "transparent", boxShadow: "none" } }} headerProps={{ sx: { display: "none" } }} footerButtons={null}><ImageContentTypeForm form={form} /><Actions saveButtonProps={{ onClick: () => void form.handleSubmit((values) => onFinishPayload(payload(values)))(), variant: "contained" }} onCancel={() => list("image-content-types")} loading={formLoading} /></Create></Box>;
};

export const ImageContentTypeEdit = () => {
  const { list } = useNavigation();
  const form = useForm<IImageContentType, HttpError, IImageContentTypeFormValues>({ refineCoreProps: { resource: "image-content-types", action: "edit", redirect: false, onMutationSuccess: () => list("image-content-types") }, warnWhenUnsavedChanges: true });
  const { refineCore: { formLoading, query } } = form;
  useEffect(() => { const type = query?.data?.data; if (type && !form.getValues("code")) form.reset({ code: type.code, name: type.name, status: type.status, sortOrder: type.sortOrder, hasItemLimit: type.maxItems !== null, maxItems: type.maxItems, fieldConfig: type.fieldConfig }); }, [form, query]);
  const onFinishPayload = form.refineCore.onFinish as unknown as (values: IUpdateImageContentTypePayload) => Promise<unknown>;
  return <Box><Typography variant="h4" fontWeight={700} mb={3}>Chỉnh sửa loại nội dung hình ảnh</Typography><Edit title="" isLoading={formLoading} wrapperProps={{ sx: { p: 0, bgcolor: "transparent", boxShadow: "none" } }} headerProps={{ sx: { display: "none" } }} footerButtons={null}><ImageContentTypeForm form={form} isEdit /><Actions saveButtonProps={{ onClick: () => void form.handleSubmit((values) => { const { code: _code, ...update } = payload(values); void onFinishPayload(update); })(), variant: "contained" }} onCancel={() => list("image-content-types")} loading={formLoading} edit /></Edit></Box>;
};

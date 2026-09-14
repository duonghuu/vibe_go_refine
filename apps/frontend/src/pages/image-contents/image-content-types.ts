import { BaseRecord } from "@refinedev/core";

export type ImageContentStatus = "ACTIVE" | "INACTIVE";
export type ImageContentFieldKey = "name" | "description" | "secondaryDescription" | "url";

export interface IImageContentFieldRule {
  enabled: boolean;
  required: boolean;
}

export interface IImageContentFieldConfig {
  name: IImageContentFieldRule;
  description: IImageContentFieldRule;
  secondaryDescription: IImageContentFieldRule;
  url: IImageContentFieldRule;
}

export interface IImageContentType extends BaseRecord {
  id: number;
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface IImageContentTypeFormValues {
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  hasItemLimit: boolean;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
}

export interface ICreateImageContentTypePayload {
  code: string;
  name: string;
  status: ImageContentStatus;
  sortOrder: number;
  maxItems: number | null;
  fieldConfig: IImageContentFieldConfig;
}

export type IUpdateImageContentTypePayload = Omit<ICreateImageContentTypePayload, "code">;

export interface IImageContentMedia extends BaseRecord {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType: string;
  status: "temporary" | "attached";
}

export interface IImageContent extends BaseRecord {
  id: number;
  typeCode: string;
  typeName: string;
  fieldConfig: IImageContentFieldConfig;
  mediaId: number;
  media: IImageContentMedia;
  name: string | null;
  description: string | null;
  secondaryDescription: string | null;
  url: string | null;
  sortOrder: number;
  status: ImageContentStatus;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface IImageContentFormValues {
  type: IImageContentType | null;
  media: IImageContentMedia | null;
  name: string;
  description: string;
  secondaryDescription: string;
  url: string;
  status: ImageContentStatus;
}

export interface ICreateImageContentPayload {
  typeCode: string;
  mediaId: number;
  name: string | null;
  description: string | null;
  secondaryDescription: string | null;
  url: string | null;
  status: ImageContentStatus;
}

export interface IUpdateImageContentPayload {
  mediaId: number;
  name?: string | null;
  description?: string | null;
  secondaryDescription?: string | null;
  url?: string | null;
  status: ImageContentStatus;
}

export interface IImageContentOrderPayload {
  typeCode: string;
  items: Array<{ id: number; sortOrder: number }>;
}

export type ImageContentErrorCode =
  | "VALIDATION_ERROR" | "TYPE_NOT_FOUND" | "IMAGE_CONTENT_NOT_FOUND" | "MEDIA_NOT_FOUND"
  | "TYPE_CODE_CONFLICT" | "TYPE_CONFIG_CONFLICT" | "TYPE_IN_USE" | "MAX_ITEMS_EXCEEDED"
  | "FIELD_REQUIRED" | "FIELD_NOT_ENABLED" | "MEDIA_TYPE_INVALID" | "FORBIDDEN"
  | "ORDER_CONFLICT" | "NOT_FOUND" | "INTERNAL_ERROR";

export interface IImageContentApiError {
  error: string;
  code: ImageContentErrorCode;
  message: string;
  fields?: Partial<Record<ImageContentFieldKey | "code" | "maxItems", string>>;
}

export const IMAGE_CONTENT_FIELD_KEYS: ImageContentFieldKey[] = ["name", "description", "secondaryDescription", "url"];

export const EMPTY_FIELD_CONFIG: IImageContentFieldConfig = {
  name: { enabled: true, required: false },
  description: { enabled: false, required: false },
  secondaryDescription: { enabled: false, required: false },
  url: { enabled: false, required: false },
};

export const fieldLabel: Record<ImageContentFieldKey, string> = {
  name: "Tên",
  description: "Mô tả chính",
  secondaryDescription: "Mô tả phụ",
  url: "Liên kết",
};

export const mediaPreviewUrl = (media: IImageContentMedia | null | undefined): string =>
  media?.thumbnailUrl ?? media?.mediumUrl ?? media?.originalUrl ?? "";


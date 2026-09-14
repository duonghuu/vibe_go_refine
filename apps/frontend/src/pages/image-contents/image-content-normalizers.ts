import { BACKEND_URL } from "../../providers/constants";
import {
  EMPTY_FIELD_CONFIG, IImageContent, IImageContentFieldConfig, IImageContentMedia, IImageContentType,
  ImageContentStatus,
} from "./image-content-types";

type JsonRecord = Record<string, unknown>;
const isRecord = (value: unknown): value is JsonRecord => typeof value === "object" && value !== null;
const stringValue = (record: JsonRecord, key: string, fallback = ""): string => typeof record[key] === "string" ? record[key] as string : fallback;
const numberValue = (record: JsonRecord, key: string, fallback = 0): number => typeof record[key] === "number" ? record[key] as number : fallback;
const nullableString = (record: JsonRecord, key: string): string | null => typeof record[key] === "string" ? record[key] as string : null;
const statusValue = (value: unknown): ImageContentStatus => value === "INACTIVE" ? "INACTIVE" : "ACTIVE";
const resolveUrl = (value: string): string => {
  if (!value) return "";
  try { return new URL(value, BACKEND_URL).toString(); } catch { return value; }
};

export const normalizeFieldConfig = (value: unknown): IImageContentFieldConfig => {
  if (!isRecord(value)) return EMPTY_FIELD_CONFIG;
  const readRule = (key: string) => {
    const rule = isRecord(value[key]) ? value[key] : {};
    return { enabled: rule.enabled === true, required: rule.required === true && rule.enabled === true };
  };
  return { name: readRule("name"), description: readRule("description"), secondaryDescription: readRule("secondaryDescription"), url: readRule("url") };
};

export const normalizeMedia = (value: unknown): IImageContentMedia => {
  const record = isRecord(value) ? value : {};
  return {
    id: numberValue(record, "id"),
    fileName: stringValue(record, "fileName", stringValue(record, "file_name", "media")),
    originalUrl: resolveUrl(stringValue(record, "originalUrl", stringValue(record, "original_url"))),
    thumbnailUrl: resolveUrl(stringValue(record, "thumbnailUrl", stringValue(record, "thumbnail_url"))) || undefined,
    mediumUrl: resolveUrl(stringValue(record, "mediumUrl", stringValue(record, "medium_url"))) || undefined,
    mimeType: stringValue(record, "mimeType", "image/*"),
    status: stringValue(record, "status") === "temporary" ? "temporary" : "attached",
  };
};

export const normalizeImageContentType = (value: unknown): IImageContentType => {
  const record = isRecord(value) ? value : {};
  const maxItems = typeof record.maxItems === "number" ? record.maxItems : null;
  return {
    id: numberValue(record, "id"), code: stringValue(record, "code"), name: stringValue(record, "name"),
    status: statusValue(record.status), sortOrder: numberValue(record, "sortOrder", numberValue(record, "sort_order")),
    maxItems, fieldConfig: normalizeFieldConfig(record.fieldConfig ?? record.field_config),
    itemCount: numberValue(record, "itemCount", numberValue(record, "item_count")),
    createdAt: stringValue(record, "createdAt", stringValue(record, "created_at")),
    updatedAt: stringValue(record, "updatedAt", stringValue(record, "updated_at")),
  };
};

export const normalizeImageContent = (value: unknown): IImageContent => {
  const record = isRecord(value) ? value : {};
  return {
    id: numberValue(record, "id"), typeCode: stringValue(record, "typeCode", stringValue(record, "type_code")),
    typeName: stringValue(record, "typeName", stringValue(record, "type_name")),
    fieldConfig: normalizeFieldConfig(record.fieldConfig ?? record.field_config),
    mediaId: numberValue(record, "mediaId", numberValue(record, "media_id")), media: normalizeMedia(record.media),
    name: nullableString(record, "name"), description: nullableString(record, "description"),
    secondaryDescription: nullableString(record, "secondaryDescription"), url: nullableString(record, "url") ?? nullableString(record, "targetUrl"),
    sortOrder: numberValue(record, "sortOrder", numberValue(record, "sort_order")), status: statusValue(record.status),
    createdBy: numberValue(record, "createdBy", numberValue(record, "created_by")),
    createdAt: stringValue(record, "createdAt", stringValue(record, "created_at")), updatedAt: stringValue(record, "updatedAt", stringValue(record, "updated_at")),
  };
};

export const unwrapData = (value: unknown): unknown => isRecord(value) && "data" in value ? value.data : value;
export const unwrapList = (value: unknown): { data: unknown[]; total: number } => {
  const record = isRecord(value) ? value : {};
  const data = Array.isArray(record.data) ? record.data : Array.isArray(value) ? value : [];
  return { data, total: typeof record.total === "number" ? record.total : data.length };
};


import { BaseRecord, CreateParams, DataProvider, DeleteOneParams, GetListParams, GetOneParams, HttpError, UpdateParams } from "@refinedev/core";
import { API_URL } from "./constants";
import { customRequest } from "./data";
import { normalizeImageContent, normalizeImageContentType, normalizeMedia, unwrapData, unwrapList } from "../pages/image-contents/image-content-normalizers";
import { IImageContent, IImageContentMedia, IImageContentType } from "../pages/image-contents/image-content-types";

type JsonRecord = Record<string, unknown>;
const isRecord = (value: unknown): value is JsonRecord => typeof value === "object" && value !== null;
const filterValue = (params: GetListParams, field: string): string => {
  const filter = params.filters?.find((item) => "field" in item && item.field === field);
  return filter && "value" in filter && (typeof filter.value === "string" || typeof filter.value === "number") ? String(filter.value) : "";
};
const queryValue = (params: GetListParams, field: string): string => {
  const value = params.meta?.[field];
  return typeof value === "string" || typeof value === "number" ? String(value) : "";
};
const responseError = async (response: Response): Promise<HttpError> => {
  const body: unknown = await response.json().catch(() => ({}));
  const record = isRecord(body) ? body : {};
  const message = typeof record.message === "string" ? record.message : typeof record.error === "string" ? record.error : "Không thể xử lý yêu cầu.";
  return { message, statusCode: response.status, code: typeof record.code === "string" ? record.code : undefined, fields: isRecord(record.fields) ? record.fields : undefined } as HttpError;
};
const read = async (response: Response): Promise<unknown> => { if (!response.ok) throw await responseError(response); return response.status === 204 ? undefined : response.json().catch(() => ({})); };
const listQuery = (params: GetListParams, kind: "types" | "contents" | "media"): string => {
  const currentPage = params.pagination?.currentPage ?? 1; const pageSize = params.pagination?.pageSize ?? 20;
  const query = new URLSearchParams();
  if (kind === "media") { query.set("current", String(currentPage)); query.set("pageSize", String(pageSize)); query.set("mimeType", "image"); }
  else { query.set("_start", String((currentPage - 1) * pageSize)); query.set("_end", String(currentPage * pageSize)); }
  const q = filterValue(params, "q") || filterValue(params, "title_like"); if (q) query.set("q", q);
  const status = filterValue(params, "status"); if (status) query.set("status", status);
  const typeCode = filterValue(params, "typeCode"); if (typeCode) query.set("typeCode", typeCode);
  const sorter = params.sorters?.[0]; if (sorter) { query.set("_sort", String(sorter.field)); query.set("_order", sorter.order === "asc" ? "ASC" : "DESC"); }
  const search = queryValue(params, "search"); if (search) query.set("q", search);
  return query.toString();
};

export const imageContentDataProvider: DataProvider = {
  getApiUrl: () => API_URL,
  getList: async <TData extends BaseRecord = BaseRecord>(params: GetListParams) => {
    const resource = params.resource;
    const kind = resource === "image-content-types" ? "types" : resource === "media" ? "media" : "contents";
    const path = kind === "types" ? "image-content-types" : kind === "media" ? "media" : "image-contents";
    const body = await read(await customRequest({ url: `${API_URL}/admin/${path}?${listQuery(params, kind)}` }));
    const list = unwrapList(body);
    const data = kind === "types" ? list.data.map(normalizeImageContentType) : kind === "media" ? list.data.map(normalizeMedia) : list.data.map(normalizeImageContent);
    return { data: data as unknown as TData[], total: list.total };
  },
  getOne: async <TData extends BaseRecord = BaseRecord>(params: GetOneParams) => {
    const path = params.resource === "image-content-types" ? "image-content-types" : params.resource === "media" ? "media" : "image-contents";
    const body = unwrapData(await read(await customRequest({ url: `${API_URL}/admin/${path}/${params.id}` })));
    const data = path === "image-content-types" ? normalizeImageContentType(body) : path === "media" ? normalizeMedia(body) : normalizeImageContent(body);
    return { data: data as unknown as TData };
  },
  create: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: CreateParams<TVariables>) => {
    const path = params.resource === "image-content-types" ? "image-content-types" : "image-contents";
    const body = unwrapData(await read(await customRequest({ url: `${API_URL}/admin/${path}`, method: "POST", payload: params.variables })));
    return { data: (path === "image-content-types" ? normalizeImageContentType(body) : normalizeImageContent(body)) as unknown as TData };
  },
  update: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: UpdateParams<TVariables>) => {
    const path = params.resource === "image-content-types" ? "image-content-types" : "image-contents";
    const body = unwrapData(await read(await customRequest({ url: `${API_URL}/admin/${path}/${params.id}`, method: "PUT", payload: params.variables })));
    return { data: (path === "image-content-types" ? normalizeImageContentType(body) : normalizeImageContent(body)) as unknown as TData };
  },
  deleteOne: async <TData extends BaseRecord = BaseRecord, TVariables = unknown>(params: DeleteOneParams<TVariables>) => {
    const path = params.resource === "image-content-types" ? "image-content-types" : "image-contents";
    await read(await customRequest({ url: `${API_URL}/admin/${path}/${params.id}`, method: "DELETE" }));
    return { data: { id: params.id } as TData };
  },
};

export const imageContentMediaFromUpload = async (file: File): Promise<IImageContentMedia> => {
  const form = new FormData(); form.append("file", file);
  const body = unwrapData(await read(await customRequest({ url: `${API_URL}/media/upload`, method: "POST", payload: form })));
  return normalizeMedia(body);
};

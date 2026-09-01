import { BaseRecord, CreateParams, DataProvider, DeleteOneParams, GetListParams, GetOneParams, HttpError, UpdateParams } from "@refinedev/core";
import { API_URL } from "./constants";
import { customRequest } from "./data";
import { PageSectionItemType } from "../pages/pages/sections/page-section-types";

type JsonRecord = Record<string, unknown>;
const isRecord = (value: unknown): value is JsonRecord => typeof value === "object" && value !== null;
const numberMeta = (params: GetListParams, key: string): number | undefined => { const value: unknown = params.meta?.[key]; return typeof value === "number" ? value : undefined; };
const stringMeta = (params: GetListParams, key: string): string | undefined => { const value: unknown = params.meta?.[key]; return typeof value === "string" ? value : undefined; };
const filterValue = (params: GetListParams, field: string): string => { const found = params.filters?.find((item) => "field" in item && item.field === field); return found && "value" in found && typeof found.value === "string" ? found.value : ""; };
const pageParams = (params: GetListParams): string => { const current = params.pagination?.currentPage ?? 1; const pageSize = params.pagination?.pageSize ?? 100; return `current=${current}&pageSize=${pageSize}`; };
const unwrap = (body: unknown): unknown => isRecord(body) && "data" in body ? body.data : body;
const readJson = async (response: Response): Promise<unknown> => {
  const body: unknown = await response.json().catch(() => ({}));
  if (!response.ok) { const record = isRecord(body) ? body : {}; const message = typeof record.message === "string" ? record.message : "Không thể xử lý yêu cầu."; const error: HttpError = { message, statusCode: response.status }; throw error; }
  return body;
};
const pageSectionPath = (pageId: number): string => `${API_URL}/admin/pages/${pageId}/sections`;
const sourceOption = (value: unknown, itemType: PageSectionItemType): BaseRecord => {
  const item = isRecord(value) ? value : {};
  const id = typeof item.id === "number" ? item.id : 0;
  const name = typeof item.name === "string" ? item.name : "";
  const title = typeof item.title === "string" ? item.title : "";
  const slug = typeof item.slug === "string" ? item.slug : "";
  const fileName = typeof item.fileName === "string" ? item.fileName : "";
  const originalUrl = typeof item.originalUrl === "string" ? item.originalUrl : "";
  const thumbnailUrl = typeof item.thumbnailUrl === "string" ? item.thumbnailUrl : "";
  const mediumUrl = typeof item.mediumUrl === "string" ? item.mediumUrl : "";
  const label = itemType === "CATEGORY" ? name : itemType === "POST" ? title : fileName;
  const secondary = itemType === "MEDIA" ? (typeof item.mimeType === "string" ? item.mimeType : "") : slug;
  return { id, itemType, label, secondary, slug, imageUrl: thumbnailUrl || (typeof item.imageUrl === "string" ? item.imageUrl : ""), originalUrl, thumbnailUrl, mediumUrl, mimeType: typeof item.mimeType === "string" ? item.mimeType : undefined, status: typeof item.status === "string" ? item.status : undefined, mediaStatus: typeof item.status === "string" ? item.status : undefined, typeCode: typeof item.typeCode === "string" ? item.typeCode : undefined };
};

export const pageSectionDataProvider: DataProvider = {
  getList: async <TData extends BaseRecord = BaseRecord>(params: GetListParams) => {
    const pageId = numberMeta(params, "pageId");
    if (params.resource === "page-sections" && pageId !== undefined) {
      const response = await customRequest({ url: `${pageSectionPath(pageId)}?${pageParams(params)}` });
      const body = await readJson(response); const data = unwrap(body); const record = isRecord(body) ? body : {};
      return { data: (Array.isArray(data) ? data : []) as TData[], total: typeof record.total === "number" ? record.total : Array.isArray(data) ? data.length : 0 };
    }
    if (params.resource === "page-section-items" && pageId !== undefined) {
      const sectionId = numberMeta(params, "sectionId"); const collection = stringMeta(params, "collection");
      if (sectionId === undefined || collection === undefined) return { data: [] as TData[], total: 0 };
      const response = await customRequest({ url: `${pageSectionPath(pageId)}/${sectionId}/items?collection=${encodeURIComponent(collection)}&${pageParams(params)}` });
      const body = await readJson(response); const data = unwrap(body); const record = isRecord(body) ? body : {};
      return { data: (Array.isArray(data) ? data : []) as TData[], total: typeof record.total === "number" ? record.total : Array.isArray(data) ? data.length : 0 };
    }
    if (params.resource === "section-picker-options") {
      const itemType = stringMeta(params, "itemType") as PageSectionItemType | undefined; const search = filterValue(params, "search");
      const endpoint = itemType === "CATEGORY" ? "categories" : itemType === "POST" ? "posts" : "media";
      const current = params.pagination?.currentPage ?? 1; const size = params.pagination?.pageSize ?? 5;
      const query = itemType === "CATEGORY" ? new URLSearchParams({ _start: String((current - 1) * size), _end: String(current * size) }) : new URLSearchParams({ current: String(current), pageSize: String(size) });
      if (search) query.set(itemType === "POST" ? "title_like" : "q", search);
      if (itemType === "MEDIA") query.set("mimeType", "image");
      const response = await customRequest({ url: `${API_URL}/admin/${endpoint}?${query.toString()}` }); const body = await readJson(response); const data = unwrap(body); const record = isRecord(body) ? body : {};
      return { data: (Array.isArray(data) ? data.map((item) => sourceOption(item, itemType ?? "MEDIA")) : []) as TData[], total: typeof record.total === "number" ? record.total : Array.isArray(data) ? data.length : 0 };
    }
    return { data: [] as TData[], total: 0 };
  },
  getOne: async <TData extends BaseRecord = BaseRecord>(params: GetOneParams) => {
    const pageId = numberMeta(params as unknown as GetListParams, "pageId"); const sectionId = Number(params.id);
    if (pageId === undefined) throw { message: "Thiếu Page ID.", statusCode: 400 } satisfies HttpError;
    const response = await customRequest({ url: `${pageSectionPath(pageId)}/${sectionId}` }); const body = await readJson(response); return { data: unwrap(body) as TData };
  },
  create: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: CreateParams<TVariables>) => {
    const values = params.variables as JsonRecord; const pageId = typeof values.pageId === "number" ? values.pageId : undefined;
    if (params.resource !== "page-sections" || pageId === undefined) throw { message: "Thiếu Page ID.", statusCode: 400 } satisfies HttpError;
    const payload = { key: values.key, name: values.name, title: values.title ?? "", description: values.description ?? "", backgroundColor: values.backgroundColor ?? null, backgroundMediaId: values.backgroundMediaId ?? null, featureMediaId: values.featureMediaId ?? null, sortOrder: values.sortOrder, status: values.status };
    const response = await customRequest({ url: pageSectionPath(pageId), method: "POST", payload }); const body = await readJson(response); return { data: unwrap(body) as TData };
  },
  update: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: UpdateParams<TVariables>) => {
    const values = params.variables as JsonRecord; const pageId = typeof values.pageId === "number" ? values.pageId : numberMeta(params as GetListParams, "pageId");
    if (params.resource === "page-section-order") { if (pageId === undefined) throw { message: "Thiếu Page ID.", statusCode: 400 } satisfies HttpError; const response = await customRequest({ url: `${API_URL}/admin/pages/${pageId}/section-order`, method: "PUT", payload: { sections: values.sections } }); const body = await readJson(response); return { data: unwrap(body) as TData }; }
    if (params.resource === "page-section-items") { const sectionId = typeof values.sectionId === "number" ? values.sectionId : undefined; const collection = typeof values.collection === "string" ? values.collection : undefined; if (pageId === undefined || sectionId === undefined || collection === undefined) throw { message: "Thiếu thông tin collection.", statusCode: 400 } satisfies HttpError; const itemIds = Array.isArray(values.itemIds) ? values.itemIds.filter((item): item is number => typeof item === "number") : []; const itemType = typeof values.itemType === "string" ? values.itemType : "MEDIA"; const response = await customRequest({ url: `${pageSectionPath(pageId)}/${sectionId}/items/${encodeURIComponent(collection)}`, method: "PUT", payload: { itemType, items: itemIds.map((itemId, sortOrder) => ({ itemId, sortOrder })) } }); const body = await readJson(response); return { data: unwrap(body) as TData }; }
    const sectionId = Number(params.id); if (pageId === undefined || !Number.isInteger(sectionId)) throw { message: "Thiếu thông tin Section.", statusCode: 400 } satisfies HttpError; const payload = { key: values.key, name: values.name, title: values.title ?? "", description: values.description ?? "", backgroundColor: values.backgroundColor ?? null, backgroundMediaId: values.backgroundMediaId ?? null, featureMediaId: values.featureMediaId ?? null, status: values.status }; const response = await customRequest({ url: `${pageSectionPath(pageId)}/${sectionId}`, method: "PUT", payload }); const body = await readJson(response); return { data: unwrap(body) as TData };
  },
  deleteOne: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: DeleteOneParams<TVariables>) => { const pageId = numberMeta(params as unknown as GetListParams, "pageId"); if (pageId === undefined) throw { message: "Thiếu Page ID.", statusCode: 400 } satisfies HttpError; const response = await customRequest({ url: `${pageSectionPath(pageId)}/${Number(params.id)}`, method: "DELETE" }); const body = response.status === 204 ? {} : await readJson(response); return { data: unwrap(body) as TData }; },
  getApiUrl: () => `${API_URL}/admin`,
};

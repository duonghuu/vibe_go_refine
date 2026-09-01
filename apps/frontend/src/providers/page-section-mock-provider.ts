import {
  BaseRecord,
  CreateParams,
  CrudFilter,
  DataProvider,
  DeleteOneParams,
  GetListParams,
  GetOneParams,
  HttpError,
  UpdateParams,
} from "@refinedev/core";
import rawDataset from "../mock-data/page-sections.json";
import {
  IPageSection,
  IPageSectionItem,
  IPageSectionMockDataset,
  IReorderPageSectionsPayload,
  ISectionPickerOption,
  ISyncSectionCollectionPayload,
  IUpsertPageSectionPayload,
  PageSectionItemType,
  pickerOptionToSectionItem,
} from "../pages/pages/sections/page-section-types";

const dataset = rawDataset as IPageSectionMockDataset;
let sections: IPageSection[] = dataset.sections.map((section) => ({
  ...section,
  collections: section.collections.map((collection) => ({ ...collection })),
}));
let items: IPageSectionItem[] = dataset.items.map((item) => ({ ...item, data: { ...item.data } })) as IPageSectionItem[];
const options: ISectionPickerOption[] = dataset.options.map((option) => ({ ...option }));
let nextSectionId = Math.max(...sections.map((section) => section.id), 0) + 1;
let nextItemId = Math.max(...items.map((item) => item.id), 0) + 1;

const mockError = (message: string, statusCode = 400): HttpError => ({ message, statusCode });

const asData = <TData extends BaseRecord>(record: BaseRecord): TData => record as TData;
const asDataList = <TData extends BaseRecord>(records: BaseRecord[]): TData[] => records as TData[];

const readMetaNumber = (params: GetListParams, key: string): number | undefined => {
  const value: unknown = params.meta?.[key];
  return typeof value === "number" && Number.isFinite(value) ? value : undefined;
};

const readMetaString = (params: GetListParams, key: string): string | undefined => {
  const value: unknown = params.meta?.[key];
  return typeof value === "string" ? value : undefined;
};

const getFilterValue = (filters: CrudFilter[] | undefined, field: string): unknown => {
  const filter = filters?.find((candidate) => "field" in candidate && candidate.field === field);
  return filter && "value" in filter ? filter.value as unknown : undefined;
};

const paginate = <T,>(records: T[], params: GetListParams): T[] => {
  const currentPage = params.pagination?.currentPage ?? 1;
  const pageSize = params.pagination?.pageSize ?? (records.length || 1);
  const start = (currentPage - 1) * pageSize;
  return records.slice(start, start + pageSize);
};

const refreshSectionCollections = (sectionId: number): void => {
  const section = sections.find((candidate) => candidate.id === sectionId);
  if (!section) return;
  const grouped = new Map<string, { itemType: PageSectionItemType; total: number }>();
  items.filter((item) => item.sectionId === sectionId).forEach((item) => {
    const current = grouped.get(item.collection);
    grouped.set(item.collection, { itemType: item.itemType, total: (current?.total ?? 0) + 1 });
  });
  section.collections = Array.from(grouped.entries()).map(([collection, value]) => ({
    collection,
    itemType: value.itemType,
    total: value.total,
  }));
};

const isUpsertPayload = (value: unknown): value is IUpsertPageSectionPayload => {
  if (typeof value !== "object" || value === null) return false;
  const record = value as Record<string, unknown>;
  return typeof record.pageId === "number" && typeof record.key === "string" && typeof record.name === "string";
};

const isReorderPayload = (value: unknown): value is IReorderPageSectionsPayload => {
  if (typeof value !== "object" || value === null) return false;
  const record = value as Record<string, unknown>;
  return typeof record.pageId === "number" && Array.isArray(record.sections);
};

const isSyncPayload = (value: unknown): value is ISyncSectionCollectionPayload => {
  if (typeof value !== "object" || value === null) return false;
  const record = value as Record<string, unknown>;
  return typeof record.sectionId === "number" && typeof record.collection === "string" && Array.isArray(record.itemIds);
};

export const pageSectionMockDataProvider: DataProvider = {
  getList: async <TData extends BaseRecord = BaseRecord>(params: GetListParams) => {
    if (params.resource === "page-sections") {
      const pageId = readMetaNumber(params, "pageId");
      const data = sections
        .filter((section) => pageId === undefined || section.pageId === pageId)
        .sort((left, right) => left.sortOrder - right.sortOrder || left.id - right.id);
      return { data: asDataList<TData>(paginate(data, params)), total: data.length };
    }
    if (params.resource === "page-section-items") {
      const sectionId = readMetaNumber(params, "sectionId");
      const collection = readMetaString(params, "collection");
      const data = items
        .filter((item) => (sectionId === undefined || item.sectionId === sectionId) && (!collection || item.collection === collection))
        .sort((left, right) => left.sortOrder - right.sortOrder || left.id - right.id);
      return { data: asDataList<TData>(paginate(data, params)), total: data.length };
    }
    if (params.resource === "section-picker-options") {
      const itemType = readMetaString(params, "itemType");
      const searchValue = getFilterValue(params.filters, "search");
      const search = typeof searchValue === "string" ? searchValue.trim().toLocaleLowerCase("vi") : "";
      const data = options.filter((option) =>
        (!itemType || option.itemType === itemType)
        && (!search || option.label.toLocaleLowerCase("vi").includes(search) || option.secondary.toLocaleLowerCase("vi").includes(search)),
      );
      return { data: asDataList<TData>(paginate(data, params)), total: data.length };
    }
    return { data: [], total: 0 };
  },

  getOne: async <TData extends BaseRecord = BaseRecord>(params: GetOneParams) => {
    const { resource, id } = params;
    const record = resource === "page-sections"
      ? sections.find((section) => section.id === Number(id))
      : items.find((item) => item.id === Number(id));
    if (!record) throw mockError("Không tìm thấy dữ liệu Section.", 404);
    return { data: asData<TData>(record) };
  },

  create: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: CreateParams<TVariables>) => {
    const { resource, variables } = params;
    const value: unknown = variables;
    if (resource !== "page-sections" || !isUpsertPayload(value)) {
      throw mockError("Dữ liệu tạo Section không hợp lệ.");
    }
    if (sections.some((section) => section.pageId === value.pageId && section.key === value.key)) {
      throw mockError("Key Section đã tồn tại trong trang.", 409);
    }
    const now = new Date().toISOString();
    const section: IPageSection = {
      ...value,
      id: nextSectionId++,
      collections: [],
      createdAt: now,
      updatedAt: now,
    };
    sections.push(section);
    return { data: asData<TData>(section) };
  },

  update: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: UpdateParams<TVariables>) => {
    const { resource, id, variables } = params;
    const value: unknown = variables;
    if (resource === "page-sections" && isUpsertPayload(value)) {
      const index = sections.findIndex((section) => section.id === Number(id));
      if (index < 0) throw mockError("Không tìm thấy Section.", 404);
      if (sections.some((section) => section.id !== Number(id) && section.pageId === value.pageId && section.key === value.key)) {
        throw mockError("Key Section đã tồn tại trong trang.", 409);
      }
      sections[index] = { ...sections[index], ...value, id: Number(id), updatedAt: new Date().toISOString() };
      return { data: asData<TData>(sections[index]) };
    }
    if (resource === "page-section-order" && isReorderPayload(value)) {
      value.sections.forEach((orderItem) => {
        const section = sections.find((candidate) => candidate.id === orderItem.id && candidate.pageId === value.pageId);
        if (section) section.sortOrder = orderItem.sortOrder;
      });
      return { data: asData<TData>({ id: value.pageId }) };
    }
    if (resource === "page-section-items" && isSyncPayload(value)) {
      items = items.filter((item) => !(item.sectionId === value.sectionId && item.collection === value.collection));
      const selectedOptions = value.itemIds
        .map((itemId) => options.find((option) => option.id === itemId && option.itemType === value.itemType))
        .filter((option): option is ISectionPickerOption => Boolean(option));
      const newItems = selectedOptions.map((option, index) => pickerOptionToSectionItem(
        option,
        value.sectionId,
        value.collection,
        index,
        nextItemId++,
      ));
      items.push(...newItems);
      refreshSectionCollections(value.sectionId);
      return { data: asData<TData>({ id: String(id) }) };
    }
    throw mockError("Dữ liệu cập nhật Section không hợp lệ.");
  },

  deleteOne: async <TData extends BaseRecord = BaseRecord, TVariables = Record<string, unknown>>(params: DeleteOneParams<TVariables>) => {
    const { resource, id } = params;
    if (resource !== "page-sections") throw mockError("Resource không hỗ trợ xóa.");
    const index = sections.findIndex((section) => section.id === Number(id));
    if (index < 0) throw mockError("Không tìm thấy Section.", 404);
    const [deleted] = sections.splice(index, 1);
    items = items.filter((item) => item.sectionId !== deleted.id);
    sections
      .filter((section) => section.pageId === deleted.pageId)
      .sort((left, right) => left.sortOrder - right.sortOrder)
      .forEach((section, sortOrder) => { section.sortOrder = sortOrder; });
    return { data: asData<TData>(deleted) };
  },

  getApiUrl: () => "mock://page-sections",
};

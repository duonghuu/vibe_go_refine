import { BaseRecord } from "@refinedev/core";

export type PageSectionStatus = "ACTIVE" | "INACTIVE";
export type PageSectionItemType = "CATEGORY" | "POST" | "MEDIA";

export interface ISectionMediaSummary {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType: string;
}

export interface IPageSectionCollectionSummary {
  collection: string;
  itemType: PageSectionItemType;
  total: number;
}

export interface IPageSection extends BaseRecord {
  id: number;
  pageId: number;
  key: string;
  name: string;
  title: string | null;
  description: string | null;
  backgroundColor: string | null;
  backgroundMediaId: number | null;
  backgroundMedia: ISectionMediaSummary | null;
  featureMediaId: number | null;
  featureMedia: ISectionMediaSummary | null;
  sortOrder: number;
  status: PageSectionStatus;
  collections: IPageSectionCollectionSummary[];
  createdAt: string;
  updatedAt: string;
}

export interface IPageSectionFormValues {
  key: string;
  name: string;
  title: string;
  description: string;
  backgroundColor: string;
  backgroundMedia: ISectionMediaSummary | null;
  featureMedia: ISectionMediaSummary | null;
  status: PageSectionStatus;
}

export interface IUpsertPageSectionPayload {
  pageId: number;
  key: string;
  name: string;
  title: string | null;
  description: string | null;
  backgroundColor: string | null;
  backgroundMediaId: number | null;
  backgroundMedia: ISectionMediaSummary | null;
  featureMediaId: number | null;
  featureMedia: ISectionMediaSummary | null;
  sortOrder: number;
  status: PageSectionStatus;
}

export interface IReorderPageSectionsPayload {
  pageId: number;
  sections: Array<{ id: number; sortOrder: number }>;
}

export interface ICategorySectionData {
  id: number;
  name: string;
  slug: string;
  imageUrl?: string;
  status: "ACTIVE" | "HIDDEN";
}

export interface IPostSectionData {
  id: number;
  title: string;
  slug: string;
  typeCode: string;
  thumbnailUrl?: string;
}

export interface IMediaSectionData extends ISectionMediaSummary {
  status: "temporary" | "attached";
}

interface IPageSectionItemBase extends BaseRecord {
  id: number;
  sectionId: number;
  itemId: number;
  collection: string;
  sortOrder: number;
}

export type IPageSectionItem =
  | (IPageSectionItemBase & { itemType: "CATEGORY"; data: ICategorySectionData })
  | (IPageSectionItemBase & { itemType: "POST"; data: IPostSectionData })
  | (IPageSectionItemBase & { itemType: "MEDIA"; data: IMediaSectionData });

export interface ISectionPickerOption extends BaseRecord {
  id: number;
  itemType: PageSectionItemType;
  label: string;
  secondary: string;
  imageUrl?: string;
  slug?: string;
  status?: "ACTIVE" | "HIDDEN";
  typeCode?: string;
  originalUrl?: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType?: string;
  mediaStatus?: "temporary" | "attached";
}

export interface ISyncSectionCollectionPayload {
  sectionId: number;
  collection: string;
  itemType: PageSectionItemType;
  itemIds: number[];
}

export interface IPageSectionMockDataset {
  sections: IPageSection[];
  items: IPageSectionItem[];
  options: ISectionPickerOption[];
}

export const pickerOptionToMedia = (option: ISectionPickerOption): ISectionMediaSummary => ({
  id: option.id,
  fileName: option.label,
  originalUrl: option.originalUrl ?? option.imageUrl ?? "",
  thumbnailUrl: option.thumbnailUrl ?? option.imageUrl,
  mediumUrl: option.mediumUrl,
  mimeType: option.mimeType ?? "image/webp",
});

export const pickerOptionToSectionItem = (
  option: ISectionPickerOption,
  sectionId: number,
  collection: string,
  sortOrder: number,
  id: number,
): IPageSectionItem => {
  const base = { id, sectionId, itemId: option.id, collection, sortOrder };
  if (option.itemType === "CATEGORY") {
    return {
      ...base,
      itemType: "CATEGORY",
      data: {
        id: option.id,
        name: option.label,
        slug: option.slug ?? "",
        imageUrl: option.imageUrl,
        status: option.status ?? "ACTIVE",
      },
    };
  }
  if (option.itemType === "POST") {
    return {
      ...base,
      itemType: "POST",
      data: {
        id: option.id,
        title: option.label,
        slug: option.slug ?? "",
        typeCode: option.typeCode ?? "NEWS",
        thumbnailUrl: option.imageUrl,
      },
    };
  }
  return {
    ...base,
    itemType: "MEDIA",
    data: {
      ...pickerOptionToMedia(option),
      status: option.mediaStatus ?? "attached",
    },
  };
};

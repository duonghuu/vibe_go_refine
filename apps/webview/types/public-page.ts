export type PublicPageSectionItemType = "CATEGORY" | "POST" | "MEDIA";

export interface PublicPageSectionItemDto {
  itemType: PublicPageSectionItemType;
  itemId: number;
  sortOrder: number;
  data: Record<string, unknown>;
}

export interface PublicPageSectionDto {
  id: number;
  key: string;
  title: string;
  description: string;
  sortOrder: number;
  collections: Record<string, PublicPageSectionItemDto[]>;
}

export interface PublicPageDto {
  id: number;
  slug: string;
  sections: PublicPageSectionDto[];
}

export interface PublicPageResponse {
  data: PublicPageDto;
}

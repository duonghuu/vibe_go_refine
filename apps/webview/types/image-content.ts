export interface PublicImageContentMedia {
  id: number;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  mimeType: string;
}

export interface PublicImageContentItem {
  id: number;
  name?: string | null;
  description?: string | null;
  secondaryDescription?: string | null;
  url?: string | null;
  sortOrder: number;
  media: PublicImageContentMedia;
}

export interface PublicImageContentGroup {
  typeCode: string;
  name: string;
  items: PublicImageContentItem[];
}

export interface PublicImageContentResponse {
  data: PublicImageContentGroup[];
}

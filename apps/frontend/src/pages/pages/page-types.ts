export type PageStatus = "DRAFT" | "PUBLISHED";

export interface IPageListItem {
  id: number;
  title: string;
  slug: string;
  status: PageStatus;
  authorId: number;
  createdAt: string;
  updatedAt: string;
}

export interface IPageListResponse {
  data: IPageListItem[];
  total: number;
}

export type PageStatusFilter = PageStatus | "";

import { BaseRecord } from "@refinedev/core";

export type PageStatus = "DRAFT" | "PUBLISHED";

export interface IPageCreateFormValues {
  title: string;
  slug: string;
  content: string;
}

export interface IPageCreatePayload extends IPageCreateFormValues {
  status: PageStatus;
}

export type IPageUpdatePayload = IPageCreatePayload;

export interface IPageResponse extends BaseRecord {
  id: number;
  title: string;
  slug: string;
  content: string;
  status: PageStatus;
  authorId: number;
  createdAt: string;
  updatedAt: string;
}

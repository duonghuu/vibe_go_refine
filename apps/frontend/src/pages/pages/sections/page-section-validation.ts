import { IPageSectionFormValues, IUpsertPageSectionPayload } from "./page-section-types";

export const sectionKeyPattern = /^[a-z0-9]+(?:_[a-z0-9]+)*$/;
export const sectionColorPattern = /^#[0-9a-fA-F]{6}(?:[0-9a-fA-F]{2})?$/;

export const toSectionPayload = (
  pageId: number,
  sortOrder: number,
  values: IPageSectionFormValues,
): IUpsertPageSectionPayload => ({
  pageId,
  key: values.key.trim(),
  name: values.name.trim(),
  title: values.title.trim() || null,
  description: values.description.trim() || null,
  backgroundColor: values.backgroundColor.trim() || null,
  backgroundMediaId: values.backgroundMedia?.id ?? null,
  backgroundMedia: values.backgroundMedia,
  featureMediaId: values.featureMedia?.id ?? null,
  featureMedia: values.featureMedia,
  sortOrder,
  status: values.status,
});

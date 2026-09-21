import { useCallback, useEffect, useMemo, useState } from "react";
import { HttpError, useList, useNotification, useUpdate } from "@refinedev/core";
import { IPageSectionItem, ISectionPickerOption, ISyncSectionCollectionPayload, PageSectionItemType, pickerOptionToSectionItem } from "./page-section-types";

export const useSectionCollection = (pageId: number, sectionId: number, collection: string, itemType: PageSectionItemType) => {
  const { open } = useNotification();
  const list = useList<IPageSectionItem, HttpError>({ resource: "page-section-items", dataProviderName: "pageSections", pagination: { mode: "off" }, meta: { pageId, sectionId, collection, itemType } });
  const mutation = useUpdate<IPageSectionItem, HttpError, ISyncSectionCollectionPayload>({ dataProviderName: "pageSections" });
  const [items, setItems] = useState<IPageSectionItem[]>([]);
  const [savedIds, setSavedIds] = useState<number[]>([]);
  useEffect(() => { if (list.result.data) { setItems(list.result.data); setSavedIds(list.result.data.map((item) => item.itemId)); } }, [list.result.data]);
  const select = useCallback((options: ISectionPickerOption[]) => setItems(options.map((option, index) => pickerOptionToSectionItem(option, sectionId, collection, index, -index - 1))), [collection, sectionId]);
  const remove = useCallback((itemId: number) => setItems((current) => current.filter((item) => item.itemId !== itemId).map((item, sortOrder) => ({ ...item, sortOrder }))), []);
  const move = useCallback((from: number, to: number) => setItems((current) => { if (to < 0 || to >= current.length) return current; const next = [...current]; const [item] = next.splice(from, 1); next.splice(to, 0, item); return next.map((entry, sortOrder) => ({ ...entry, sortOrder })); }), []);
  const save = useCallback(async () => { await mutation.mutateAsync({ resource: "page-section-items", id: `${sectionId}:${collection}`, values: { pageId, sectionId, collection, itemType, itemIds: items.map((item) => item.itemId) } }); await list.query.refetch(); open?.({ type: "success", message: `Đã lưu collection ${collection}.` }); }, [collection, itemType, items, list.query, mutation, open, pageId, sectionId]);
  const reset = useCallback(() => { void list.query.refetch(); }, [list.query]);
  const isDirty = useMemo(() => items.map((item) => item.itemId).join() !== savedIds.join(), [items, savedIds]);
  return { items, query: list.query, isDirty, isSaving: mutation.mutation.isPending, select, remove, move, save, reset };
};

import { useCallback, useEffect, useMemo, useState } from "react";
import { HttpError, useCreate, useDelete, useList, useNotification, useUpdate } from "@refinedev/core";
import { IPageSection, IReorderPageSectionsPayload, IUpsertPageSectionPayload } from "./page-section-types";

export const usePageSections = (pageId: number) => {
  const { open } = useNotification();
  const list = useList<IPageSection, HttpError>({ resource: "page-sections", dataProviderName: "pageSections", pagination: { mode: "off" }, meta: { pageId } });
  const create = useCreate<IPageSection, HttpError, IUpsertPageSectionPayload>({ dataProviderName: "pageSections" });
  const update = useUpdate<IPageSection, HttpError, IUpsertPageSectionPayload>({ dataProviderName: "pageSections" });
  const reorder = useUpdate<IPageSection, HttpError, IReorderPageSectionsPayload>({ dataProviderName: "pageSections" });
  const remove = useDelete<IPageSection, HttpError>();
  const [sections, setSections] = useState<IPageSection[]>([]);
  const [savedOrder, setSavedOrder] = useState<number[]>([]);

  useEffect(() => {
    if (!list.result.data) return;
    setSections(list.result.data);
    setSavedOrder(list.result.data.map((section) => section.id));
  }, [list.result.data]);

  const move = useCallback((from: number, to: number) => {
    setSections((current) => {
      if (to < 0 || to >= current.length || from === to) return current;
      const next = [...current];
      const [item] = next.splice(from, 1);
      next.splice(to, 0, item);
      return next.map((section, sortOrder) => ({ ...section, sortOrder }));
    });
  }, []);

  const reload = useCallback(async () => { await list.query.refetch(); }, [list.query]);
  const createSection = useCallback(async (values: IUpsertPageSectionPayload) => {
    await create.mutateAsync({ resource: "page-sections", values });
    await reload();
    open?.({ type: "success", message: "Đã tạo Section." });
  }, [create, open, reload]);
  const updateSection = useCallback(async (id: number, values: IUpsertPageSectionPayload) => {
    await update.mutateAsync({ resource: "page-sections", id, values });
    await reload();
    open?.({ type: "success", message: "Đã cập nhật Section." });
  }, [open, reload, update]);
  const deleteSection = useCallback(async (id: number) => {
    await remove.mutateAsync({ resource: "page-sections", id, dataProviderName: "pageSections", meta: { pageId } });
    await reload();
    open?.({ type: "success", message: "Đã xóa Section." });
  }, [open, reload, remove]);
  const toggleSection = useCallback(async (section: IPageSection) => {
    const status = section.status === "ACTIVE" ? "INACTIVE" : "ACTIVE";
    setSections((current) => current.map((item) => item.id === section.id ? { ...item, status } : item));
    try {
      await update.mutateAsync({ resource: "page-sections", id: section.id, values: { ...section, status } });
      open?.({ type: "success", message: status === "ACTIVE" ? "Đã bật Section." : "Đã tắt Section." });
    } catch (error) {
      await reload();
      open?.({ type: "error", message: "Không thể đổi trạng thái Section.", description: error instanceof Error ? error.message : "Vui lòng thử lại." });
    }
  }, [open, reload, update]);
  const saveOrder = useCallback(async () => {
    await reorder.mutateAsync({ resource: "page-section-order", id: pageId, values: { pageId, sections: sections.map((section, sortOrder) => ({ id: section.id, sortOrder })) } });
    setSavedOrder(sections.map((section) => section.id));
    open?.({ type: "success", message: "Đã lưu thứ tự Section." });
  }, [open, pageId, reorder, sections]);
  const resetOrder = useCallback(() => {
    setSections((current) => savedOrder.map((id) => current.find((section) => section.id === id)).filter((section): section is IPageSection => Boolean(section)));
  }, [savedOrder]);

  const isOrderDirty = useMemo(() => sections.map((section) => section.id).join() !== savedOrder.join(), [savedOrder, sections]);
  const isMutating = create.mutation.isPending || update.mutation.isPending || reorder.mutation.isPending || remove.mutation.isPending;
  return { sections, query: list.query, isOrderDirty, isMutating, move, createSection, updateSection, deleteSection, toggleSection, saveOrder, resetOrder, reload };
};

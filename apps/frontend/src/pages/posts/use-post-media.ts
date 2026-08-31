import { useCallback, useEffect, useMemo, useState } from "react";
import { useNotification } from "@refinedev/core";
import {
  deleteTemporaryMedia,
  getPostMedia,
  syncPostMedia,
  uploadPostMedia,
} from "./post-media-api";
import {
  IUploadedMedia,
  PostMediaCollection,
  PostMediaState,
} from "./post-media-types";

const acceptedMimeTypes = new Set(["image/jpeg", "image/png", "image/webp"]);
const maxFileSize = 5 * 1024 * 1024;

const mediaFromItem = (item: { mediaId: number; media?: IUploadedMedia }): IUploadedMedia =>
  item.media ?? {
    id: item.mediaId,
    fileName: `media-${item.mediaId}`,
    originalUrl: "",
    status: "attached",
  };

export interface UsePostMediaResult extends PostMediaState {
  uploadThumbnail: (file: File) => Promise<void>;
  uploadGallery: (files: FileList | File[]) => Promise<void>;
  removeThumbnail: () => Promise<void>;
  removeGalleryItem: (mediaId: number) => Promise<void>;
  reorderGallery: (fromIndex: number, toIndex: number) => void;
  load: () => Promise<void>;
  sync: (postId: number) => Promise<void>;
  markDirty: () => void;
}

export const usePostMedia = (postId?: number, entityPath = "posts"): UsePostMediaResult => {
  const { open } = useNotification();
  const [thumbnail, setThumbnail] = useState<IUploadedMedia | null>(null);
  const [gallery, setGallery] = useState<IUploadedMedia[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadingCount, setUploadingCount] = useState(0);
  const [isLoading, setIsLoading] = useState(Boolean(postId));
  const [isSyncing, setIsSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isDirty, setIsDirty] = useState(false);

  const validateFile = useCallback((file: File) => {
    if (!acceptedMimeTypes.has(file.type)) {
      throw new Error("Chỉ hỗ trợ ảnh JPG, PNG hoặc WEBP.");
    }
    if (file.size > maxFileSize) {
      throw new Error("Dung lượng mỗi ảnh không được vượt quá 5MB.");
    }
  }, []);

  const notifyError = useCallback((message: string) => {
    open?.({ type: "error", message: "Không thể cập nhật hình ảnh", description: message });
  }, [open]);

  const upload = useCallback(async (file: File, collection: PostMediaCollection) => {
    validateFile(file);
    setIsUploading(true);
    setUploadingCount((count) => count + 1);
    setError(null);
    try {
      const media = await uploadPostMedia(file);
      if (collection === "thumbnail") {
        setThumbnail(media);
      } else {
        setGallery((items) =>
          items.some((item) => item.id === media.id) ? items : [...items, media],
        );
      }
      setIsDirty(true);
    } finally {
      setUploadingCount((count) => Math.max(0, count - 1));
      setIsUploading(false);
    }
  }, [validateFile]);

  const uploadThumbnail = useCallback(async (file: File) => {
    try {
      await upload(file, "thumbnail");
    } catch (uploadError) {
      const message = uploadError instanceof Error ? uploadError.message : "Không thể tải ảnh lên.";
      setError(message);
      notifyError(message);
    }
  }, [notifyError, upload]);

  const uploadGallery = useCallback(async (files: FileList | File[]) => {
    const fileArray = Array.from(files);
    for (const file of fileArray) {
      try {
        await upload(file, "gallery");
      } catch (uploadError) {
        const message = uploadError instanceof Error ? uploadError.message : "Không thể tải ảnh lên.";
        setError(message);
        notifyError(message);
      }
    }
  }, [notifyError, upload]);

  const removeThumbnail = useCallback(async () => {
    const current = thumbnail;
    setThumbnail(null);
    setIsDirty(true);
    if (current?.status === "temporary") {
      try {
        await deleteTemporaryMedia(current.id);
      } catch (removeError) {
        const message = removeError instanceof Error ? removeError.message : "Không thể xóa ảnh tạm thời.";
        setError(message);
        notifyError(message);
      }
    }
  }, [notifyError, thumbnail]);

  const removeGalleryItem = useCallback(async (mediaId: number) => {
    const current = gallery.find((item) => item.id === mediaId);
    setGallery((items) => items.filter((item) => item.id !== mediaId));
    setIsDirty(true);
    if (current?.status === "temporary") {
      try {
        await deleteTemporaryMedia(mediaId);
      } catch (removeError) {
        const message = removeError instanceof Error ? removeError.message : "Không thể xóa ảnh tạm thời.";
        setError(message);
        notifyError(message);
      }
    }
  }, [gallery, notifyError]);

  const reorderGallery = useCallback((fromIndex: number, toIndex: number) => {
    setGallery((items) => {
      if (fromIndex === toIndex || fromIndex < 0 || toIndex < 0 || fromIndex >= items.length || toIndex >= items.length) {
        return items;
      }
      const next = [...items];
      const [moved] = next.splice(fromIndex, 1);
      next.splice(toIndex, 0, moved);
      return next;
    });
    setIsDirty(true);
  }, []);

  const load = useCallback(async () => {
    if (!postId) return;
    setIsLoading(true);
    setError(null);
    try {
      const response = await getPostMedia(postId, undefined, entityPath);
      const thumbnailItem = response.data.find((item) => item.collection === "thumbnail");
      const galleryItems = response.data
        .filter((item) => item.collection === "gallery")
        .sort((a, b) => a.sortOrder - b.sortOrder);
      setThumbnail(thumbnailItem ? mediaFromItem(thumbnailItem) : null);
      setGallery(galleryItems.map(mediaFromItem));
      setIsDirty(false);
    } catch (loadError) {
      const message = loadError instanceof Error ? loadError.message : "Không thể tải hình ảnh bài viết.";
      setError(message);
      notifyError(message);
    } finally {
      setIsLoading(false);
    }
  }, [entityPath, notifyError, postId]);

  useEffect(() => {
    void load();
  }, [load]);

  const sync = useCallback(async (targetPostId: number) => {
    setIsSyncing(true);
    setError(null);
    try {
      await syncPostMedia(
        targetPostId,
        "thumbnail",
        thumbnail ? [{ id: thumbnail.id, sortOrder: 0 }] : [],
        entityPath,
      );
      await syncPostMedia(
        targetPostId,
        "gallery",
        // Backend validates sort_order as a contiguous zero-based index.
        gallery.map((item, index) => ({ id: item.id, sortOrder: index })),
        entityPath,
      );
      setIsDirty(false);
    } catch (syncError) {
      const message = syncError instanceof Error ? syncError.message : "Không thể đồng bộ hình ảnh.";
      setError(message);
      throw syncError;
    } finally {
      setIsSyncing(false);
    }
  }, [entityPath, gallery, thumbnail]);

  const markDirty = useCallback(() => setIsDirty(true), []);

  return useMemo(() => ({
    thumbnail,
    gallery,
    isUploading,
    uploadingCount,
    isLoading,
    isSyncing,
    error,
    isDirty,
    uploadThumbnail,
    uploadGallery,
    removeThumbnail,
    removeGalleryItem,
    reorderGallery,
    load,
    sync,
    markDirty,
  }), [
    thumbnail, gallery, isUploading, uploadingCount, isLoading, isSyncing, error, isDirty,
    uploadThumbnail, uploadGallery, removeThumbnail, removeGalleryItem, reorderGallery, load, sync, markDirty,
  ]);
};

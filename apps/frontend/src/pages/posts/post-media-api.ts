import { customRequest } from "../../providers/data";
import { API_URL, BACKEND_URL } from "../../providers/constants";
import {
  IPostMediaResponse,
  IPostMediaItem,
  IUploadedMedia,
  IMediaSyncItem,
  IUploadMediaResponse,
  PostMediaCollection,
} from "./post-media-types";

const getErrorMessage = async (response: Response, fallback: string): Promise<string> => {
  try {
    const body: unknown = await response.json();
    if (typeof body === "object" && body !== null && "error" in body) {
      const error = (body as { error?: unknown }).error;
      if (typeof error === "string" && error.length > 0) return error;
    }
  } catch {
    // Keep the fallback for non-JSON responses.
  }
  return fallback;
};

const unwrapData = <T,>(body: unknown): T => {
  if (typeof body === "object" && body !== null && "data" in body) {
    return (body as { data: T }).data;
  }
  return body as T;
};

const resolveMediaUrl = (value?: string | null): string | undefined => {
  if (!value) return undefined;
  try {
    return new URL(value, BACKEND_URL).toString();
  } catch {
    return value;
  }
};

const toUploadedMedia = (media: IUploadMediaResponse): IUploadedMedia => ({
  id: media.id,
  fileName: media.fileName,
  originalUrl: resolveMediaUrl(media.originalUrl) ?? "",
  thumbnailUrl: resolveMediaUrl(media.thumbnailUrl),
  mediumUrl: resolveMediaUrl(media.mediumUrl),
  status: media.status,
});

interface RawMedia {
  id: number;
  fileName?: string;
  file_name?: string;
  originalUrl?: string;
  original_url?: string;
  thumbnailUrl?: string;
  thumbnail_url?: string;
  mediumUrl?: string;
  medium_url?: string;
  status: IUploadedMedia["status"];
}

interface RawPostMediaItem {
  id: number;
  postId?: number;
  post_id?: number;
  mediaId?: number;
  media_id?: number;
  collection: PostMediaCollection;
  sortOrder?: number;
  sort_order?: number;
  media?: RawMedia;
}

const normalizeMedia = (media?: RawMedia): IUploadedMedia | undefined => {
  if (!media) return undefined;
  return toUploadedMedia({
    id: media.id,
    fileName: media.fileName ?? media.file_name ?? `media-${media.id}`,
    originalUrl: media.originalUrl ?? media.original_url ?? "",
    thumbnailUrl: media.thumbnailUrl ?? media.thumbnail_url,
    mediumUrl: media.mediumUrl ?? media.medium_url,
    status: media.status,
  });
};

const normalizePostMedia = (item: RawPostMediaItem): IPostMediaItem => ({
  id: item.id,
  postId: item.postId ?? item.post_id ?? 0,
  mediaId: item.mediaId ?? item.media_id ?? 0,
  collection: item.collection,
  sortOrder: item.sortOrder ?? item.sort_order ?? 0,
  media: normalizeMedia(item.media),
});

export const uploadPostMedia = async (file: File): Promise<IUploadedMedia> => {
  const formData = new FormData();
  formData.append("file", file);

  const response = await customRequest({ url: `${API_URL}/media/upload`, method: "POST", payload: formData });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Không thể tải ảnh lên."));
  }

  const body: unknown = await response.json();
  return toUploadedMedia(unwrapData<IUploadMediaResponse>(body));
};

export const getPostMedia = async (
  postId: number,
  collection?: PostMediaCollection,
): Promise<IPostMediaResponse> => {
  const query = collection ? `?collection=${encodeURIComponent(collection)}` : "";
  const response = await customRequest({ url: `${API_URL}/admin/posts/${postId}/media${query}` });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Không thể tải hình ảnh bài viết."));
  }

  const body: unknown = await response.json();
  const payload = unwrapData<IPostMediaResponse | IPostMediaResponse["data"]>(body);
  if (Array.isArray(payload)) {
    return { data: payload.map((item) => normalizePostMedia(item as RawPostMediaItem)), total: payload.length };
  }
  return {
    data: payload.data.map((item) => normalizePostMedia(item as RawPostMediaItem)),
    total: payload.total,
  };
};

export const syncPostMedia = async (
  postId: number,
  collection: PostMediaCollection,
  items: IMediaSyncItem[],
): Promise<void> => {
  const response = await customRequest({
    url: `${API_URL}/admin/posts/${postId}/media/${collection}`,
    method: "PUT",
    payload: { media: items.map((item) => ({ id: item.id, sort_order: item.sortOrder })) },
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Không thể đồng bộ hình ảnh bài viết."));
  }
};

export const deleteTemporaryMedia = async (mediaId: number): Promise<void> => {
  const response = await customRequest({ url: `${API_URL}/media/${mediaId}`, method: "DELETE" });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Không thể xóa ảnh tạm thời."));
  }
};

export type MediaStatus = "temporary" | "attached";
export type PostMediaCollection = "thumbnail" | "gallery";

export interface IUploadedMedia {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  status: MediaStatus;
}

export interface IPostMediaItem {
  id: number;
  postId: number;
  mediaId: number;
  collection: PostMediaCollection;
  sortOrder: number;
  media?: IUploadedMedia;
}

export interface IPostMediaResponse {
  data: IPostMediaItem[];
  total: number;
}

export interface IMediaSyncItem {
  id: number;
  sortOrder: number;
}

export interface IUploadMediaResponse {
  id: number;
  fileName: string;
  originalUrl: string;
  thumbnailUrl?: string;
  mediumUrl?: string;
  status: MediaStatus;
}

export interface PostMediaState {
  thumbnail: IUploadedMedia | null;
  gallery: IUploadedMedia[];
  isUploading: boolean;
  uploadingCount: number;
  isLoading: boolean;
  isSyncing: boolean;
  error: string | null;
  isDirty: boolean;
}

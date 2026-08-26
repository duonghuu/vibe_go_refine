export type SeoRobots =
  | "index,follow"
  | "noindex,follow"
  | "noindex,nofollow";

export type SeoOgType = "article" | "website" | "product";
export type SeoTwitterCard = "summary" | "summary_large_image";

export type JsonPrimitive = string | number | boolean | null;
export type JsonValue =
  | JsonPrimitive
  | JsonValue[]
  | { [key: string]: JsonValue };

export interface IPostSeoFormValues {
  metaTitle: string;
  metaDescription: string;
  metaKeywords: string;
  canonicalUrl: string;
  ogTitle: string;
  ogDescription: string;
  ogImage: string;
  ogType: SeoOgType | "";
  twitterTitle: string;
  twitterDescription: string;
  twitterImage: string;
  twitterCard: SeoTwitterCard | "";
  robots: SeoRobots;
  schemaJsonText: string;
}

export interface IPostSeoResponse {
  entityType: "post";
  entityId: number;
  metaTitle: string | null;
  metaDescription: string | null;
  metaKeywords: string | null;
  canonicalUrl: string | null;
  ogTitle: string | null;
  ogDescription: string | null;
  ogImage: string | null;
  ogType: SeoOgType | null;
  twitterTitle: string | null;
  twitterDescription: string | null;
  twitterImage: string | null;
  twitterCard: SeoTwitterCard | null;
  robots: SeoRobots;
  schemaJson: JsonValue | null;
  resolvedTitle: string;
  resolvedDescription: string;
  resolvedCanonicalUrl: string;
  resolvedOgTitle: string;
  resolvedOgDescription: string;
  resolvedOgImage: string | null;
  resolvedTwitterTitle: string;
  resolvedTwitterDescription: string;
  resolvedTwitterImage: string | null;
}

export interface IPostSeoPreview {
  title: string;
  description: string;
  canonicalUrl: string;
  ogImage: string | null;
  isTitleFallback: boolean;
  isDescriptionFallback: boolean;
  isCanonicalFallback: boolean;
  isOgImageFallback: boolean;
}

export interface IPostSeoContext {
  title: string;
  slug: string;
  content: string;
  thumbnailUrl?: string | null;
}

export interface PostSeoState {
  data: IPostSeoResponse | null;
  isLoading: boolean;
  isSaving: boolean;
  isError: boolean;
  error: Error | null;
  isDirty: boolean;
}

export const emptyPostSeoValues: IPostSeoFormValues = {
  metaTitle: "",
  metaDescription: "",
  metaKeywords: "",
  canonicalUrl: "",
  ogTitle: "",
  ogDescription: "",
  ogImage: "",
  ogType: "",
  twitterTitle: "",
  twitterDescription: "",
  twitterImage: "",
  twitterCard: "",
  robots: "index,follow",
  schemaJsonText: "",
};


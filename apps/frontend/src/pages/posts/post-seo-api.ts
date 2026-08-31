import {
  IPostSeoFormValues,
  IPostSeoResponse,
  JsonValue,
  SeoOgType,
  SeoRobots,
  SeoTwitterCard,
} from "./post-seo-types";

export interface RawPostSeo {
  entity_type?: string;
  entity_id?: number;
  entityType?: string;
  entityId?: number;
  meta_title?: string | null;
  metaTitle?: string | null;
  meta_description?: string | null;
  metaDescription?: string | null;
  meta_keywords?: string | null;
  metaKeywords?: string | null;
  canonical_url?: string | null;
  canonicalUrl?: string | null;
  og_title?: string | null;
  ogTitle?: string | null;
  og_description?: string | null;
  ogDescription?: string | null;
  og_image?: string | null;
  ogImage?: string | null;
  og_type?: string | null;
  ogType?: string | null;
  twitter_title?: string | null;
  twitterTitle?: string | null;
  twitter_description?: string | null;
  twitterDescription?: string | null;
  twitter_image?: string | null;
  twitterImage?: string | null;
  twitter_card?: string | null;
  twitterCard?: string | null;
  robots?: string | null;
  schema_json?: JsonValue | null;
  schemaJson?: JsonValue | null;
  resolved_title?: string;
  resolvedTitle?: string;
  resolved_description?: string;
  resolvedDescription?: string;
  resolved_canonical_url?: string;
  resolvedCanonicalUrl?: string;
  resolved_og_title?: string;
  resolvedOgTitle?: string;
  resolved_og_description?: string;
  resolvedOgDescription?: string;
  resolved_og_image?: string | null;
  resolvedOgImage?: string | null;
  resolved_twitter_title?: string;
  resolvedTwitterTitle?: string;
  resolved_twitter_description?: string;
  resolvedTwitterDescription?: string;
  resolved_twitter_image?: string | null;
  resolvedTwitterImage?: string | null;
}

const asSeoRobots = (value: string | null | undefined): SeoRobots => {
  if (value === "noindex,follow" || value === "noindex,nofollow") return value;
  return "index,follow";
};

const asOgType = (value: string | null | undefined): SeoOgType | null => {
  if (value === "article" || value === "website" || value === "product") return value;
  return null;
};

const asTwitterCard = (value: string | null | undefined): SeoTwitterCard | null => {
  if (value === "summary" || value === "summary_large_image") return value;
  return null;
};

export const normalizePostSeo = (raw: RawPostSeo, entityType: "post" | "page" = "post"): IPostSeoResponse => ({
  entityType,
  entityId: raw.entityId ?? raw.entity_id ?? 0,
  metaTitle: raw.metaTitle ?? raw.meta_title ?? null,
  metaDescription: raw.metaDescription ?? raw.meta_description ?? null,
  metaKeywords: raw.metaKeywords ?? raw.meta_keywords ?? null,
  canonicalUrl: raw.canonicalUrl ?? raw.canonical_url ?? null,
  ogTitle: raw.ogTitle ?? raw.og_title ?? null,
  ogDescription: raw.ogDescription ?? raw.og_description ?? null,
  ogImage: raw.ogImage ?? raw.og_image ?? null,
  ogType: asOgType(raw.ogType ?? raw.og_type),
  twitterTitle: raw.twitterTitle ?? raw.twitter_title ?? null,
  twitterDescription: raw.twitterDescription ?? raw.twitter_description ?? null,
  twitterImage: raw.twitterImage ?? raw.twitter_image ?? null,
  twitterCard: asTwitterCard(raw.twitterCard ?? raw.twitter_card),
  robots: asSeoRobots(raw.robots),
  schemaJson: raw.schemaJson ?? raw.schema_json ?? null,
  resolvedTitle: raw.resolvedTitle ?? raw.resolved_title ?? raw.metaTitle ?? raw.meta_title ?? "",
  resolvedDescription:
    raw.resolvedDescription ?? raw.resolved_description ?? raw.metaDescription ?? raw.meta_description ?? "",
  resolvedCanonicalUrl:
    raw.resolvedCanonicalUrl ?? raw.resolved_canonical_url ?? raw.canonicalUrl ?? raw.canonical_url ?? "",
  resolvedOgTitle: raw.resolvedOgTitle ?? raw.resolved_og_title ?? raw.ogTitle ?? raw.og_title ?? "",
  resolvedOgDescription:
    raw.resolvedOgDescription ?? raw.resolved_og_description ?? raw.ogDescription ?? raw.og_description ?? "",
  resolvedOgImage: raw.resolvedOgImage ?? raw.resolved_og_image ?? raw.ogImage ?? raw.og_image ?? null,
  resolvedTwitterTitle:
    raw.resolvedTwitterTitle ?? raw.resolved_twitter_title ?? raw.twitterTitle ?? raw.twitter_title ?? "",
  resolvedTwitterDescription:
    raw.resolvedTwitterDescription ?? raw.resolved_twitter_description ?? raw.twitterDescription ?? raw.twitter_description ?? "",
  resolvedTwitterImage:
    raw.resolvedTwitterImage ?? raw.resolved_twitter_image ?? raw.twitterImage ?? raw.twitter_image ?? null,
});

const toNullable = (value: string): string | null => {
  const normalized = value.trim();
  return normalized.length > 0 ? normalized : null;
};

const parseSchema = (value: string): JsonValue | null => {
  const normalized = value.trim();
  if (!normalized) return null;
  const parsed: unknown = JSON.parse(normalized);
  return parsed as JsonValue;
};

export const toPostSeoPayload = (values: IPostSeoFormValues) => {
  let schemaJson: JsonValue | null = null;
  try {
    schemaJson = parseSchema(values.schemaJsonText);
  } catch {
    throw new Error("Schema JSON-LD không hợp lệ.");
  }

  const payload = {
    meta_title: toNullable(values.metaTitle),
    meta_description: toNullable(values.metaDescription),
    meta_keywords: toNullable(values.metaKeywords),
    canonical_url: toNullable(values.canonicalUrl),
    og_title: toNullable(values.ogTitle),
    og_description: toNullable(values.ogDescription),
    og_image: toNullable(values.ogImage),
    og_type: toNullable(values.ogType),
    twitter_title: toNullable(values.twitterTitle),
    twitter_description: toNullable(values.twitterDescription),
    twitter_image: toNullable(values.twitterImage),
    twitter_card: toNullable(values.twitterCard),
    robots: values.robots,
    schema_json: schemaJson,
  };

  return payload;
};

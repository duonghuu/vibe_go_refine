import { getWebviewApiEndpoint, getWebviewPublicApiOrigin } from "@/lib/config";
import type { HeroContent, HeroSectionContent } from "@/types/home";
import type {
  PublicImageContentGroup,
  PublicImageContentItem,
  PublicImageContentResponse,
} from "@/types/image-content";

export type HeroImageState = "success" | "empty" | "error";

export interface HeroContentResult {
  content: HeroSectionContent;
  backgroundImage: string | null;
  state: HeroImageState;
}

const EMPTY_HERO_CONTENT: HeroSectionContent = {
  eyebrow: null,
  title: null,
  description: null,
  cta: null,
};

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function isOptionalString(value: unknown): value is string | null | undefined {
  return value === undefined || value === null || typeof value === "string";
}

function isPublicImageContentResponse(value: unknown): value is PublicImageContentResponse {
  if (!isRecord(value) || !Array.isArray(value.data)) {
    return false;
  }

  return value.data.every((group: unknown) => {
    if (
      !isRecord(group) ||
      typeof group.typeCode !== "string" ||
      typeof group.name !== "string" ||
      !Array.isArray(group.items)
    ) {
      return false;
    }

    return group.items.every((item: unknown) => {
      if (
        !isRecord(item) ||
        typeof item.id !== "number" ||
        typeof item.sortOrder !== "number" ||
        !isOptionalString(item.name) ||
        !isOptionalString(item.description) ||
        !isOptionalString(item.secondaryDescription) ||
        !isOptionalString(item.url) ||
        !isRecord(item.media)
      ) {
        return false;
      }

      return (
        typeof item.media.originalUrl === "string" &&
        typeof item.media.mimeType === "string" &&
        isOptionalString(item.media.thumbnailUrl) &&
        isOptionalString(item.media.mediumUrl)
      );
    });
  });
}

function resolveMediaUrl(value: string, publicBackendOrigin: string): string | null {
  if (value.startsWith("//")) {
    return null;
  }

  try {
    const url = new URL(value, publicBackendOrigin);
    return url.protocol === "http:" || url.protocol === "https:" ? url.toString() : null;
  } catch {
    return null;
  }
}

function getOptionalText(value: string | null | undefined): string | null {
  const text = value?.trim();
  return text || null;
}

function resolveHeroHref(value: string | null | undefined): string | null {
  const href = getOptionalText(value);
  if (!href || href.startsWith("//")) {
    return null;
  }

  if (href.startsWith("/")) {
    return href;
  }

  try {
    const url = new URL(href);
    return url.protocol === "http:" || url.protocol === "https:" ? url.toString() : null;
  } catch {
    return null;
  }
}

function selectHeroItem(
  group: PublicImageContentGroup,
  publicBackendOrigin: string,
): { item: PublicImageContentItem; backgroundImage: string } | undefined {
  return group.items
    .map((item) => ({
      item,
      backgroundImage: resolveMediaUrl(
        item.media.mediumUrl || item.media.originalUrl,
        publicBackendOrigin,
      ),
    }))
    .find((value): value is { item: PublicImageContentItem; backgroundImage: string } =>
      Boolean(value.backgroundImage),
    );
}

export async function getHeroContent(hero: HeroContent): Promise<HeroContentResult> {
  try {
    const endpoint = getWebviewApiEndpoint("/public/image-contents");
    endpoint.searchParams.set("typeCodes", "HERO");
    
    const response = await fetch(endpoint, { cache: "no-store" });
    if (!response.ok) {
      return { content: EMPTY_HERO_CONTENT, backgroundImage: null, state: "error" };
    }

    const payload: unknown = await response.json();
    if (!isPublicImageContentResponse(payload)) {
      return { content: EMPTY_HERO_CONTENT, backgroundImage: null, state: "error" };
    }

    const group = payload.data.find((value) => value.typeCode === "HERO");
    if (!group || group.items.length === 0) {
      return { content: EMPTY_HERO_CONTENT, backgroundImage: null, state: "empty" };
    }

    const publicBackendOrigin = getWebviewPublicApiOrigin();
    const selected = selectHeroItem(group, publicBackendOrigin);
    if (!selected) {
      return { content: EMPTY_HERO_CONTENT, backgroundImage: null, state: "empty" };
    }

    const href = resolveHeroHref(selected.item.url);
    const content: HeroSectionContent = {
      eyebrow: getOptionalText(selected.item.name),
      title: getOptionalText(selected.item.description),
      description: getOptionalText(selected.item.secondaryDescription),
      cta: href ? { label: hero.cta.label, href } : null,
    };

    return {
      content,
      backgroundImage: selected.backgroundImage,
      state: "success",
    };
  } catch (err) {
    return { content: EMPTY_HERO_CONTENT, backgroundImage: null, state: "error" };
  }
}

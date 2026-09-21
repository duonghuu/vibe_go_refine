import { getWebviewApiEndpoint, getWebviewPublicApiOrigin } from "@/lib/config";
import type { BenefitItem, IntroContent } from "@/types/home";
import type {
  PublicPageResponse,
  PublicPageSectionDto,
  PublicPageSectionItemDto,
  PublicPageSectionItemType,
} from "@/types/public-page";

const HOME_PAGE_SLUG = "trang-chu";
const INTRO_SECTION_KEY = "home_creative";
const INTRO_BENEFITS_COLLECTION = "postt";

export type IntroContentState = "success" | "empty" | "error";

export interface IntroContentResult {
  state: IntroContentState;
  content: IntroContent | null;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isFiniteNumber(value: unknown): value is number {
  return typeof value === "number" && Number.isFinite(value);
}

function isPublicPageSectionItemType(value: unknown): value is PublicPageSectionItemType {
  return value === "CATEGORY" || value === "POST" || value === "MEDIA";
}

function isPublicPageSectionItem(value: unknown): value is PublicPageSectionItemDto {
  if (!isRecord(value) || !isRecord(value.data)) {
    return false;
  }

  return (
    isPublicPageSectionItemType(value.itemType) &&
    isFiniteNumber(value.itemId) &&
    isFiniteNumber(value.sortOrder)
  );
}

function isPublicPageSectionCollections(
  value: unknown,
): value is Record<string, PublicPageSectionItemDto[]> {
  if (!isRecord(value)) {
    return false;
  }

  return Object.values(value).every(
    (items) => Array.isArray(items) && items.every(isPublicPageSectionItem),
  );
}

function isPublicPageSection(value: unknown): value is PublicPageSectionDto {
  return (
    isRecord(value) &&
    isFiniteNumber(value.id) &&
    typeof value.key === "string" &&
    typeof value.title === "string" &&
    typeof value.description === "string" &&
    isFiniteNumber(value.sortOrder) &&
    isPublicPageSectionCollections(value.collections)
  );
}

function isPublicPageResponse(value: unknown): value is PublicPageResponse {
  if (!isRecord(value) || !isRecord(value.data)) {
    return false;
  }

  return (
    isFiniteNumber(value.data.id) &&
    typeof value.data.slug === "string" &&
    Array.isArray(value.data.sections) &&
    value.data.sections.every(isPublicPageSection)
  );
}

function getIntroHeading(section: PublicPageSectionDto): IntroContent["heading"] | null {
  const eyebrow = section.title.trim();
  const title = section.description.trim();

  if (!eyebrow || !title) {
    return null;
  }

  return { eyebrow, title };
}

function getText(data: Record<string, unknown>, keys: string[]): string {
  for (const key of keys) {
    const value = data[key];
    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
  }

  return "";
}

function resolveMediaUrl(value: string, publicBackendOrigin: string): string | null {
  if (!value || value.startsWith("//")) {
    return null;
  }

  try {
    const url = new URL(value, publicBackendOrigin);
    return url.protocol === "http:" || url.protocol === "https:" ? url.toString() : null;
  } catch {
    return null;
  }
}

function mapBenefitItems(
  section: PublicPageSectionDto,
  publicBackendOrigin: string,
): BenefitItem[] {
  const items = section.collections[INTRO_BENEFITS_COLLECTION] ?? [];
  return items.flatMap((item) => {
    const title = getText(item.data, ["title", "name", "fileName"]);
    const description = getText(item.data, ["description", "secondaryDescription"]);
    const imageValue = getText(item.data, ["imageUrl", "thumbnailUrl", "mediumUrl", "originalUrl"]);
    const imageSrc = resolveMediaUrl(imageValue, publicBackendOrigin);
    if (!title || !imageSrc) {
      return [];
    }

    return [{
      image: { src: imageSrc, alt: title },
      title,
      description,
    }];
  });
}

export async function getIntroContent(): Promise<IntroContentResult> {
  try {
    const endpoint = getWebviewApiEndpoint(`/pages/${HOME_PAGE_SLUG}`);
    const response = await fetch(endpoint, { cache: "no-store" });

    if (response.status === 404) {
      return { state: "empty", content: null };
    }

    if (!response.ok) {
      return { state: "error", content: null };
    }

    const payload: unknown = await response.json();
    if (!isPublicPageResponse(payload) || payload.data.slug !== HOME_PAGE_SLUG) {
      return { state: "error", content: null };
    }

    const section = payload.data.sections.find((item) => item.key === INTRO_SECTION_KEY);
    if (!section) {
      return { state: "empty", content: null };
    }

    const heading = getIntroHeading(section);
    if (!heading) {
      return { state: "empty", content: null };
    }

    const benefits = mapBenefitItems(section, getWebviewPublicApiOrigin());
    if (benefits.length === 0) {
      return { state: "empty", content: null };
    }

    return { state: "success", content: { heading, benefits } };
  } catch {
    return { state: "error", content: null };
  }
}

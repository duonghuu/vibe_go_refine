import { useCallback, useEffect, useMemo, useState } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import { getPostSeo, savePostSeo } from "./post-seo-api";
import {
  emptyPostSeoValues,
  IPostSeoContext,
  IPostSeoFormValues,
  IPostSeoPreview,
  IPostSeoResponse,
  PostSeoState,
} from "./post-seo-types";

export interface UsePostSeoResult extends PostSeoState {
  form: UseFormReturn<IPostSeoFormValues>;
  preview: IPostSeoPreview;
  save: (targetPostId: number) => Promise<void>;
  retry: () => Promise<void>;
}

const stripHtml = (value: string): string =>
  value.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();

const truncate = (value: string, maxLength: number): string =>
  value.length > maxLength ? `${value.slice(0, maxLength - 1).trim()}…` : value;

const responseToFormValues = (data: IPostSeoResponse): IPostSeoFormValues => ({
  metaTitle: data.metaTitle ?? "",
  metaDescription: data.metaDescription ?? "",
  metaKeywords: data.metaKeywords ?? "",
  canonicalUrl: data.canonicalUrl ?? "",
  ogTitle: data.ogTitle ?? "",
  ogDescription: data.ogDescription ?? "",
  ogImage: data.ogImage ?? "",
  ogType: data.ogType ?? "",
  twitterTitle: data.twitterTitle ?? "",
  twitterDescription: data.twitterDescription ?? "",
  twitterImage: data.twitterImage ?? "",
  twitterCard: data.twitterCard ?? "",
  robots: data.robots,
  schemaJsonText: data.schemaJson ? JSON.stringify(data.schemaJson, null, 2) : "",
});

const getPreview = (
  values: IPostSeoFormValues,
  context: IPostSeoContext,
  data: IPostSeoResponse | null,
): IPostSeoPreview => {
  const fallbackDescription = truncate(stripHtml(context.content), 160);
  const fallbackCanonical = context.slug ? `/${context.slug}` : "";
  const title = values.metaTitle.trim() || data?.resolvedTitle || context.title;
  const description =
    values.metaDescription.trim() ||
    data?.resolvedDescription ||
    fallbackDescription;
  const canonicalUrl =
    values.canonicalUrl.trim() ||
    data?.resolvedCanonicalUrl ||
    fallbackCanonical;
  const ogImage = values.ogImage.trim() || data?.resolvedOgImage || context.thumbnailUrl || null;

  return {
    title,
    description,
    canonicalUrl,
    ogImage,
    isTitleFallback: !values.metaTitle.trim(),
    isDescriptionFallback: !values.metaDescription.trim(),
    isCanonicalFallback: !values.canonicalUrl.trim(),
    isOgImageFallback: !values.ogImage.trim(),
  };
};

export const usePostSeo = (
  postId: number | undefined,
  context: IPostSeoContext,
): UsePostSeoResult => {
  const form = useForm<IPostSeoFormValues>({
    defaultValues: emptyPostSeoValues,
    mode: "onBlur",
  });
  const [data, setData] = useState<IPostSeoResponse | null>(null);
  const [isLoading, setIsLoading] = useState(Boolean(postId));
  const [isSaving, setIsSaving] = useState(false);
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const values = form.watch();

  const load = useCallback(async () => {
    if (!postId) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setIsError(false);
    setError(null);

    try {
      const result = await getPostSeo(postId);
      setData(result);
      form.reset(result ? responseToFormValues(result) : emptyPostSeoValues);
    } catch (loadError) {
      const normalizedError =
        loadError instanceof Error ? loadError : new Error("Không thể tải thông tin SEO.");
      setIsError(true);
      setError(normalizedError);
    } finally {
      setIsLoading(false);
    }
  }, [form, postId]);

  useEffect(() => {
    void load();
  }, [load]);

  const save = useCallback(async (targetPostId: number) => {
    const valid = await form.trigger();
    if (!valid) {
      throw new Error("Vui lòng kiểm tra các trường SEO.");
    }

    setIsSaving(true);
    setIsError(false);
    setError(null);

    try {
      const result = await savePostSeo(targetPostId, form.getValues());
      setData(result);
      form.reset(responseToFormValues(result));
    } catch (saveError) {
      const normalizedError =
        saveError instanceof Error ? saveError : new Error("Không thể lưu thông tin SEO.");
      setIsError(true);
      setError(normalizedError);
      throw normalizedError;
    } finally {
      setIsSaving(false);
    }
  }, [form]);

  const preview = useMemo(
    () => getPreview(values, context, data),
    [context, data, values],
  );

  return {
    data,
    isLoading,
    isSaving,
    isError,
    error,
    isDirty: form.formState.isDirty,
    form,
    preview,
    save,
    retry: load,
  };
};

import { useCallback, useEffect, useMemo, useState } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import { HttpError, useNotification, useOne, useUpdate } from "@refinedev/core";
import { normalizePostSeo, RawPostSeo, toPostSeoPayload } from "./post-seo-api";
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
  const { open } = useNotification();
  const seoQuery = useOne<RawPostSeo, HttpError>({
    resource: "seo-meta",
    id: postId ? `post/${postId}` : undefined,
    queryOptions: { enabled: Boolean(postId), retry: false },
  });
  const seoMutation = useUpdate<RawPostSeo, HttpError, ReturnType<typeof toPostSeoPayload>>();
  const [data, setData] = useState<IPostSeoResponse | null>(null);
  const values = form.watch();

  useEffect(() => {
    if (seoQuery.result) {
      const normalized = normalizePostSeo(seoQuery.result);
      setData(normalized);
      form.reset(responseToFormValues(normalized));
    }
  }, [form, seoQuery.result]);

  useEffect(() => {
    if (seoQuery.query.error) {
      open?.({ type: "error", message: "Không thể tải thông tin SEO.", description: seoQuery.query.error.message });
    }
  }, [open, seoQuery.query.error]);

  const save = useCallback(async (targetPostId: number) => {
    const valid = await form.trigger();
    if (!valid) {
      throw new Error("Vui lòng kiểm tra các trường SEO.");
    }

    try {
      const result = await seoMutation.mutateAsync({
        resource: "seo-meta",
        id: `post/${targetPostId}`,
        values: toPostSeoPayload(form.getValues()),
      });
      const normalized = normalizePostSeo(result.data);
      setData(normalized);
      form.reset(responseToFormValues(normalized));
      open?.({ type: "success", message: "Đã lưu thông tin SEO." });
    } catch (saveError) {
      const normalizedError =
        saveError instanceof Error ? saveError : new Error("Không thể lưu thông tin SEO.");
      open?.({ type: "error", message: "Không thể lưu thông tin SEO.", description: normalizedError.message });
      throw normalizedError;
    }
  }, [form, open, seoMutation]);

  const preview = useMemo(
    () => getPreview(values, context, data),
    [context, data, values],
  );

  return {
    data,
    isLoading: Boolean(postId) && seoQuery.query.isLoading,
    isSaving: seoMutation.mutation.isPending,
    isError: Boolean(seoQuery.query.error),
    error: seoQuery.query.error ? new Error(seoQuery.query.error.message) : null,
    isDirty: form.formState.isDirty,
    form,
    preview,
    save,
    retry: async () => { await seoQuery.query.refetch(); },
  };
};

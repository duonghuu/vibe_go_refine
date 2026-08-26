import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest";
import { API_URL } from "./constants";

import { refreshTokenFn } from "./axiosInstance";

export const authenticatedFetch = async (url: string, options: RequestInit = {}): Promise<Response> => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    options.headers = {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    };
  }

  let response = await fetch(url, options);

  if (response.status === 401) {
    try {
      const newToken = await refreshTokenFn();
      if (newToken) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${newToken}`,
        };
        response = await fetch(url, options);
        if (response.status === 401) {
          throw new Error("Unauthorized after refresh");
        }
      }
    } catch (e) {
      // Failed to refresh, forcefully logout
      localStorage.removeItem("accessToken");
      localStorage.removeItem("user");
      window.location.href = "/login";
      throw e;
    }
  }
  return response;
};

/**
 * Gateway for custom Refine resources (uploads and nested post-media actions).
 * Keeping authentication and JSON/FormData handling here prevents components
 * and feature hooks from spreading raw fetch calls across the application.
 */
export const customRequest = async (params: {
  url: string;
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  payload?: unknown;
  headers?: HeadersInit;
}): Promise<Response> => {
  const isFormData = typeof FormData !== "undefined" && params.payload instanceof FormData;
  const headers = new Headers(params.headers);
  if (!isFormData && params.payload !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const body = isFormData || typeof params.payload === "string"
    ? params.payload as BodyInit | undefined
    : params.payload === undefined
      ? undefined
      : JSON.stringify(params.payload);

  return authenticatedFetch(params.url, {
    method: params.method ?? "GET",
    headers,
    body,
  });
};

const simpleRest = createSimpleRestDataProvider({
  apiURL: `${API_URL}/admin`,
  kyOptions: {
    fetch: authenticatedFetch as typeof fetch,
  },
});

export const kyInstance = simpleRest.kyInstance; // Kept for backward compatibility if used elsewhere


const parseError = (error: any) => {
  if (error && typeof error.message === "string") {
    try {
      const parsed = JSON.parse(error.message);
      if (parsed && parsed.error) {
        error.message = parsed.error;
      }
    } catch (e) {
      // ignore
    }
  }
};

export const dataProvider = {
  ...simpleRest.dataProvider,
  getList: async (params: any) => {
    try {
      const response = await simpleRest.dataProvider.getList(params);
      const body = response.data as any;
      if (body && !Array.isArray(body) && "data" in body) {
        return {
          data: body.data,
          total: body.total || 0,
        };
      }
      return response;
    } catch (error) {
      parseError(error);
      throw error;
    }
  },
  getOne: async (params: any) => {
    try {
      const response = await simpleRest.dataProvider.getOne(params);
      const body = response.data as any;
      if (body && !Array.isArray(body) && "data" in body) {
        return {
          data: body.data,
        };
      }
      return response;
    } catch (error) {
      parseError(error);
      throw error;
    }
  },
  create: async (params: any) => {
    try {
      const response = await simpleRest.dataProvider.create(params);
      const body = response.data as any;
      if (body && !Array.isArray(body) && "data" in body) {
        return { data: body.data };
      }
      return response;
    } catch (error) {
      parseError(error);
      throw error;
    }
  },
  update: async (params: any) => {
    try {
      const response = await simpleRest.dataProvider.update(params);
      const body = response.data as any;
      if (body && !Array.isArray(body) && "data" in body) {
        return { data: body.data };
      }
      return response;
    } catch (error) {
      parseError(error);
      throw error;
    }
  },
  deleteOne: async (params: any) => {
    try {
      const response = await simpleRest.dataProvider.deleteOne(params);
      const body = response.data as any;
      if (body && !Array.isArray(body) && "data" in body) {
        return { data: body.data };
      }
      return response;
    } catch (error) {
      parseError(error);
      throw error;
    }
  },
};

import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest";
import { API_URL } from "./constants";

import { refreshTokenFn } from "./axiosInstance";

const customHttpClient = async (url: string, options: RequestInit = {}) => {
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
      }
    } catch (e) {
      // Failed to refresh, will return the original 401 response and let AuthProvider.onError handle it
    }
  }
  return response as any;
};

const simpleRest = createSimpleRestDataProvider({
  apiURL: `${API_URL}/admin`,
  httpClient: customHttpClient,
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

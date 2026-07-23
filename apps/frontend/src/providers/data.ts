import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest";
import { API_URL } from "./constants";

const simpleRest = createSimpleRestDataProvider({
  apiURL: `${API_URL}/admin`,
});

export const kyInstance = simpleRest.kyInstance;

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

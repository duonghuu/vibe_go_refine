import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest";
import { API_URL } from "./constants";

const simpleRest = createSimpleRestDataProvider({
  apiURL: API_URL,
});

export const kyInstance = simpleRest.kyInstance;

export const dataProvider = {
  ...simpleRest.dataProvider,
  getList: async (params: any) => {
    const response = await simpleRest.dataProvider.getList(params);
    const body = response.data as any;
    if (body && !Array.isArray(body) && "data" in body) {
      return {
        data: body.data,
        total: body.total || 0,
      };
    }
    return response;
  },
  getOne: async (params: any) => {
    const response = await simpleRest.dataProvider.getOne(params);
    const body = response.data as any;
    if (body && !Array.isArray(body) && "data" in body) {
      return {
        data: body.data,
      };
    }
    return response;
  },
  create: async (params: any) => {
    const response = await simpleRest.dataProvider.create(params);
    const body = response.data as any;
    if (body && !Array.isArray(body) && "data" in body) {
      return { data: body.data };
    }
    return response;
  },
  update: async (params: any) => {
    const response = await simpleRest.dataProvider.update(params);
    const body = response.data as any;
    if (body && !Array.isArray(body) && "data" in body) {
      return { data: body.data };
    }
    return response;
  },
  deleteOne: async (params: any) => {
    const response = await simpleRest.dataProvider.deleteOne(params);
    const body = response.data as any;
    if (body && !Array.isArray(body) && "data" in body) {
      return { data: body.data };
    }
    return response;
  },
};

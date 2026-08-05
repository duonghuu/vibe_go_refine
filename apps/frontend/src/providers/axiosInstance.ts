import axios from "axios";
import { API_URL } from "./constants";

export const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
});

axiosInstance.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

// Flag to prevent multiple simultaneous refresh requests
let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise(function (resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers["Authorization"] = "Bearer " + token;
            return axiosInstance(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const response = await axios.post(
          `${API_URL}/auth/refresh-token`,
          {},
          { withCredentials: true }
        );
        const accessToken =
          response.data?.accessToken || response.data?.data?.accessToken;

        if (accessToken) {
          localStorage.setItem("accessToken", accessToken);
          axiosInstance.defaults.headers.common["Authorization"] =
            "Bearer " + accessToken;
          originalRequest.headers["Authorization"] = "Bearer " + accessToken;
          processQueue(null, accessToken);
          return axiosInstance(originalRequest);
        }
      } catch (refreshError) {
        processQueue(refreshError, null);
        localStorage.removeItem("accessToken");
        localStorage.removeItem("user");
        window.location.href = "/login";
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    return Promise.reject(error);
  }
);

export const refreshTokenFn = async () => {
  try {
    const response = await axios.post(
      `${API_URL}/auth/refresh-token`,
      {},
      { withCredentials: true }
    );
    const accessToken =
      response.data?.accessToken || response.data?.data?.accessToken;
    if (accessToken) {
      localStorage.setItem("accessToken", accessToken);
      return accessToken;
    }
    throw new Error("Missing accessToken in response");
  } catch (e) {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("user");
    window.location.href = "/login";
    throw e;
  }
};

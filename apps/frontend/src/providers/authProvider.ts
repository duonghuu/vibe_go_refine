import { AuthProvider } from "@refinedev/core";
import { axiosInstance } from "./axiosInstance";

export const authProvider: AuthProvider = {
  login: async ({ email, password }) => {
    try {
      const response = await axiosInstance.post(`/auth/login`, {
        email,
        password,
      });

      const data = response.data?.data || response.data;
      const { accessToken, user } = data;

      if (accessToken) {
        localStorage.setItem("accessToken", accessToken);
      }
      if (user) {
        localStorage.setItem("user", JSON.stringify(user));
      }

      return {
        success: true,
        redirectTo: "/",
      };
    } catch (error: any) {
      return {
        success: false,
        error: {
          name: "Lỗi đăng nhập",
          message:
            error.response?.data?.error ||
            error.response?.data?.message ||
            "Email hoặc mật khẩu không chính xác",
        },
      };
    }
  },
  logout: async () => {
    try {
      await axiosInstance.post(`/auth/logout`);
    } catch (error) {
      // ignore errors during logout (e.g., token already invalid)
    }
    localStorage.removeItem("accessToken");
    localStorage.removeItem("user");
    return {
      success: true,
      redirectTo: "/login",
    };
  },
  check: async () => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      return {
        authenticated: true,
      };
    }
    return {
      authenticated: false,
      redirectTo: "/login",
      logout: true,
    };
  },
  getPermissions: async () => null,
  getIdentity: async () => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch (e) {
        return null;
      }
    }
    return null;
  },
  onError: async (error) => {
    if (error.response?.status === 401 || error.status === 401) {
      return {
        logout: true,
        redirectTo: "/login",
      };
    }
    return {};
  },
};

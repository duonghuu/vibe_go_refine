const DEFAULT_WEBVIEW_API_URL = "http://localhost:8080/api/v1";
const DEFAULT_WEBVIEW_PUBLIC_API_URL = "http://localhost:8081";

function readWebviewApiUrl(): string {
  const configuredUrl = process.env.WEBVIEW_API_URL?.trim();
  if (!configuredUrl) {
    return DEFAULT_WEBVIEW_API_URL;
  }

  try {
    const url = new URL(configuredUrl);
    if (url.protocol !== "http:" && url.protocol !== "https:") {
      return DEFAULT_WEBVIEW_API_URL;
    }

    return configuredUrl.replace(/\/+$/, "");
  } catch {
    return DEFAULT_WEBVIEW_API_URL;
  }
}

function readWebviewPublicApiUrl(): string {
  const configuredUrl = process.env.WEBVIEW_PUBLIC_API_URL?.trim();

  if (!configuredUrl) {
    return DEFAULT_WEBVIEW_PUBLIC_API_URL;
  }

  try {
    const url = new URL(configuredUrl);
    if (url.protocol !== "http:" && url.protocol !== "https:") {
      return DEFAULT_WEBVIEW_PUBLIC_API_URL;
    }

    return configuredUrl.replace(/\/+$/, "");
  } catch {
    return DEFAULT_WEBVIEW_PUBLIC_API_URL;
  }
}

export const WEBVIEW_API_URL = readWebviewApiUrl();
export const WEBVIEW_PUBLIC_API_URL = readWebviewPublicApiUrl();

export function getWebviewApiEndpoint(pathname: string): URL {
  return new URL(pathname.replace(/^\/+/, ""), `${WEBVIEW_API_URL}/`);
}

export function getWebviewPublicApiOrigin(): string {
  return new URL(WEBVIEW_PUBLIC_API_URL).origin;
}

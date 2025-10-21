import type {
  FetchClient,
  ExtendedRequestInit
} from "@/core/entities/FetchClient";
import camelcaseKeys from "camelcase-keys";

export const API_BASE_URL = "http://localhost:8080/api/v1";

let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

export function createFetchClient(baseUrl: string): FetchClient {
  return async function <T>(
    path: string,
    options: ExtendedRequestInit = {}
  ): Promise<T> {
    let safeBody: BodyInit | null | undefined;

    if (
      options.body &&
      typeof options.body === "object" &&
      !(options.body instanceof URLSearchParams)
    ) {
      const params = new URLSearchParams();
      for (const [key, value] of Object.entries(options.body)) {
        params.append(key, String(value));
      }
      safeBody = params;
    } else {
      safeBody = options.body as BodyInit | null | undefined;
    }

    const safeOptions: RequestInit = {
      method: options.method,
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        ...options.headers
      },
      credentials: "include",
      body: safeBody
    };

    let res = await fetch(`${baseUrl}${path}`, safeOptions);

    if (res.status === 401) {
      try {
        await refreshToken();

        const retryOptions: RequestInit = {
          ...safeOptions,
          headers: {
            "Content-Type": "application/json",
            ...options.headers
          }
        };

        res = await fetch(`${baseUrl}${path}`, retryOptions);
      } catch (err) {
        throw new Error(`Unauthorized and refresh failed: ${err}`);
      }
    }

    const json = await res.json();

    if (!res.ok) {
      throw new Error(`${json.message}`);
    }

    return camelcaseKeys(json, { deep: true }) as T;
  };
}

async function refreshToken(): Promise<void> {
  if (!isRefreshing) {
    isRefreshing = true;
    refreshPromise = (async () => {
      const res = await fetch(`${API_BASE_URL}/accounts/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json"
        }
      });

      if (!res.ok) {
        throw new Error("Failed to refresh token");
      }

      isRefreshing = false;
    })();
  }
  return refreshPromise!;
}

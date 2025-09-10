import { FetchClient } from "@/core/entities/FetchClient";
import camelcaseKeys from "camelcase-keys";

export const API_BASE_URL = "http://localhost:8080/api/v1";

let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

export function createFetchClient(baseUrl: string): FetchClient {
  return async function <T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    let res = await fetch(`${baseUrl}${path}`, {
      ...options,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options.headers
      }
    });

    if (res.status === 401) {
      try {
        await refreshToken();

        res = await fetch(`${baseUrl}${path}`, {
          ...options,
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            ...options.headers
          }
        });
      } catch (err) {
        throw new Error(`Unauthorized and refresh failed: ${err}`);
      }
    }

    if (!res.ok) {
      throw new Error(`HTTP error ${res.status}: ${await res.text()}`);
    }

    const response = await res.json();
    return camelcaseKeys(response, { deep: true }) as T;
  };
}

async function refreshToken(): Promise<void> {
  if (!isRefreshing) {
    isRefreshing = true;
    refreshPromise = (async () => {
      const res = await fetch(`${API_BASE_URL}/google/auth/refresh`, {
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

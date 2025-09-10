import { createFetchClient } from "@/infrastructure/http/fetchClient";

describe("createFetchClient", () => {
  const mockFetch = vi.fn();

  beforeEach(() => {
    vi.resetAllMocks();
    global.fetch = mockFetch;
  });

  it("should call fetch with baseUrl + path", async () => {
    const client = createFetchClient("http://api.test");
    const mockResponse = {
      ok: true,
      json: () => Promise.resolve({ data: 123 })
    };
    mockFetch.mockResolvedValueOnce(mockResponse);

    const result = await client<{ data: number }>("/endpoint");

    expect(mockFetch).toHaveBeenCalledWith(
      "http://api.test/endpoint",
      expect.any(Object)
    );
    expect(result).toEqual({ data: 123 });
  });

  it("should include default headers and merge with custom headers", async () => {
    const client = createFetchClient("http://api.test");
    const mockResponse = { ok: true, json: () => Promise.resolve({}) };
    mockFetch.mockResolvedValueOnce(mockResponse);

    await client("/endpoint", {
      headers: { Authorization: "Bearer token" },
      method: "POST"
    });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://api.test/endpoint",
      expect.objectContaining({
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer token"
        }
      })
    );
  });

  it("should throw an error if response is not ok", async () => {
    const client = createFetchClient("http://api.test");
    const mockResponse = {
      ok: false,
      status: 500,
      text: () => Promise.resolve("Server error")
    };
    mockFetch.mockResolvedValueOnce(mockResponse);

    await expect(client("/endpoint")).rejects.toThrow(
      "HTTP error 500: Server error"
    );
  });

  it("should refresh token and retry request if 401", async () => {
    const client = createFetchClient("http://api.test");

    const unauthorizedResponse = {
      ok: false,
      status: 401,
      text: () => Promise.resolve("Unauthorized")
    };

    const refreshResponse = {
      ok: true,
      json: () => Promise.resolve({ refreshed: true })
    };

    const successResponse = {
      ok: true,
      json: () => Promise.resolve({ data: "after-refresh" })
    };

    mockFetch
      .mockResolvedValueOnce(unauthorizedResponse)
      .mockResolvedValueOnce(refreshResponse)
      .mockResolvedValueOnce(successResponse);

    const result = await client<{ data: string }>("/endpoint");

    expect(mockFetch).toHaveBeenNthCalledWith(
      1,
      "http://api.test/endpoint",
      expect.any(Object)
    );
    expect(mockFetch).toHaveBeenNthCalledWith(
      2,
      "http://localhost:8080/api/v1/google/auth/refresh",
      expect.any(Object)
    );
    expect(mockFetch).toHaveBeenNthCalledWith(
      3,
      "http://api.test/endpoint",
      expect.any(Object)
    );

    expect(result).toEqual({ data: "after-refresh" });
  });

  it("should throw if refresh fails", async () => {
    const client = createFetchClient("http://api.test");

    const unauthorizedResponse = {
      ok: false,
      status: 401,
      text: () => Promise.resolve("Unauthorized")
    };

    const failedRefreshResponse = {
      ok: false,
      status: 400,
      text: () => Promise.resolve("Bad refresh")
    };

    mockFetch
      .mockResolvedValueOnce(unauthorizedResponse)
      .mockResolvedValueOnce(failedRefreshResponse);

    await expect(client("/endpoint")).rejects.toThrow(
      /Unauthorized and refresh failed/
    );
  });
});

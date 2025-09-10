import { vi } from "vitest";
import { MovementsRepository } from "@/infrastructure/repositories/movements";
import {
  createFetchClient,
  API_BASE_URL
} from "@/infrastructure/http/fetchClient";
import { type MovementByYear } from "@/core/entities/Movement";

vi.mock("@/infrastructure/http/fetchClient", () => {
  return {
    API_BASE_URL: "http://fake-api",
    createFetchClient: vi.fn()
  };
});

describe("MovementsRepository", () => {
  const mockClient = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (createFetchClient as unknown as ReturnType<typeof vi.fn>).mockReturnValue(
      mockClient
    );
  });

  it("should call createFetchClient with API_BASE_URL", () => {
    new MovementsRepository();
    expect(createFetchClient).toHaveBeenCalledWith(API_BASE_URL);
  });

  it("should call client with the correct URL when getMovementsByYear is called", async () => {
    const repo = new MovementsRepository();
    const mockResponse: MovementByYear = {
      totalIncome: 100,
      totalOutcome: 40,
      balance: 60,
      months: []
    };

    mockClient.mockResolvedValueOnce(mockResponse);

    const result = await repo.getMovementsByYear(2023);

    expect(mockClient).toHaveBeenCalledWith("/movements/years/2023");

    expect(result).toEqual(mockResponse);
  });

  it("should propagate errors from client", async () => {
    const repo = new MovementsRepository();
    mockClient.mockRejectedValueOnce(new Error("Network error"));

    await expect(repo.getMovementsByYear(2023)).rejects.toThrow(
      "Network error"
    );
  });
});

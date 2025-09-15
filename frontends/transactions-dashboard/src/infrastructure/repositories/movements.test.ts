import { MovementsRepository } from "@/infrastructure/repositories/movements";
import {
  createFetchClient,
  API_BASE_URL
} from "@/infrastructure/http/fetchClient";
import type {
  MovementByYear,
  MovementsResponse,
  MovementByMonth
} from "@/core/entities/Movement";

vi.mock("@/infrastructure/http/fetchClient", () => ({
  API_BASE_URL: "http://fake-api",
  createFetchClient: vi.fn()
}));

describe("MovementsRepository", () => {
  const mockClient = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (createFetchClient as unknown as ReturnType<typeof vi.fn>).mockReturnValue(
      mockClient
    );
  });

  it("calls createFetchClient with API_BASE_URL on instantiation", () => {
    new MovementsRepository();
    expect(createFetchClient).toHaveBeenCalledWith(API_BASE_URL);
  });

  it("calls client with correct URL when getMovements is called", async () => {
    const repo = new MovementsRepository();
    const mockResponse: MovementsResponse = {
      page: 0,
      totalPages: 1,
      movements: []
    };
    mockClient.mockResolvedValueOnce(mockResponse);

    const result = await repo.getMovements(0);

    expect(mockClient).toHaveBeenCalledWith("/movements?page=0");
    expect(result).toEqual(mockResponse);
  });

  it("calls client with correct URL when getMovementsByYear is called", async () => {
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

  it("calls client with correct URL when getMovementsByMonth is called", async () => {
    const repo = new MovementsRepository();
    const mockResponse: MovementByMonth = {
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      days: []
    };
    mockClient.mockResolvedValueOnce(mockResponse);

    const result = await repo.getMovementsByMonth(2023, 0);

    expect(mockClient).toHaveBeenCalledWith("/movements/years/2023/months/1");
    expect(result).toEqual(mockResponse);
  });

  it("propagates errors from client for getMovementsByYear", async () => {
    const repo = new MovementsRepository();
    mockClient.mockRejectedValueOnce(new Error("Network error"));

    await expect(repo.getMovementsByYear(2023)).rejects.toThrow(
      "Network error"
    );
  });

  it("propagates errors from client for getMovementsByMonth", async () => {
    const repo = new MovementsRepository();
    mockClient.mockRejectedValueOnce(new Error("Network error"));

    await expect(repo.getMovementsByMonth(2023, 1)).rejects.toThrow(
      "Network error"
    );
  });
});

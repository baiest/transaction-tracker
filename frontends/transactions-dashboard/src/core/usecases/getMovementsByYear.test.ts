import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import type {
  IMovementsRepository,
  MovementByYear
} from "@/core/entities/Movement";

describe("GetMovementsByYear", () => {
  let mockRepository: IMovementsRepository;
  let usecase: GetMovementsByYear;

  beforeEach(() => {
    mockRepository = {
      getMovementsByYear: vi.fn(),
      getMovementsByMonth: vi.fn()
    };
    usecase = new GetMovementsByYear(mockRepository);
  });

  it("should call repository.getMovementsByYear with a single correct year and return the result", async () => {
    const mockData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: [{ income: 100, outcome: 50 }]
    };
    (
      mockRepository.getMovementsByYear as ReturnType<typeof vi.fn>
    ).mockResolvedValue(mockData);

    const result = await usecase.excecute([2023]);

    expect(mockRepository.getMovementsByYear).toHaveBeenCalledWith(2023);
    expect(result).toEqual([mockData]);
  });

  it("should call repository.getMovementsByYear for multiple years and return all results", async () => {
    const mockData2023: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: [{ income: 100, outcome: 50 }]
    };
    const mockData2024: MovementByYear = {
      totalIncome: 1200,
      totalOutcome: 600,
      balance: 600,
      months: [{ income: 120, outcome: 60 }]
    };
    (mockRepository.getMovementsByYear as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce(mockData2023)
      .mockResolvedValueOnce(mockData2024);

    const result = await usecase.excecute([2023, 2024]);

    expect(mockRepository.getMovementsByYear).toHaveBeenCalledWith(2023);
    expect(mockRepository.getMovementsByYear).toHaveBeenCalledWith(2024);
    expect(result).toEqual([mockData2023, mockData2024]);
  });

  it("should handle null or undefined results from the repository by returning default data", async () => {
    const mockData2023: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: [{ income: 100, outcome: 50 }]
    };
    (mockRepository.getMovementsByYear as ReturnType<typeof vi.fn>)
      .mockResolvedValueOnce(mockData2023)
      .mockResolvedValueOnce(null)
      .mockResolvedValueOnce(undefined);

    const result = await usecase.excecute([2023, 2024, 2025]);

    expect(result).toEqual([
      mockData2023,
      {
        balance: 0,
        totalIncome: 0,
        totalOutcome: 0,
        months: [{ income: 0, outcome: 0 }]
      },
      {
        balance: 0,
        totalIncome: 0,
        totalOutcome: 0,
        months: [{ income: 0, outcome: 0 }]
      }
    ]);
  });

  it("should throw an error for an invalid year (not a number)", async () => {
    await expect(usecase.excecute([Number("abc")])).rejects.toThrow(
      "year is not a number"
    );
  });

  it("should throw an error for an invalid year (out of range)", async () => {
    await expect(usecase.excecute([999])).rejects.toThrow("invalid year");
    await expect(usecase.excecute([10000])).rejects.toThrow("invalid year");
  });

  it("should propagate errors from the repository", async () => {
    (
      mockRepository.getMovementsByYear as ReturnType<typeof vi.fn>
    ).mockRejectedValue(new Error("DB connection failed"));
    await expect(usecase.excecute([2023])).rejects.toThrow(
      "DB connection failed"
    );
  });
});

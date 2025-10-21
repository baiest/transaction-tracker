import { GetMovementsByMonth } from "./getMovementsByMonth";
import type {
  IMovementsRepository,
  MovementByMonth
} from "@/core/entities/Movement";

describe("GetMovementsByMonth", () => {
  let mockRepository: IMovementsRepository;
  let usecase: GetMovementsByMonth;

  beforeEach(() => {
    mockRepository = {
      createMovement: vi.fn(),
      getMovements: vi.fn(),
      getMovementsByYear: vi.fn(),
      getMovementsByMonth: vi.fn()
    };
    usecase = new GetMovementsByMonth(mockRepository);
  });

  it("should call repository.getMovementsByMonth with the correct year and month", async () => {
    const mockData: MovementByMonth = {
      totalIncome: 100,
      totalExpense: 50,
      balance: 50,
      days: [{ day: 15, income: 100, outcome: 50 }]
    };

    (
      mockRepository.getMovementsByMonth as ReturnType<typeof vi.fn>
    ).mockResolvedValue(mockData);

    const year = 2023;
    const month = 1;

    const result = await usecase.excecute(year, month);

    expect(mockRepository.getMovementsByMonth).toHaveBeenCalledWith(
      year,
      month
    );
    expect(result).toEqual({
      totalIncome: 100,
      totalExpense: 50,
      balance: 50,
      days: [
        ...Array(14).fill({ day: 0, income: 0, outcome: 0 }),
        { day: 15, income: 100, outcome: 50 },
        ...Array(16).fill({ day: 0, income: 0, outcome: 0 })
      ]
    });
  });

  it("should handle multiple entries in the days array", async () => {
    const mockData: MovementByMonth = {
      totalIncome: 200,
      totalExpense: 75,
      balance: 125,
      days: [
        { day: 5, income: 50, outcome: 10 },
        { day: 20, income: 150, outcome: 65 }
      ]
    };

    (
      mockRepository.getMovementsByMonth as ReturnType<typeof vi.fn>
    ).mockResolvedValue(mockData);

    const result = await usecase.excecute(2023, 1);

    const expectedDays = [
      ...Array(4).fill({ day: 0, income: 0, outcome: 0 }),
      { day: 5, income: 50, outcome: 10 },
      ...Array(14).fill({ day: 0, income: 0, outcome: 0 }),
      { day: 20, income: 150, outcome: 65 },
      ...Array(11).fill({ day: 0, income: 0, outcome: 0 })
    ];

    expect(result.days).toEqual(expectedDays);
  });

  it("should handle an empty days array from the repository", async () => {
    const mockData: MovementByMonth = {
      totalIncome: 0,
      totalExpense: 0,
      balance: 0,
      days: []
    };
    (
      mockRepository.getMovementsByMonth as ReturnType<typeof vi.fn>
    ).mockResolvedValue(mockData);

    const result = await usecase.excecute(2023, 1);

    expect(result.totalIncome).toBe(0);
    expect(result.days.length).toBe(31);
    expect(
      result.days.every((day) => day.income === 0 && day.outcome === 0)
    ).toBe(true);
  });

  it("should throw an error for an invalid year (not a number)", async () => {
    await expect(usecase.excecute(Number("abc"), 1)).rejects.toThrow(
      "year is not a number"
    );
  });

  it("should throw an error for an invalid year (out of range)", async () => {
    await expect(usecase.excecute(999, 1)).rejects.toThrow("invalid year");
    await expect(usecase.excecute(10000, 1)).rejects.toThrow("invalid year");
  });

  it("should throw an error for an invalid month (out of range)", async () => {
    await expect(usecase.excecute(2023, -1)).rejects.toThrow("invalid month");
    await expect(usecase.excecute(2023, 12)).rejects.toThrow("invalid month");
  });

  it("should propagate errors from the repository", async () => {
    (
      mockRepository.getMovementsByMonth as ReturnType<typeof vi.fn>
    ).mockRejectedValue(new Error("Database connection failed"));
    await expect(usecase.excecute(2023, 1)).rejects.toThrow(
      "Database connection failed"
    );
  });
});

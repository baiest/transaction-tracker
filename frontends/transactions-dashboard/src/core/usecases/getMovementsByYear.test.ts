import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import type {
  IMovementsRepository,
  MovementByYear
} from "@/core/entities/Movement";

describe("GetMovementsByYear", () => {
  it("should call repository.getMovementsByYear with the correct year", async () => {
    const mockRepository: IMovementsRepository = {
      getMovementsByYear: vi.fn().mockResolvedValue({
        totalIncome: 100,
        totalOutcome: 50,
        balance: 50,
        months: []
      } as MovementByYear)
    };

    const usecase = new GetMovementsByYear(mockRepository);

    const result = await usecase.excecute(2023);

    expect(mockRepository.getMovementsByYear).toHaveBeenCalledWith(2023);
    expect(result).toEqual({
      totalIncome: 100,
      totalOutcome: 50,
      balance: 50,
      months: []
    });
  });

  it("should propagate errors from repository", async () => {
    const mockRepository: IMovementsRepository = {
      getMovementsByYear: vi.fn().mockRejectedValue(new Error("DB error"))
    };

    const usecase = new GetMovementsByYear(mockRepository);

    await expect(usecase.excecute(2023)).rejects.toThrow("DB error");
  });
});

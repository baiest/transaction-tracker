import { useMovementsStore } from "@/infrastructure/store/movements";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import { MovementByYear } from "@/core/entities/Movement";

vi.mock("@/core/usecases/getMovementsByYear", () => {
  return {
    GetMovementsByYear: vi.fn().mockImplementation(() => ({
      excecute: vi.fn()
    }))
  };
});

describe("useMovementsStore", () => {
  let mockExecute: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    useMovementsStore.setState({
      movementsByYear: {
        totalIncome: 0,
        totalOutcome: 0,
        balance: 0,
        months: []
      },
      year: 0,
      isLoading: false,
      error: null,
      timeSelected: "year",

      fetchMomentsByYear: useMovementsStore.getState().fetchMomentsByYear
    });

    mockExecute = (GetMovementsByYear as ReturnType<typeof vi.fn>).mock
      .results[0].value.excecute;
  });

  it("should have an initial state", () => {
    const state = useMovementsStore.getState();

    expect(state.movementsByYear).toEqual({
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      months: []
    });
    expect(state.allYearsRaw).toEqual([]);
    expect(state.timeSelected).toBeTruthy();
    expect(state.year).toBe(0);
    expect(state.isLoading).toBeFalsy();
    expect(state.error).toBeNull();
  });

  it("should fetch movements successfully", async () => {
    const fakeData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: []
    };

    mockExecute.mockResolvedValueOnce([fakeData]);

    await useMovementsStore.getState().fetchMomentsByYear(2024);

    const store = useMovementsStore.getState();

    expect(mockExecute).toHaveBeenCalledWith([2024]);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movementsByYear).toEqual(fakeData);
    expect(store.year).toBe(2024);
  });

  it("should fetch all years movements successfully", async () => {
    const fakeData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: []
    };

    mockExecute.mockResolvedValue([fakeData, fakeData, fakeData]);

    await useMovementsStore.getState().fetchAllYearsData([2024, 2023, 2022]);

    const store = useMovementsStore.getState();

    expect(mockExecute).toHaveBeenCalledWith([2024, 2023, 2022]);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.allYearsRaw).length(3);
    expect(store.allYearsRaw).toEqual([fakeData, fakeData, fakeData]);
  });

  it("should handle errors correctly", async () => {
    mockExecute.mockRejectedValueOnce(new Error("Failed to load"));

    await useMovementsStore.getState().fetchMomentsByYear(2025);

    let store = useMovementsStore.getState();

    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load");

    mockExecute.mockRejectedValueOnce(new Error("Failed to load"));

    await useMovementsStore.getState().fetchAllYearsData([2025]);

    store = useMovementsStore.getState();

    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load");
  });

  it("should update the year when setYear is called", () => {
    const { setYear } = useMovementsStore.getState();

    setYear(2030);

    const store = useMovementsStore.getState();
    expect(store.year).toBe(2030);
  });
});

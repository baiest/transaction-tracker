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
      fetchMomentesByYear: useMovementsStore.getState().fetchMomentesByYear
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
    expect(state.year).toBe(0);
    expect(state.isLoading).toBe(false);
    expect(state.error).toBeNull();
  });

  it("should have the correct initial state", () => {
    const store = useMovementsStore.getState();

    expect(store.year).toBe(0);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movementsByYear).toEqual({
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      months: []
    });
  });

  it("should fetch movements successfully", async () => {
    const fakeData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: []
    };

    mockExecute.mockResolvedValueOnce(fakeData);

    await useMovementsStore.getState().fetchMomentesByYear(2024);

    const store = useMovementsStore.getState();

    expect(mockExecute).toHaveBeenCalledWith(2024);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movementsByYear).toEqual(fakeData);
    expect(store.year).toBe(2024);
  });

  it("should handle errors correctly", async () => {
    mockExecute.mockRejectedValueOnce(new Error("Failed to load"));

    await useMovementsStore.getState().fetchMomentesByYear(2025);

    const store = useMovementsStore.getState();

    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load");
  });
});

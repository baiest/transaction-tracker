import { useMovementsStore } from "@/infrastructure/store/movements";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import type {
  MovementByYear,
  MovementsResponse
} from "@/core/entities/Movement";
import { GetMovementsByMonth } from "@/core/usecases/getMovementsByMonth";
import { GetMovements } from "@/core/usecases/getMovements";

vi.mock("@/core/usecases/getMovementsByYear", () => {
  return {
    GetMovementsByYear: vi.fn().mockImplementation(() => ({
      excecute: vi.fn()
    }))
  };
});

vi.mock("@/core/usecases/getMovementsByMonth", () => {
  return {
    GetMovementsByMonth: vi.fn().mockImplementation(() => ({
      excecute: vi.fn()
    }))
  };
});

vi.mock("@/core/usecases/getMovements", () => {
  return {
    GetMovements: vi.fn().mockImplementation(() => ({
      excecute: vi.fn(),
      totalPages: 0
    }))
  };
});

describe("useMovementsStore", () => {
  let mockExecuteYear: ReturnType<typeof vi.fn>;
  let mockExecuteMonth: ReturnType<typeof vi.fn>;
  let mockExecuteMovements: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockExecuteYear = (GetMovementsByYear as ReturnType<typeof vi.fn>).mock
      .results[0].value.excecute;
    mockExecuteMonth = (GetMovementsByMonth as ReturnType<typeof vi.fn>).mock
      .results[0].value.excecute;
    mockExecuteMovements = (GetMovements as ReturnType<typeof vi.fn>).mock
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
    expect(state.movementsByMonth).toEqual({
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      days: []
    });
    expect(state.allYearsRaw).toEqual([]);
    expect(state.movements).toEqual([]);
    expect(state.totalPages).toBe(0);
    expect(state.timeSelected).toBe("year");
    expect(state.year).toBe(new Date().getFullYear());
    expect(state.month).toBe(new Date().getMonth() - 1);
    expect(state.isLoading).toBeFalsy();
    expect(state.error).toBeNull();
  });

  it("should fetch movements by year successfully", async () => {
    const fakeData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: []
    };

    mockExecuteYear.mockResolvedValueOnce([fakeData]);

    await useMovementsStore.getState().fetchMomentsByYear(2024);

    const store = useMovementsStore.getState();

    expect(mockExecuteYear).toHaveBeenCalledWith([2024]);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movementsByYear).toEqual(fakeData);
    expect(store.year).toBe(2024);
  });

  it("should fetch movements by month successfully", async () => {
    const fakeData = {
      totalIncome: 200,
      totalOutcome: 50,
      balance: 150,
      days: []
    };

    mockExecuteMonth.mockResolvedValueOnce(fakeData);

    await useMovementsStore.getState().fetchMomentsByMonth(2024, 1);

    const store = useMovementsStore.getState();

    expect(mockExecuteMonth).toHaveBeenCalledWith(2024, 1);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movementsByMonth).toEqual(fakeData);
  });

  it("should fetch movements successfully", async () => {
    const fakeData = [
      { id: "1", amount: 10, type: "income" },
      { id: "2", amount: 20, type: "outcome" }
    ];
    const getMovementsInstance = (GetMovements as ReturnType<typeof vi.fn>).mock
      .results[0].value;
    (getMovementsInstance as MovementsResponse).totalPages = 5;

    mockExecuteMovements.mockResolvedValueOnce(fakeData);

    await useMovementsStore.getState().fetchMovements(1);

    const store = useMovementsStore.getState();

    expect(mockExecuteMovements).toHaveBeenCalledWith(1);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.movements).toEqual(fakeData);
    expect(store.totalPages).toBe(5);
  });

  it("should fetch all years movements successfully", async () => {
    const fakeData: MovementByYear = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: []
    };

    mockExecuteYear.mockResolvedValue([fakeData, fakeData, fakeData]);

    await useMovementsStore.getState().fetchAllYearsData([2024, 2023, 2022]);

    const store = useMovementsStore.getState();

    expect(mockExecuteYear).toHaveBeenCalledWith([2024, 2023, 2022]);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.allYearsRaw.length).toBe(3);
    expect(store.allYearsRaw).toEqual([fakeData, fakeData, fakeData]);
  });

  it("should handle errors correctly for all fetch methods", async () => {
    mockExecuteYear.mockRejectedValueOnce(new Error("Failed to load year"));
    await useMovementsStore.getState().fetchMomentsByYear(2025);
    let store = useMovementsStore.getState();
    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load year");

    mockExecuteMonth.mockRejectedValueOnce(new Error("Failed to load month"));
    await useMovementsStore.getState().fetchMomentsByMonth(2025, 1);
    store = useMovementsStore.getState();
    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load month");

    mockExecuteMovements.mockRejectedValueOnce(
      new Error("Failed to load movements")
    );
    await useMovementsStore.getState().fetchMovements(1);
    store = useMovementsStore.getState();
    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load movements");

    mockExecuteYear.mockRejectedValueOnce(
      new Error("Failed to load all years")
    );
    await useMovementsStore.getState().fetchAllYearsData([2025]);
    store = useMovementsStore.getState();
    expect(store.isLoading).toBe(false);
    expect(store.error).toBe("Failed to load all years");
  });

  it("should update the year when setYear is called", () => {
    const { setYear } = useMovementsStore.getState();
    setYear(2030);
    const store = useMovementsStore.getState();
    expect(store.year).toBe(2030);
  });

  it("should update the month when setMonth is called", () => {
    const { setMonth } = useMovementsStore.getState();
    setMonth(5);
    const store = useMovementsStore.getState();
    expect(store.month).toBe(5);
  });

  it("should update the timeSelected when setTimeSelected is called", () => {
    const { setTimeSelected } = useMovementsStore.getState();
    setTimeSelected("month");
    const store = useMovementsStore.getState();
    expect(store.timeSelected).toBe("month");
  });
});

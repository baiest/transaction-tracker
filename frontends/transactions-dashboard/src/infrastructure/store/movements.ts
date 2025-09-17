import { GetMovementsByMonth } from "@/core/usecases/getMovementsByMonth";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import { MovementsRepository } from "@/infrastructure/repositories/movements";
import { create } from "zustand";
import { MovementsStore, Time } from "./models";
import { GetMovements } from "@/core/usecases/getMovements";

export const useMovementsStore = create<MovementsStore>((set) => {
  const movementsRepository = new MovementsRepository();
  const getMovementsByYear = new GetMovementsByYear(movementsRepository);
  const getMovementsByMonth = new GetMovementsByMonth(movementsRepository);
  const getMovements = new GetMovements(movementsRepository);

  return {
    totalPages: 0,
    movements: [],
    movementsByYear: {
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      months: []
    },
    movementsByMonth: {
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      days: []
    },
    allYearsRaw: [],
    month: new Date().getMonth() - 1,
    year: new Date().getFullYear(),
    timeSelected: "year",
    isLoading: false,
    error: null,

    setYear: (year: number) => set({ year }),
    setMonth: (month: number) => set({ month }),
    setTimeSelected: (timeSelected: Time) => set({ timeSelected }),
    fetchMomentsByYear: async (year: number): Promise<void> => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute([year]);
        set({ movementsByYear: data[0], isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchMomentsByMonth: async (year: number, month: number): Promise<void> => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByMonth.excecute(year, month);
        set({ movementsByMonth: data, isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchAllYearsData: async (years: number[]) => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute(years);

        set({
          allYearsRaw: data,
          isLoading: false
        });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchMovements: async (page: number) => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovements.excecute(page);

        set({
          movements: data,
          totalPages: getMovements.totalPages,
          isLoading: false
        });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    }
  };
});

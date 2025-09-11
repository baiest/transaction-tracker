import type { MovementYear, MovementByYear } from "@/core/entities/Movement";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import { MovementsRepository } from "@/infrastructure/repositories/movements";
import { create } from "zustand";

export interface MovementsStore {
  movementsByYear: MovementByYear;
  allYearsRaw: MovementByYear[];
  year: number;
  showAllYears: boolean;
  isLoading: boolean;
  error: string | null;

  setYear: (year: number) => void;
  setShowAllYears: (value: boolean) => void;
  fetchMomentsByYear: (year: number) => Promise<void>;
  fetchAllYearsData: (years: number[]) => Promise<void>;
}

export const useMovementsStore = create<MovementsStore>((set, get) => {
  const movementsRepository = new MovementsRepository();
  const getMovementsByYear = new GetMovementsByYear(movementsRepository);

  return {
    movementsByYear: {
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      months: []
    },
    allYearsRaw: [],
    year: new Date().getFullYear(),
    showAllYears: false,
    isLoading: false,
    error: null,

    setYear: (year: number) => set({ year }),
    setShowAllYears: (value: boolean) => set({ showAllYears: value }),
    fetchMomentsByYear: async (year: number): Promise<void> => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute(year);
        set({ movementsByYear: data, isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchAllYearsData: async (years: number[]) => {
      set({ isLoading: true, error: null });
      try {
        const data = await Promise.all(
          years.map((y) => getMovementsByYear.excecute(y))
        );

        const validData = data.map((d: MovementByYear | null | undefined) =>
          !!d
            ? d
            : {
                balance: 0,
                totalIncome: 0,
                totalOutcome: 0,
                months: [{ income: 0, outcome: 0 }]
              }
        );

        set({
          allYearsRaw: validData,
          isLoading: false
        });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    }
  };
});

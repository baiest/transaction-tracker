import { type MovementByYear } from "@/core/entities/Movement";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import { MovementsRepository } from "@/infrastructure/repositories/movements";
import { create } from "zustand";

export interface MovementsStore {
  movementsByYear: MovementByYear;
  year: number;
  isLoading: boolean;
  error: string | null;

  fetchMomentesByYear: (year: number) => Promise<void>;
}

export const useMovementsStore = create<MovementsStore>((set) => {
  const movementsRepository = new MovementsRepository();
  const getMovementsByYear = new GetMovementsByYear(movementsRepository);

  return {
    movementsByYear: {
      totalIncome: 0,
      totalOutcome: 0,
      balance: 0,
      months: []
    },
    year: 0,
    isLoading: false,
    error: null,

    fetchMomentesByYear: async (year: number) => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute(year);
        set({ movementsByYear: data, isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    }
  };
});

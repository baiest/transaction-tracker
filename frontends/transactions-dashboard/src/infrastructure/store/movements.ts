import { GetMovementsByMonth } from "@/core/usecases/getMovementsByMonth";
import { GetMovementsByYear } from "@/core/usecases/getMovementsByYear";
import { MovementsRepository } from "@/infrastructure/repositories/movements";
import { create } from "zustand";
import { MovementsStore, Time } from "./models";
import { GetMovements } from "@/core/usecases/getMovements";
import { MovementRequest } from "@/core/entities/Movement";
import { CreateMovement } from "@/core/usecases/createMovement";

export const useMovementsStore = create<MovementsStore>((set, get) => {
  const movementsRepository = new MovementsRepository();
  const getMovementsByYear = new GetMovementsByYear(movementsRepository);
  const getMovementsByMonth = new GetMovementsByMonth(movementsRepository);
  const createMovement = new CreateMovement(movementsRepository);
  const getMovements = new GetMovements(movementsRepository);

  return {
    totalPages: 0,
    movements: [],
    institutionsSelected: [],
    movementsByYear: {
      totalIncome: 0,
      totalExpense: 0,
      balance: 0,
      months: []
    },
    movementsByMonth: {
      totalIncome: 0,
      totalExpense: 0,
      balance: 0,
      days: []
    },
    allYearsRaw: [],
    month: new Date().getMonth() - 1,
    year: new Date().getFullYear(),
    timeSelected: "year",
    isLoading: false,
    error: null,

    setInstitutionsSelected: (institutionsSelected: string[]) => {
      set({ institutionsSelected: [...new Set(institutionsSelected)] });
    },
    setYear: (year: number) => set({ year }),
    setMonth: (month: number) => set({ month }),
    setTimeSelected: (timeSelected: Time) => set({ timeSelected }),
    fetchMomentsByYear: async (year: number): Promise<void> => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute(
          [year],
          get().institutionsSelected
        );

        set({ movementsByYear: data[0], isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchMomentsByMonth: async (year: number, month: number): Promise<void> => {
      set({ isLoading: true, error: null });

      try {
        const data = await getMovementsByMonth.excecute(
          year,
          month,
          get().institutionsSelected
        );
        set({ movementsByMonth: data, isLoading: false, year });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    fetchAllYearsData: async (years: number[]) => {
      set({ isLoading: true, error: null });
      try {
        const data = await getMovementsByYear.excecute(
          years,
          get().institutionsSelected
        );

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
        const data = await getMovements.excecute(
          page,
          get().institutionsSelected
        );

        set({
          movements: data,
          totalPages: getMovements.totalPages,
          isLoading: false
        });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    },

    createMovement: async (movement: MovementRequest) => {
      set({ isLoading: true, error: null });

      try {
        await createMovement.excecute(movement);
        set({
          isLoading: false
        });
      } catch (err: unknown) {
        set({ error: (err as Error).message, isLoading: false });
      }
    }
  };
});

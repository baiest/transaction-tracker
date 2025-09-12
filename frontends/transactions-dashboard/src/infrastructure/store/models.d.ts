import type { MovementByMonth, MovementByYear } from "@/core/entities/Movement";

export type Time = "all_years" | "year" | "month";

export interface MovementsStore {
  movementsByYear: MovementByYear;
  movementsByMonth: MovementByMonth;
  allYearsRaw: MovementByYear[];
  year: number;
  month: number;
  timeSelected: Time;
  isLoading: boolean;
  error: string | null;

  setYear: (year: number) => void;
  setMonth: (month: number) => void;
  setTimeSelected: (value: Time) => void;
  fetchMomentsByYear: (year: number) => Promise<void>;
  fetchMomentsByMonth: (year: number, month: number) => Promise<void>;
  fetchAllYearsData: (years: number[]) => Promise<void>;
}

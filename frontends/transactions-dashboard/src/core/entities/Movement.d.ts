export interface MovementYear {
  income: number;
  outcome: number;
}

export interface MovementMonth {
  day: number;
  income: number;
  outcome: number;
}

export interface MovementByYear {
  totalIncome: number;
  totalOutcome: number;
  balance: number;
  months: MovementYear[];
}

export interface MovementByMonth {
  totalIncome: number;
  totalOutcome: number;
  balance: number;
  days: MovementMonth[];
}

export interface IMovementsRepository {
  getMovementsByYear: (year: number) => Promise<MovementByYear>;
  getMovementsByMonth: (
    year: number,
    month: number
  ) => Promise<MovementByMonth>;
}

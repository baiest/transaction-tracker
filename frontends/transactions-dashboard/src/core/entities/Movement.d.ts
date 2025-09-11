export interface MovementYear {
  income: number;
  outcome: number;
}

export interface MovementByYear {
  totalIncome: number;
  totalOutcome: number;
  balance: number;
  months: MovementYear[];
}

export interface IMovementsRepository {
  getMovementsByYear: (year: number) => Promise<MovementByYear>;
}

export interface Movement {
  id: string;
  email: string;
  messageId: string;
  date: Date;
  value: number;
  isNegative: boolean;
  topic: strting;
  detail: string;
}

export interface MovementResponse extends Omit<Movement, "date"> {
  date: string;
}

export interface MovementsResponse {
  totalPages: number;
  page: number;
  movements: MovementResponse[];
}

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
  getMovements: (page: number) => Promise<MovementsResponse>;
  getMovementsByYear: (year: number) => Promise<MovementByYear>;
  getMovementsByMonth: (
    year: number,
    month: number
  ) => Promise<MovementByMonth>;
}

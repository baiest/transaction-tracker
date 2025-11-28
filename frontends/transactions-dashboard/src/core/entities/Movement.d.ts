export interface Movement {
  id: string;
  email: string;
  messageId: string;
  date: Date;
  amount: number;
  type: string;
  category: string;
  topic: strting;
  description: string;
}

export interface MovementResponse extends Omit<Movement, "date"> {
  date: string;
}

export interface MovementsResponse {
  totalPages: number;
  page: number;
  movements: MovementResponse[];
}

export interface MovementRequest {
  amount: number;
  category: string;
  type: string;
  date: string;
  description: string;
}

export interface MovementYear {
  month?: number;
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
  totalExpense: number;
  balance: number;
  months: MovementYear[];
}

export interface MovementByMonth {
  totalIncome: number;
  totalExpense: number;
  balance: number;
  days: MovementMonth[];
}

export interface IMovementsRepository {
  createMovement: (movement: MovementRequest) => Promise<Movement>;
  getMovements: (
    page: number,
    institutionIDs: string[]
  ) => Promise<MovementsResponse>;
  getMovementsByYear: (
    year: number,
    institutionIDs: string[]
  ) => Promise<MovementByYear>;
  getMovementsByMonth: (
    year: number,
    month: number,
    institutionIDs: string[]
  ) => Promise<MovementByMonth>;
}

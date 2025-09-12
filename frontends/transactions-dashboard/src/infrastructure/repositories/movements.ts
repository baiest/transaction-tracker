import type {
  IMovementsRepository,
  MovementByYear
} from "@/core/entities/Movement";
import { createFetchClient, API_BASE_URL } from "../http/fetchClient";
import { FetchClient } from "@/core/entities/FetchClient";
import { MovementByMonth } from "@/core/entities/Movement.d";

export class MovementsRepository implements IMovementsRepository {
  private client: FetchClient;

  constructor() {
    this.client = createFetchClient(API_BASE_URL);
  }

  getMovementsByYear(year: number): Promise<MovementByYear> {
    return this.client<MovementByYear>(`/movements/years/${year}`);
  }

  getMovementsByMonth(year: number, month: number): Promise<MovementByMonth> {
    month += 1;

    return this.client<MovementByMonth>(
      `/movements/years/${year}/months/${month}`
    );
  }
}

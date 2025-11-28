import type {
  IMovementsRepository,
  Movement,
  MovementByYear,
  MovementRequest,
  MovementsResponse
} from "@/core/entities/Movement";
import { createFetchClient, API_BASE_URL } from "../http/fetchClient";
import { ExtendedBody, FetchClient } from "@/core/entities/FetchClient";
import { MovementByMonth } from "@/core/entities/Movement.d";

export class MovementsRepository implements IMovementsRepository {
  private client: FetchClient;

  constructor() {
    this.client = createFetchClient(API_BASE_URL);
  }

  createMovement(movement: MovementRequest): Promise<Movement> {
    return this.client<Movement>(`/movements`, {
      method: "POST",
      body: movement
    });
  }

  getMovements(page: number, intitutionIDs: string[]): Promise<MovementsResponse> {
    return this.client<MovementsResponse>(`/movements?page=${page}&limit=20&institution_ids=${intitutionIDs.join(",")}`);
  }

  getMovementsByYear(year: number, intitutionIDs: string[]): Promise<MovementByYear> {
    return this.client<MovementByYear>(`/movements/years/${year}?institution_ids=${intitutionIDs.join(",")}`);
  }

  getMovementsByMonth(year: number, month: number, intitutionIDs: string[]): Promise<MovementByMonth> {
    month += 1;

    return this.client<MovementByMonth>(
      `/movements/years/${year}/months/${month}?institution_ids=${intitutionIDs.join(",")}`
    );
  }
}

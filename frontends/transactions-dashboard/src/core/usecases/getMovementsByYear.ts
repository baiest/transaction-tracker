import type {
  IMovementsRepository,
  MovementByYear
} from "@/core/entities/Movement";

export class GetMovementsByYear {
  constructor(private repository: IMovementsRepository) {}

  async excecute(year: number): Promise<MovementByYear> {
    return this.repository.getMovementsByYear(year);
  }
}

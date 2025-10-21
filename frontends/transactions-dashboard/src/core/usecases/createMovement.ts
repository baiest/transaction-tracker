import type { IMovementsRepository, Movement, MovementRequest } from "@/core/entities/Movement";

export class CreateMovement {
  page = 0;
  totalPages = 0;

  constructor(private repository: IMovementsRepository) {}

  async excecute(movement: MovementRequest): Promise<Movement> {
    return this.repository.createMovement(movement);
  }
}

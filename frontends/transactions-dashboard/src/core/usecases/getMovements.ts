import type { IMovementsRepository, Movement } from "@/core/entities/Movement";

export class GetMovements {
  page = 0;
  totalPages = 0;

  constructor(private repository: IMovementsRepository) {}

  validatePage(page: number) {
    if (page < 0) {
      throw Error("invalid page number");
    }
  }

  async excecute(page = 0, institutionIDs: string[]): Promise<Movement[]> {
    this.validatePage(page);

    const data = await this.repository.getMovements(page, institutionIDs);

    this.page = data.page;
    this.totalPages = data.totalPages;

    return data.movements.map((m) => ({
      ...m,
      date: new Date(m.date)
    }));
  }
}

import type {
  IMovementsRepository,
  MovementByYear
} from "@/core/entities/Movement";

export class GetMovementsByYear {
  constructor(private repository: IMovementsRepository) {}

  validateYear(year: number) {
    const yearNumber = Number(year);
    if (Number.isNaN(yearNumber)) {
      throw Error("year is not a number");
    }

    if (yearNumber > 9999 || yearNumber < 1000) {
      throw Error("invalid year");
    }
  }

  async excecute(
    years: number[],
    institutionIDs: string[]
  ): Promise<MovementByYear[]> {
    years.forEach((y) => {
      this.validateYear(y);
    });

    if (years.length === 1) {
      return [
        await this.repository.getMovementsByYear(years[0], institutionIDs)
      ];
    }

    const data = await Promise.all(
      years.map((y) => this.repository.getMovementsByYear(y, institutionIDs))
    );

    const validData = data.map((d: MovementByYear | null | undefined) =>
      !!d
        ? d
        : {
            balance: 0,
            totalIncome: 0,
            totalExpense: 0,
            months: [{ income: 0, outcome: 0 }]
          }
    );

    return validData;
  }
}

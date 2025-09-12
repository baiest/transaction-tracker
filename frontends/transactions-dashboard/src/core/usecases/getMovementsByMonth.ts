import type {
  IMovementsRepository,
  MovementByMonth
} from "@/core/entities/Movement";

export class GetMovementsByMonth {
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

  validateMonth(month: number) {
    if (month < 0 || month > 11) {
      throw Error("invalid month");
    }
  }

  async excecute(year: number, month: number): Promise<MovementByMonth> {
    this.validateYear(year);
    this.validateMonth(month);

    const movementsByMonth = await this.repository.getMovementsByMonth(
      year,
      month
    );

    let currentIndex = 0;

    const days = Array.from({ length: 31 }, (_, index) => {
      if (movementsByMonth.days[currentIndex]?.day === index + 1) {
        currentIndex += 1;
        return movementsByMonth.days[currentIndex - 1];
      }

      return { income: 0, outcome: 0, day: 0 };
    });

    movementsByMonth.days = days;

    return movementsByMonth;
  }
}

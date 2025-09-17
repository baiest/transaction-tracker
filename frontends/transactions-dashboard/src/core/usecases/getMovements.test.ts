import { describe, it, expect, vi, beforeEach } from "vitest";
import { GetMovements } from "./getMovements";
import { IMovementsRepository } from "../entities/Movement";

const mockMovements = [
  { id: "1", amount: 100, date: "2025-01-01T00:00:00.000Z" },
  { id: "2", amount: 200, date: "2025-02-01T00:00:00.000Z" }
];

const repositoryMock: IMovementsRepository = {
  getMovements: vi.fn(),
  getMovementsByMonth: vi.fn(),
  getMovementsByYear: vi.fn()
};

describe("GetMovements", () => {
  let service: GetMovements;

  beforeEach(() => {
    (repositoryMock.getMovements as ReturnType<typeof vi.fn>).mockReset();
    service = new GetMovements(repositoryMock);
  });

  it("throws error for negative page", () => {
    expect(() => service.validatePage(-1)).toThrow("invalid page number");
  });

  it("sets page and totalPages after execution", async () => {
    (repositoryMock.getMovements as ReturnType<typeof vi.fn>).mockResolvedValue(
      {
        page: 1,
        totalPages: 3,
        movements: mockMovements
      }
    );

    const result = await service.excecute(1);

    expect(service.page).toBe(1);
    expect(service.totalPages).toBe(3);
    expect(result).toHaveLength(2);
    expect(result[0].date).toBeInstanceOf(Date);
    expect(result[0].id).toBe("1");
  });

  it("calls repository.getMovements with correct page", async () => {
    (repositoryMock.getMovements as ReturnType<typeof vi.fn>).mockResolvedValue(
      {
        page: 0,
        totalPages: 1,
        movements: []
      }
    );

    await service.excecute(0);
    expect(repositoryMock.getMovements).toHaveBeenCalledWith(0);
  });

  it("defaults to page 0 if no page is passed", async () => {
    (repositoryMock.getMovements as ReturnType<typeof vi.fn>).mockResolvedValue(
      {
        page: 0,
        totalPages: 1,
        movements: []
      }
    );

    await service.excecute();
    expect(repositoryMock.getMovements).toHaveBeenCalledWith(0);
  });
});

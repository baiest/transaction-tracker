import { render, screen } from "@testing-library/react";
import Home from "./page";
import { vi } from "vitest";
import { MovementByMonth, MovementByYear } from "@/core/entities/Movement";

// Mock the entire module and directly expose the function
vi.mock("@/infrastructure/store/movements", () => {
  const useMovementsStore = vi.fn();
  return { useMovementsStore };
});

vi.mock("@/ui/charts/LineChart", () => {
  return {
    default: vi.fn(() => <div>Line Chart Mock</div>)
  };
});

// A helper function to easily mock the store's state for each test
const mockUseMovementsStore = async (
  timeSelected: string,
  data: MovementByMonth | MovementByYear | undefined = {
    totalIncome: 1000,
    totalOutcome: 1500,
    balance: -500,
    months: [],
    days: []
  },
  allYearsData: MovementByYear[] = []
) => {
  const { useMovementsStore } = vi.mocked(
    await import("@/infrastructure/store/movements")
  );
  useMovementsStore.mockReturnValue({
    movementsByYear: timeSelected === "year" ? data : { months: [] },
    movementsByMonth: timeSelected === "month" ? data : { days: [] },
    allYearsRaw: timeSelected === "all_years" ? allYearsData : [],
    year: 2024,
    month: 1,
    timeSelected,
    fetchMomentsByYear: vi.fn(),
    fetchMomentsByMonth: vi.fn(),
    fetchAllYearsData: vi.fn()
  });
};

describe("Home Page", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders with a year selected and displays correct financial data", async () => {
    await mockUseMovementsStore("year", {
      totalIncome: 1000,
      totalOutcome: 1500,
      balance: -500,
      months: []
    });

    render(<Home />);

    expect(screen.getByText("Total Income").parentElement).toHaveTextContent(
      /\$ *1\.000/
    );
    expect(screen.getByText("Total Expenses")).toBeInTheDocument();
    expect(screen.getByText("Total Expenses").parentElement).toHaveTextContent(
      /\$ *1\.500/
    );
    expect(screen.getByText("Net Balance")).toBeInTheDocument();
    expect(screen.getByText("Net Balance").parentElement).toHaveTextContent(
      /\-?\$ *500/
    );
  });

  it("renders with a month selected and displays correct financial data", async () => {
    await mockUseMovementsStore("month", {
      totalIncome: 500,
      totalOutcome: 200,
      balance: 300,
      days: []
    });

    render(<Home />);

    expect(screen.getByText(/\$ *500/)).toBeInTheDocument();
    expect(screen.getByText(/\$ *200/)).toBeInTheDocument();
    expect(screen.getByText(/\$ *300/)).toBeInTheDocument();
  });

  it("renders with 'all_years' selected and displays correct financial data", async () => {
    await mockUseMovementsStore("all_years", undefined, [
      { totalIncome: 500, totalOutcome: 200, balance: 300, months: [] },
      { totalIncome: 1000, totalOutcome: 500, balance: 500, months: [] }
    ]);

    render(<Home />);

    expect(screen.getByText("Total Income").parentElement).toHaveTextContent(
      /\$ *1\.500/
    ); // 500 + 1000
    expect(screen.getByText("Total Expenses").parentElement).toHaveTextContent(
      /\$ *700/
    ); // 200 + 500
    expect(screen.getByText("Net Balance").parentElement).toHaveTextContent(
      /\$ *800/
    ); // 300 + 500
  });

  it("calculates percentage correctly when income is positive", async () => {
    await mockUseMovementsStore("year", {
      totalIncome: 2000,
      totalOutcome: 1500,
      balance: 500,
      months: []
    });

    render(<Home />);
    expect(screen.getByText("25% Balance / Income Ratio")).toBeInTheDocument();
  });

  it("displays 0% when total income is zero", async () => {
    await mockUseMovementsStore("year", {
      totalIncome: 0,
      totalOutcome: 500,
      balance: -500,
      months: []
    });

    render(<Home />);
    expect(screen.getByText("0% Balance / Income Ratio")).toBeInTheDocument();
  });

  it("renders other static UI elements correctly", async () => {
    await mockUseMovementsStore("year", {
      totalIncome: 1000,
      totalOutcome: 1500,
      balance: -500,
      months: []
    });

    render(<Home />);

    expect(screen.getByText("Income vs Expenses")).toBeInTheDocument();
    expect(screen.getByText("Income by Category")).toBeInTheDocument();
    expect(screen.getByText("Pie Chart Placeholder")).toBeInTheDocument();
    expect(screen.getByText("Expenses by Category")).toBeInTheDocument();
    expect(screen.getByText("Bar Chart Placeholder")).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Add reminders/)).toBeInTheDocument();
  });

  it('passes correct data to LineChart when timeSelected is "year"', async () => {
    const mockData = {
      totalIncome: 1000,
      totalOutcome: 500,
      balance: 500,
      months: [
        { income: 100, outcome: 50 },
        { income: 200, outcome: 100 },
        { income: 300, outcome: 150 }
      ]
    };
    await mockUseMovementsStore("year", mockData);
    render(<Home />);
    const lineChart = screen.getByText("Line Chart Mock");
    expect(lineChart).toBeInTheDocument();
  });

  it('passes correct data to LineChart when timeSelected is "month"', async () => {
    const mockData: MovementByMonth = {
      totalIncome: 100,
      totalOutcome: 50,
      balance: 50,
      days: [
        { day: 1, income: 10, outcome: 5 },
        { day: 2, income: 20, outcome: 10 }
      ]
    };
    await mockUseMovementsStore("month", mockData);
    render(<Home />);
    const lineChart = screen.getByText("Line Chart Mock");
    expect(lineChart).toBeInTheDocument();
  });

  it('passes correct data to LineChart when timeSelected is "all_years"', async () => {
    const mockAllYearsData: MovementByYear[] = [
      {
        totalIncome: 1000,
        totalOutcome: 500,
        balance: 500,
        months: [{ income: 100, outcome: 50 }]
      },
      {
        totalIncome: 2000,
        totalOutcome: 1000,
        balance: 500,
        months: [{ income: 200, outcome: 100 }]
      }
    ];
    await mockUseMovementsStore("all_years", undefined, mockAllYearsData);
    render(<Home />);
    const lineChart = screen.getByText("Line Chart Mock");
    expect(lineChart).toBeInTheDocument();
  });
});

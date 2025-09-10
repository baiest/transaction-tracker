import { render, screen } from "@testing-library/react";
import Home from "./page";

vi.mock("@/infrastructure/store/movements", () => {
  return {
    useMovementsStore: vi.fn().mockReturnValue({
      movementsByYear: {
        totalIncome: 1000,
        totalOutcome: 1500,
        balance: 500
      } as MovementByYear,
      fetchMomentesByYear: vi.fn()
    })
  };
});

import { useMovementsStore } from "@/infrastructure/store/movements";
import type { MovementByYear } from "@/core/entities/Movement";

describe("Home Page", () => {
  it("renders Income card correctly", () => {
    render(<Home />);
    const income = screen.getByText("Income");
    expect(income).toBeInTheDocument();
    expect(income.parentElement).toHaveTextContent(/\$ *1.000/);
    expect(screen.getByText("Primary")).toBeInTheDocument();
    expect(screen.getByText("Other")).toBeInTheDocument();
  });

  it("renders Expenses card correctly", () => {
    render(<Home />);
    const expenses = screen.getByText("Expenses");
    expect(expenses).toBeInTheDocument();
    expect(expenses.parentElement).toHaveTextContent(/\$ *1.500/);
    expect(screen.getByText("Fixed")).toBeInTheDocument();
    expect(screen.getByText("Variable")).toBeInTheDocument();
  });

  it("renders Net Balance card correctly", () => {
    render(<Home />);
    const balance = screen.getByText("Net Balance");
    expect(balance).toBeInTheDocument();
    expect(balance.parentElement).toHaveTextContent(/\$ *500/);
  });

  it("renders Income vs Expenses chart placeholder", () => {
    render(<Home />);
    expect(screen.getByText("Income vs Expenses")).toBeInTheDocument();
    expect(screen.getByText("Chart Placeholder")).toBeInTheDocument();
  });

  it("renders Income by Category chart placeholder", () => {
    render(<Home />);
    expect(screen.getByText("Income by Category")).toBeInTheDocument();
    expect(screen.getByText("Pie Chart Placeholder")).toBeInTheDocument();
  });

  it("renders Expenses by Category chart placeholder", () => {
    render(<Home />);
    expect(screen.getByText("Expenses by Category")).toBeInTheDocument();
    expect(screen.getByText("Bar Chart Placeholder")).toBeInTheDocument();
  });

  it("renders Notes textarea correctly", () => {
    render(<Home />);
    const textarea = screen.getByPlaceholderText(
      "Add reminders, insights, or pending items about your finances here."
    );
    expect(textarea).toBeInTheDocument();
  });

  it("calculates percentage correctly when outcome < income", () => {
    const { fetchMomentesByYear } = useMovementsStore();

    (
      fetchMomentesByYear as unknown as ReturnType<typeof vi.fn>
    ).mockReturnValue({
      movementsByYear: { totalIncome: 1000, totalOutcome: 500, balance: 500 },
      fetchMomentesByYear: vi.fn()
    });

    render(<Home />);
    expect(
      screen.getByText("50% Expenses to Income Ratio")
    ).toBeInTheDocument();
  });

  it("calculates when income is zero", () => {
    const { fetchMomentesByYear } = useMovementsStore();

    (
      fetchMomentesByYear as unknown as ReturnType<typeof vi.fn>
    ).mockReturnValue({
      movementsByYear: { totalIncome: 0, totalOutcome: 500, balance: -500 },
      fetchMomentesByYear: vi.fn()
    });

    useMovementsStore().movementsByYear = {
      totalIncome: 0,
      totalOutcome: 500,
      balance: -500,
      months: []
    };

    render(<Home />);
    expect(screen.getByText("0% Expenses to Income Ratio")).toBeInTheDocument();
  });
});

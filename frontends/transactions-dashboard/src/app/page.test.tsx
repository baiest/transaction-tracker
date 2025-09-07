// Home.test.tsx
import { render, screen } from "@testing-library/react";
import Home from "./page";

describe("Home Page", () => {
  it("renders Income card correctly", () => {
    render(<Home />);
    expect(screen.getByText("Income")).toBeInTheDocument();
    expect(screen.getByText("$12,450")).toBeInTheDocument();
    expect(screen.getByText("Primary")).toBeInTheDocument();
    expect(screen.getByText("Other")).toBeInTheDocument();
  });

  it("renders Expenses card correctly", () => {
    render(<Home />);
    expect(screen.getByText("Expenses")).toBeInTheDocument();
    expect(screen.getByText("$8,320")).toBeInTheDocument();
    expect(screen.getByText("Fixed")).toBeInTheDocument();
    expect(screen.getByText("Variable")).toBeInTheDocument();
  });

  it("renders Net Balance card correctly", () => {
    render(<Home />);
    expect(screen.getByText("Net Balance")).toBeInTheDocument();
    expect(screen.getByText("$4,130")).toBeInTheDocument();
    expect(screen.getByText("+6.2% vs last month")).toBeInTheDocument();
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
});

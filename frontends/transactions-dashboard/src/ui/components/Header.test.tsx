import { render, screen, fireEvent } from "@testing-library/react";
import Header from "./Header";

const setYearMock = vi.fn();
const setMonthMock = vi.fn();
const setTimeSelectedMock = vi.fn();

vi.mock("next/navigation", () => ({
  usePathname: vi.fn(() => "/"),
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    refresh: vi.fn()
  })
}));

const mockStore = (overrides = {}) => ({
  year: 2025,
  month: 0,
  timeSelected: "year",
  setYear: setYearMock,
  setMonth: setMonthMock,
  setTimeSelected: setTimeSelectedMock,
  ...overrides
});

vi.mock("@/infrastructure/store/movements", () => ({
  useMovementsStore: vi.fn(() => mockStore())
}));

import { usePathname } from "next/navigation";
import { useMovementsStore } from "@/infrastructure/store/movements";

describe("Header component", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/");
    (useMovementsStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue(
      mockStore()
    );
  });

  it("renders selects and buttons", () => {
    render(<Header />);
    expect(screen.getByDisplayValue("2025")).toBeInTheDocument();
    expect(screen.getByText("Year")).toBeInTheDocument();
    expect(screen.getByText("Export")).toBeInTheDocument();
    expect(screen.getByText("New Transaction")).toBeInTheDocument();
  });

  it("calls setYear when year changes", () => {
    render(<Header />);
    const select = screen.getByDisplayValue("2025");
    fireEvent.change(select, { target: { value: "2023" } });
    expect(setYearMock).toHaveBeenCalledWith(2023);
  });

  it("calls setTimeSelected when the view mode changes", () => {
    render(<Header />);
    const viewModeSelect = screen.getAllByRole("combobox")[0];
    fireEvent.change(viewModeSelect, { target: { value: "all_years" } });
    expect(setTimeSelectedMock).toHaveBeenCalledWith("all_years");

    fireEvent.change(viewModeSelect, { target: { value: "year" } });
    expect(setTimeSelectedMock).toHaveBeenCalledWith("year");
  });

  it("renders month select when timeSelected is month", () => {
    (useMovementsStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue(
      mockStore({ timeSelected: "month", month: 1 })
    );

    render(<Header />);
    const monthSelect = screen.getByDisplayValue("Feb");
    expect(monthSelect).toBeInTheDocument();

    fireEvent.change(monthSelect, { target: { value: "3" } });
    expect(setMonthMock).toHaveBeenCalledWith(3);
  });

  it("does not render selects if pathname is not root", () => {
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/settings");

    render(<Header />);
    expect(screen.queryAllByRole("combobox")).toHaveLength(0);
    expect(screen.getByText("Export")).toBeInTheDocument();
    expect(screen.getByText("New Transaction")).toBeInTheDocument();
  });
});

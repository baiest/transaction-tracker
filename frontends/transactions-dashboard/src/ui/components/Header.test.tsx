import { render, screen, fireEvent } from "@testing-library/react";
import Header from "./Header";

const setYearMock = vi.fn();
const setTimeSelectedMock = vi.fn();

vi.mock("@/infrastructure/store/movements", () => {
  return {
    useMovementsStore: () => ({
      year: 2025,
      timeSelected: "year",
      setYear: setYearMock,
      setTimeSelected: setTimeSelectedMock
    })
  };
});

describe("Header component", () => {
  beforeEach(() => {
    vi.clearAllMocks();
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

    const selects = screen.getAllByRole("combobox");

    const viewModeSelect = selects[0];

    fireEvent.change(viewModeSelect, { target: { value: "all_years" } });
    expect(setTimeSelectedMock).toHaveBeenCalledWith("all_years");

    fireEvent.change(viewModeSelect, { target: { value: "year" } });
    expect(setTimeSelectedMock).toHaveBeenCalledWith("year");
  });
});

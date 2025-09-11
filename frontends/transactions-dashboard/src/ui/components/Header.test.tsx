import { render, screen, fireEvent } from "@testing-library/react";
import Header from "./Header";

const setYearMock = vi.fn();
const setShowAllYearsMock = vi.fn();

vi.mock("@/infrastructure/store/movements", () => {
  return {
    useMovementsStore: () => ({
      year: 2025,
      showAllYears: false,
      setYear: setYearMock,
      setShowAllYears: setShowAllYearsMock
    })
  };
});

describe("Header component", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders title correctly", () => {
    render(<Header />);
    expect(screen.getByText("Transactions")).toBeInTheDocument();
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

  it("calls setShowAllYears when the view mode changes", () => {
    render(<Header />);
    const select = screen.getByDisplayValue("Year");
    fireEvent.change(select, { target: { value: "all-years" } });
    expect(setShowAllYearsMock).toHaveBeenCalledWith(true);

    fireEvent.change(select, { target: { value: "Year" } });
    expect(setShowAllYearsMock).toHaveBeenCalledWith(false);
  });
});

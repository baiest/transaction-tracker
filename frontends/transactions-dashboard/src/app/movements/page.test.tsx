import { render, screen, fireEvent } from "@testing-library/react";
import Movements from "./page";
import { useMovementsStore } from "@/infrastructure/store/movements";

vi.mock("@/infrastructure/store/movements", () => {
  return {
    useMovementsStore: vi.fn()
  };
});

vi.mock("@/hooks/useFormatCurrency", () => ({
  useFormatCurrency: () => (val: number) => `$${val}`
}));

const pushMock = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: pushMock }),
  useSearchParams: () =>
    new URLSearchParams({
      page: "1"
    })
}));

vi.mock("@tanstack/react-table", () => {
  return {
    useReactTable: vi.fn(() => ({
      getHeaderGroups: vi.fn(() => [
        {
          id: "headerGroup",
          headers: [
            {
              id: "h_detail",
              isPlaceholder: false,
              column: { columnDef: { header: "Details" } },
              getContext: () => ({})
            },
            {
              id: "h_date",
              isPlaceholder: false,
              column: { columnDef: { header: "Date" } },
              getContext: () => ({})
            },
            {
              id: "h_value",
              isPlaceholder: false,
              column: { columnDef: { header: "Amount" } },
              getContext: () => ({})
            }
          ]
        }
      ]),
      getRowModel: vi.fn(() => ({
        rows: [
          {
            id: "1",
            original: { isNegative: false },
            getIsSelected: vi.fn(),
            getValue: (key: string) =>
              key === "detail"
                ? "Salary"
                : key === "date"
                ? "2025-09-01"
                : key === "value"
                ? 5000
                : null,
            getVisibleCells: () => [
              {
                id: "1_detail",
                column: { columnDef: { cell: () => "Salary" } },
                getContext: () => ({})
              },
              {
                id: "1_date",
                column: { columnDef: { cell: () => "2025-09-01" } },
                getContext: () => ({})
              },
              {
                id: "1_value",
                column: { columnDef: { cell: () => "5000" } },
                getContext: () => ({})
              }
            ]
          },
          {
            id: "2",
            original: { isNegative: true },
            getIsSelected: vi.fn(),
            getValue: (key: string) =>
              key === "detail"
                ? "Groceries"
                : key === "date"
                ? "2025-09-05"
                : key === "value"
                ? 200
                : null,
            getVisibleCells: () => [
              {
                id: "2_detail",
                column: { columnDef: { cell: () => "Groceries" } },
                getContext: () => ({})
              },
              {
                id: "2_date",
                column: { columnDef: { cell: () => "2025-09-05" } },
                getContext: () => ({})
              },
              {
                id: "2_value",
                column: { columnDef: { cell: () => "200" } },
                getContext: () => ({})
              }
            ]
          }
        ]
      })),
      getState: vi.fn(() => ({
        pagination: { pageIndex: 0, pageSize: 10 }
      })),
      getAllLeafColumns: vi.fn(() => [])
    })),
    getCoreRowModel: vi.fn(),
    getFilteredRowModel: vi.fn(),
    flexRender: vi.fn((renderer: unknown, ctx: unknown) =>
      typeof renderer === "function" ? renderer(ctx) : renderer
    )
  };
});

describe("Movements", () => {
  const fetchMovementsMock = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    (useMovementsStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      movements: {
        1: {
          detail: "Salary",
          date: "2025-09-01",
          value: 5000,
          isNegative: false
        },
        2: {
          detail: "Groceries",
          date: "2025-09-05",
          value: 200,
          isNegative: true
        }
      },
      fetchMovements: fetchMovementsMock,
      totalPages: 3
    });
  });

  it("renders table with movement data", () => {
    render(<Movements />);

    expect(screen.getByText("Salary")).toBeInTheDocument();
    expect(screen.getByText("Groceries")).toBeInTheDocument();
    expect(screen.getByText("2025-09-01")).toBeInTheDocument();
    expect(screen.getByText("2025-09-05")).toBeInTheDocument();

    expect(screen.getByText("5000")).toBeInTheDocument();
    expect(screen.getByText("200")).toBeInTheDocument();
  });

  it("calls fetchMovements on mount with page 1", () => {
    render(<Movements />);
    expect(fetchMovementsMock).toHaveBeenCalledWith(1);
  });

  it("navigates to new page when pagination is clicked", () => {
    render(<Movements />);

    const page2 = screen.getAllByText("2")[0];
    fireEvent.click(page2);

    expect(pushMock).toHaveBeenCalledWith("/movements?page=2");
  });
});

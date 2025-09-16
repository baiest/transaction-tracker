import { render, screen, fireEvent } from "@testing-library/react";
import { DataTable } from "./DataTable";
import * as RT from "@tanstack/react-table";

interface Row {
  id: number;
  name: string;
}

const columns: RT.ColumnDef<Row>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" }
];

const data: Row[] = [
  { id: 1, name: "Alice" },
  { id: 2, name: "Bob" }
];

vi.mock("@tanstack/react-table", () => {
  return {
    useReactTable: vi.fn(() => ({
      getHeaderGroups: vi.fn(() => []),
      getRowModel: vi.fn(() => ({ rows: [] })),
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

describe("DataTable", () => {
  it("muestra input de búsqueda y botón de columnas", () => {
    render(
      <DataTable
        data={[]}
        columns={columns}
        totalPages={1}
        pageIndex={0}
        pageSize={10}
        onPageChange={() => {}}
      />
    );

    expect(screen.getByPlaceholderText("Search...")).toBeInTheDocument();
    expect(screen.getByText("Columns")).toBeInTheDocument();
    expect(screen.getByText("No results.")).toBeInTheDocument();
  });

  it("renderiza headers y filas de datos", () => {
    (RT.useReactTable as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      getHeaderGroups: vi.fn(() => [
        {
          id: "headerGroup1",
          headers: [
            {
              id: "col-id",
              isPlaceholder: false,
              column: { columnDef: { header: "ID" } },
              getContext: vi.fn()
            },
            {
              id: "col-name",
              isPlaceholder: false,
              column: { columnDef: { header: "Name" } },
              getContext: vi.fn()
            }
          ]
        }
      ]),
      getRowModel: vi.fn(() => ({
        rows: [
          {
            id: "row1",
            getIsSelected: () => false,
            getVisibleCells: () => [
              {
                id: "row1-cell1",
                column: { columnDef: { cell: () => "1" } },
                getContext: vi.fn()
              },
              {
                id: "row1-cell2",
                column: { columnDef: { cell: () => "Alice" } },
                getContext: vi.fn()
              }
            ]
          },
          {
            id: "row2",
            getIsSelected: () => false,
            getVisibleCells: () => [
              {
                id: "row2-cell1",
                column: { columnDef: { cell: () => "2" } },
                getContext: vi.fn()
              },
              {
                id: "row2-cell2",
                column: { columnDef: { cell: () => "Bob" } },
                getContext: vi.fn()
              }
            ]
          }
        ]
      })),
      getAllLeafColumns: vi.fn(() => []),
      getState: vi.fn(() => ({
        pagination: { pageIndex: 0, pageSize: 10 }
      }))
    });

    render(
      <DataTable
        data={[]}
        columns={[]}
        totalPages={1}
        pageIndex={0}
        pageSize={10}
        onPageChange={() => {}}
      />
    );

    expect(screen.getByText("ID")).toBeInTheDocument();
    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
  });

  it("deshabilita Previous en la primera página y Next en la última", () => {
    const { rerender } = render(
      <DataTable
        data={data}
        columns={columns}
        totalPages={3}
        pageIndex={0}
        pageSize={10}
        onPageChange={() => {}}
      />
    );

    expect(screen.getByText("Previous")).toBeDisabled();
    expect(screen.getByText("Next")).not.toBeDisabled();

    rerender(
      <DataTable
        data={data}
        columns={columns}
        totalPages={3}
        pageIndex={3}
        pageSize={10}
        onPageChange={() => {}}
      />
    );

    expect(screen.getByText("Next")).toBeDisabled();
  });

  it("llama onPageChange al hacer click en Next y Previous", () => {
    const onPageChange = vi.fn();

    render(
      <DataTable
        data={data}
        columns={columns}
        totalPages={3}
        pageIndex={1}
        pageSize={10}
        onPageChange={onPageChange}
      />
    );

    fireEvent.click(screen.getByText("Next"));
    expect(onPageChange).toHaveBeenCalledWith(2);

    fireEvent.click(screen.getByText("Previous"));
    expect(onPageChange).toHaveBeenCalledWith(0);
  });

  it("llama onPageChange al hacer click en un número de página", () => {
    const onPageChange = vi.fn();

    render(
      <DataTable
        data={data}
        columns={columns}
        totalPages={3}
        pageIndex={0}
        pageSize={10}
        onPageChange={onPageChange}
      />
    );

    fireEvent.click(screen.getByRole("button", { name: "2" }));
    expect(onPageChange).toHaveBeenCalledWith(1);
  });
});

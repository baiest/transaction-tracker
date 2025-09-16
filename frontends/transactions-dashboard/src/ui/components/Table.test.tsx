import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow
} from "./Table";

describe("Table Components", () => {
  it("should render the Table component with a div container and a table element", () => {
    render(<Table />);

    const container = screen.getByRole("table").parentElement;
    expect(container).toHaveAttribute("data-slot", "table-container");

    const table = screen.getByRole("table");
    expect(table).toHaveAttribute("data-slot", "table");
  });

  it("should render the TableHeader component", () => {
    render(
      <table>
        <TableHeader data-testid="table-header" />
      </table>
    );
    const tableHeader = screen.getByTestId("table-header");
    expect(tableHeader).toBeInTheDocument();
    expect(tableHeader.tagName).toBe("THEAD");
    expect(tableHeader).toHaveAttribute("data-slot", "table-header");
    expect(tableHeader).toHaveClass("[&_tr]:border-b");
  });

  it("should render the TableBody component", () => {
    render(
      <table>
        <TableBody data-testid="table-body" />
      </table>
    );
    const tableBody = screen.getByTestId("table-body");
    expect(tableBody).toBeInTheDocument();
    expect(tableBody.tagName).toBe("TBODY");
    expect(tableBody).toHaveAttribute("data-slot", "table-body");
    expect(tableBody).toHaveClass("[&_tr:last-child]:border-0");
  });

  it("should render the TableFooter component", () => {
    render(
      <table>
        <TableFooter data-testid="table-footer" />
      </table>
    );
    const tableFooter = screen.getByTestId("table-footer");
    expect(tableFooter).toBeInTheDocument();
    expect(tableFooter.tagName).toBe("TFOOT");
    expect(tableFooter).toHaveAttribute("data-slot", "table-footer");
    expect(tableFooter).toHaveClass("bg-muted/50 border-t font-medium");
  });

  it("should render the TableRow component", () => {
    render(
      <table>
        <tbody>
          <TableRow data-testid="table-row" />
        </tbody>
      </table>
    );
    const tableRow = screen.getByTestId("table-row");
    expect(tableRow).toBeInTheDocument();
    expect(tableRow.tagName).toBe("TR");
    expect(tableRow).toHaveAttribute("data-slot", "table-row");
    expect(tableRow).toHaveClass(
      "hover:bg-muted/50 data-[state=selected]:bg-muted"
    );
  });

  it("should render the TableHead component", () => {
    render(
      <table>
        <thead>
          <tr>
            <TableHead data-testid="table-head">Header</TableHead>
          </tr>
        </thead>
      </table>
    );
    const tableHead = screen.getByRole("columnheader", { name: /header/i });
    expect(tableHead).toBeInTheDocument();
    expect(tableHead.tagName).toBe("TH");
    expect(tableHead).toHaveAttribute("data-slot", "table-head");
    expect(tableHead).toHaveClass("h-10 px-2 text-left");
  });

  it("should render the TableCell component", () => {
    render(
      <table>
        <tbody>
          <tr>
            <TableCell data-testid="table-cell">Cell</TableCell>
          </tr>
        </tbody>
      </table>
    );
    const tableCell = screen.getByRole("cell", { name: /cell/i });
    expect(tableCell).toBeInTheDocument();
    expect(tableCell.tagName).toBe("TD");
    expect(tableCell).toHaveAttribute("data-slot", "table-cell");
    expect(tableCell).toHaveClass("p-2 align-middle");
  });

  it("should render the TableCaption component", () => {
    render(
      <table>
        <TableCaption data-testid="table-caption">Caption</TableCaption>
      </table>
    );
    const tableCaption = screen.getByTestId("table-caption");
    expect(tableCaption).toBeInTheDocument();
    expect(tableCaption.tagName).toBe("CAPTION");
    expect(tableCaption).toHaveAttribute("data-slot", "table-caption");
    expect(tableCaption).toHaveClass("text-muted-foreground mt-4 text-sm");
  });

  it("should render a full table structure correctly", () => {
    render(
      <Table>
        <TableCaption>A list of items.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>Item</TableHead>
            <TableHead>Price</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell>Laptop</TableCell>
            <TableCell>$1200</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>Mouse</TableCell>
            <TableCell>$50</TableCell>
          </TableRow>
        </TableBody>
        <TableFooter>
          <TableRow>
            <TableCell>Total</TableCell>
            <TableCell>$1250</TableCell>
          </TableRow>
        </TableFooter>
      </Table>
    );

    expect(screen.getByRole("table")).toBeInTheDocument();
    expect(screen.getByText("A list of items.")).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: "Item" })
    ).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "$1200" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "Total" })).toBeInTheDocument();
  });
});

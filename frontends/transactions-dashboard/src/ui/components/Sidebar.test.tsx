import { render, screen } from "@testing-library/react";
import Sidebar from "./Sidebar";

vi.mock("next/navigation", () => ({
  usePathname: vi.fn()
}));

import { usePathname } from "next/navigation";

describe("Sidebar component", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders all navigation options", () => {
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/");

    render(<Sidebar />);

    expect(screen.getByText("Dashboard")).toBeInTheDocument();
    expect(screen.getByText("Movimientos")).toBeInTheDocument();
    expect(screen.getByText("Cuentas")).toBeInTheDocument();
    expect(screen.getByText("Metas")).toBeInTheDocument();
    expect(screen.getByText("Ajustes")).toBeInTheDocument();
  });

  it("applies selected style when pathname matches", () => {
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/movements");

    render(<Sidebar />);

    const movementsLink = screen.getByText("Movimientos");

    expect(movementsLink).toHaveClass(
      "flex",
      "items-center",
      "p-3",
      "rounded-md",
      "bg-indigo-100",
      "text-indigo-600",
      "font-medium"
    );
  });

  it("does not apply selected style to non-active links", () => {
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/movements");

    render(<Sidebar />);

    const dashboardLink = screen.getByText("Dashboard");

    expect(dashboardLink).toHaveClass("text-gray-600");
    expect(dashboardLink).not.toHaveClass("bg-indigo-100");
  });

  it("applies custom className when passed", () => {
    (usePathname as ReturnType<typeof vi.fn>).mockReturnValue("/");

    render(<Sidebar className="extra-class" />);

    const aside = screen.getByRole("complementary");
    expect(aside).toHaveClass("extra-class");
  });
});

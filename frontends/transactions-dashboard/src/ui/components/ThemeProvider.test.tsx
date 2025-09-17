import { render, screen } from "@testing-library/react";
import ThemeProvider from "./ThemeProvider";

vi.mock("next-themes", () => ({
  ThemeProvider: ({ children }: unknonw) => (
    <div data-testid="mock">{children}</div>
  )
}));

describe("ThemeProvider", () => {
  it("renders children hidden before mounted", async () => {
    render(
      <ThemeProvider>
        <div>Child Content</div>
      </ThemeProvider>
    );
    const child = await screen.findByText("Child Content");
    expect(child).toBeInTheDocument();
    expect(screen.getByTestId("mock")).toBeInTheDocument();
  });

  it("renders children with NextThemesProvider after mount", async () => {
    render(
      <ThemeProvider>
        <div>Child Content</div>
      </ThemeProvider>
    );
    const child = await screen.findByText("Child Content");
    expect(child).toBeInTheDocument();
  });
});

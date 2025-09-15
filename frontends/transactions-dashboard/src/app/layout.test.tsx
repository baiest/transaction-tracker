import { render, screen } from "@testing-library/react";
import RootLayout from "./layout";

vi.mock("next-themes", () => ({
  ThemeProvider: ({ children }: any) => <div data-testid="mock">{children}</div>
}));

describe("RootLayout", () => {
  it("renders children", () => {
    render(
      <RootLayout>
        <div data-testid="child">Hello</div>
      </RootLayout>,
      { container: document.documentElement }
    );

    expect(screen.getByTestId("child")).toBeInTheDocument();
  });
});

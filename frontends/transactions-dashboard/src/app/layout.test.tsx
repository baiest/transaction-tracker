import { render, screen } from "@testing-library/react";
import RootLayout from "./layout";

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

import { render, screen } from "@testing-library/react";
import { Input } from "./Input";

describe("Input", () => {
  it("should render the input element with default classes and data-slot", () => {
    render(<Input />);
    const input = screen.getByRole("textbox");
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("data-slot", "input");
    expect(input).toHaveClass(
      "h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1"
    );
    expect(input).toHaveClass("placeholder:text-muted-foreground");
  });

  it("should apply the 'type' attribute correctly", () => {
    render(<Input type="email" />);
    const input = screen.getByRole("textbox");
    expect(input).toHaveAttribute("type", "email");
  });

  it("should merge custom class names with the default classes", () => {
    render(<Input className="custom-class" />);
    const input = screen.getByRole("textbox");
    expect(input).toHaveClass(
      "h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 custom-class"
    );
  });

  it("should apply the 'disabled' attribute and classes when disabled", () => {
    render(<Input disabled />);
    const input = screen.getByRole("textbox");
    expect(input).toBeDisabled();
    expect(input).toHaveClass(
      "disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
    );
  });

  it("should apply the 'aria-invalid' attribute and classes when aria-invalid is true", () => {
    render(<Input aria-invalid={true} />);
    const input = screen.getByRole("textbox", { name: "" });
    expect(input).toHaveAttribute("aria-invalid", "true");
    expect(input).toHaveClass(
      "aria-invalid:ring-destructive/20 aria-invalid:border-destructive"
    );
  });

  it("should display the placeholder text", () => {
    const placeholderText = "Enter your email";
    render(<Input placeholder={placeholderText} />);
    const input = screen.getByPlaceholderText(placeholderText);
    expect(input).toBeInTheDocument();
  });
});

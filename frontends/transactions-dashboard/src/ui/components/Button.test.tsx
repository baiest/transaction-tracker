import { render, screen } from "@testing-library/react";
import { Button } from "./Button";
import { Home } from "lucide-react";

describe("Button", () => {
  it("should render the button with default variant and size", () => {
    render(<Button>Click me</Button>);
    const button = screen.getByRole("button", { name: /click me/i });
    expect(button).toBeInTheDocument();
    expect(button).toHaveClass("bg-primary text-primary-foreground");
    expect(button).toHaveClass("h-9 px-4 py-2");
  });

  it("should apply the correct classes for the 'destructive' variant", () => {
    render(<Button variant="destructive">Delete</Button>);
    const button = screen.getByRole("button", { name: /delete/i });
    expect(button).toHaveClass(
      "bg-destructive text-white shadow-xs hover:bg-destructive/90"
    );
  });

  it("should apply the correct classes for the 'outline' variant", () => {
    render(<Button variant="outline">Outline</Button>);
    const button = screen.getByRole("button", { name: /outline/i });
    expect(button).toHaveClass("border bg-background shadow-xs");
  });

  it("should apply the correct classes for the 'sm' size", () => {
    render(<Button size="sm">Small</Button>);
    const button = screen.getByRole("button", { name: /small/i });
    expect(button).toHaveClass("h-8 rounded-md gap-1.5 px-3");
  });

  it("should apply the correct classes for the 'icon' size", () => {
    render(
      <Button size="icon">
        <Home />
      </Button>
    );
    const button = screen.getByRole("button");
    expect(button).toHaveClass("size-9");
  });

  it("should apply the correct classes when both variant and size are specified", () => {
    render(
      <Button variant="secondary" size="lg">
        Large Secondary
      </Button>
    );
    const button = screen.getByRole("button", { name: /large secondary/i });
    expect(button).toHaveClass("h-10 rounded-md px-6");
    expect(button).toHaveClass("bg-secondary text-secondary-foreground");
  });

  it("should render a different element when asChild is true", () => {
    render(
      <Button asChild>
        <a href="/dashboard">Dashboard</a>
      </Button>
    );
    const link = screen.getByRole("link", { name: /dashboard/i });
    expect(link).toBeInTheDocument();
    expect(link.tagName).toBe("A");
    expect(link).toHaveClass("bg-primary text-primary-foreground");
  });

  it("should apply disabled classes when the button is disabled", () => {
    render(<Button disabled>Disabled</Button>);
    const button = screen.getByRole("button", { name: /disabled/i });
    expect(button).toHaveAttribute("disabled");
    expect(button).toHaveClass("disabled:opacity-50");
  });
});

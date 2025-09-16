import { render, screen } from "@testing-library/react";

import {
  DropdownMenu,
  DropdownMenuPortal,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuCheckboxItem,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubTrigger,
  DropdownMenuSubContent
} from "./DropdownMenu";

vi.mock("@radix-ui/react-dropdown-menu", async (importOriginal) => {
  const actual: React.FC = await importOriginal();
  return {
    ...actual,
    Root: vi.fn((props) => <div data-testid="radix-root" {...props} />),
    Portal: vi.fn((props) => <div data-testid="radix-portal" {...props} />),
    Trigger: vi.fn((props) => (
      <button data-testid="radix-trigger" {...props} />
    )),
    Content: vi.fn((props) => <div data-testid="radix-content" {...props} />),
    Group: vi.fn((props) => <div data-testid="radix-group" {...props} />),
    Item: vi.fn((props) => <div data-testid="radix-item" {...props} />),
    CheckboxItem: vi.fn((props) => (
      <div data-testid="radix-checkbox-item" {...props} />
    )),
    ItemIndicator: vi.fn((props) => (
      <div data-testid="radix-item-indicator" {...props} />
    )),
    RadioGroup: vi.fn((props) => (
      <div data-testid="radix-radio-group" {...props} />
    )),
    RadioItem: vi.fn((props) => (
      <div data-testid="radix-radio-item" {...props} />
    )),
    Label: vi.fn((props) => <div data-testid="radix-label" {...props} />),
    Separator: vi.fn((props) => (
      <div data-testid="radix-separator" {...props} />
    )),
    Sub: vi.fn((props) => <div data-testid="radix-sub" {...props} />),
    SubTrigger: vi.fn((props) => (
      <div data-testid="radix-sub-trigger" {...props} />
    )),
    SubContent: vi.fn((props) => (
      <div data-testid="radix-sub-content" {...props} />
    ))
  };
});

describe("DropdownMenu Components", () => {
  it("should render DropdownMenu with correct data-slot", () => {
    render(<DropdownMenu />);
    const root = screen.getByTestId("radix-root");
    expect(root).toHaveAttribute("data-slot", "dropdown-menu");
  });

  it("should render DropdownMenuPortal with correct data-slot", () => {
    render(<DropdownMenuPortal />);
    const portal = screen.getByTestId("radix-portal");
    expect(portal).toHaveAttribute("data-slot", "dropdown-menu-portal");
  });

  it("should render DropdownMenuTrigger with correct data-slot", () => {
    render(<DropdownMenuTrigger />);
    const trigger = screen.getByTestId("radix-trigger");
    expect(trigger).toHaveAttribute("data-slot", "dropdown-menu-trigger");
  });

  it("should render DropdownMenuContent with correct data-slot and classes", () => {
    render(<DropdownMenuContent />);
    const content = screen.getByTestId("radix-content");
    expect(content).toHaveAttribute("data-slot", "dropdown-menu-content");
    expect(content).toHaveClass(
      "bg-popover text-popover-foreground z-50 min-w-[8rem] border p-1 shadow-md"
    );
  });

  it("should render DropdownMenuGroup with correct data-slot", () => {
    render(<DropdownMenuGroup />);
    const group = screen.getByTestId("radix-group");
    expect(group).toHaveAttribute("data-slot", "dropdown-menu-group");
  });

  it("should render DropdownMenuItem with correct data-slot and classes", () => {
    render(<DropdownMenuItem />);
    const item = screen.getByTestId("radix-item");
    expect(item).toHaveAttribute("data-slot", "dropdown-menu-item");
    expect(item).toHaveClass(
      "flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm"
    );
  });

  it("should apply `inset` classes to DropdownMenuItem", () => {
    render(<DropdownMenuItem inset>Item</DropdownMenuItem>);
    const item = screen.getByTestId("radix-item");
    expect(item).toHaveAttribute("data-inset", "true");
    expect(item).toHaveClass("data-[inset]:pl-8");
  });

  it("should apply `destructive` variant classes to DropdownMenuItem", () => {
    render(<DropdownMenuItem variant="destructive">Delete</DropdownMenuItem>);
    const item = screen.getByTestId("radix-item");
    expect(item).toHaveAttribute("data-variant", "destructive");
    expect(item).toHaveClass("data-[variant=destructive]:text-destructive");
  });

  it("should render DropdownMenuCheckboxItem with check icon", () => {
    render(
      <DropdownMenuCheckboxItem checked>Checkbox Item</DropdownMenuCheckboxItem>
    );
    const item = screen.getByTestId("radix-checkbox-item");
    expect(item).toHaveAttribute("data-slot", "dropdown-menu-checkbox-item");
    expect(item).toHaveClass("pl-8");
    expect(item.querySelector("svg")).toBeInTheDocument();
  });

  it("should render DropdownMenuRadioItem with circle icon", () => {
    render(
      <DropdownMenuRadioItem value="test">Radio Item</DropdownMenuRadioItem>
    );
    const item = screen.getByTestId("radix-radio-item");
    expect(item).toHaveAttribute("data-slot", "dropdown-menu-radio-item");
    expect(item).toHaveClass("pl-8");
    expect(item.querySelector("svg")).toBeInTheDocument();
  });

  it("should render DropdownMenuLabel with correct data-slot and classes", () => {
    render(<DropdownMenuLabel />);
    const label = screen.getByTestId("radix-label");
    expect(label).toHaveAttribute("data-slot", "dropdown-menu-label");
    expect(label).toHaveClass("px-2 py-1.5 text-sm font-medium");
  });

  it("should render DropdownMenuSeparator with correct data-slot and classes", () => {
    render(<DropdownMenuSeparator />);
    const separator = screen.getByTestId("radix-separator");
    expect(separator).toHaveAttribute("data-slot", "dropdown-menu-separator");
    expect(separator).toHaveClass("bg-border -mx-1 my-1 h-px");
  });

  it("should render DropdownMenuShortcut with correct data-slot and classes", () => {
    render(<DropdownMenuShortcut />);
    const shortcut = screen.getByTestId("dropdown-menu-shortcut");
    expect(shortcut).toHaveAttribute("data-slot", "dropdown-menu-shortcut");
    expect(shortcut).toHaveClass(
      "text-muted-foreground ml-auto text-xs tracking-widest"
    );
  });

  it("should render DropdownMenuSub with correct data-slot", () => {
    render(<DropdownMenuSub />);
    const sub = screen.getByTestId("radix-sub");
    expect(sub).toHaveAttribute("data-slot", "dropdown-menu-sub");
  });

  it("should render DropdownMenuSubTrigger with correct data-slot and classes", () => {
    render(<DropdownMenuSubTrigger>Submenu</DropdownMenuSubTrigger>);
    const subTrigger = screen.getByTestId("radix-sub-trigger");
    expect(subTrigger).toHaveAttribute(
      "data-slot",
      "dropdown-menu-sub-trigger"
    );
    expect(subTrigger).toHaveClass(
      "flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm"
    );
    expect(subTrigger.querySelector("svg")).toBeInTheDocument();
  });

  it("should render DropdownMenuSubContent with correct data-slot and classes", () => {
    render(<DropdownMenuSubContent />);
    const subContent = screen.getByTestId("radix-sub-content");
    expect(subContent).toHaveAttribute(
      "data-slot",
      "dropdown-menu-sub-content"
    );
    expect(subContent).toHaveClass(
      "z-50 min-w-[8rem] rounded-md border p-1 shadow-lg"
    );
  });

  it("should render DropdownMenuRadioGroup with correct data-slot", () => {
    render(<DropdownMenuRadioGroup />);
    const radioGroup = screen.getByTestId("radix-radio-group");
    expect(radioGroup).toBeInTheDocument();
    expect(radioGroup).toHaveAttribute(
      "data-slot",
      "dropdown-menu-radio-group"
    );
  });
});

import { renderHook } from "@testing-library/react";
import { useFormatCurrency } from "./useFormatCurrency";

describe("useFormatCurrency", () => {
  it("should format number in default COP without decimals", () => {
    const { result } = renderHook(() => useFormatCurrency());

    const formatCurrency = result.current;
    expect(formatCurrency(1000)).toMatch(/^\$\s?1\.000$/);
  });

  it("should format number in default COP with decimals", () => {
    const { result } = renderHook(() => useFormatCurrency());

    const formatCurrency = result.current;
    expect(formatCurrency(1000.5, { withDecimals: true })).toBe("$ 1.000,50");
  });

  it("should allow overriding currency", () => {
    const { result } = renderHook(() => useFormatCurrency());

    const formatCurrency = result.current;
    expect(formatCurrency(1000, { currency: "USD" })).toMatch(
      /^US\$\s?1\.000$/
    );
  });

  it("should allow overriding locale", () => {
    const { result } = renderHook(() => useFormatCurrency());

    const formatCurrency = result.current;
    expect(formatCurrency(1000, { locale: "en-US", currency: "USD" })).toBe(
      "$1,000"
    );
  });

  it("should respect decimals in custom locale", () => {
    const { result } = renderHook(() => useFormatCurrency());

    const formatCurrency = result.current;
    expect(
      formatCurrency(1234.56, {
        locale: "en-US",
        currency: "USD",
        withDecimals: true
      })
    ).toBe("$1,234.56");
  });

  it("should memoize formatter and return same function instance", () => {
    const { result, rerender } = renderHook(() => useFormatCurrency());
    const fn1 = result.current;

    rerender();

    const fn2 = result.current;
    expect(fn1).toBe(fn2);
  });
});

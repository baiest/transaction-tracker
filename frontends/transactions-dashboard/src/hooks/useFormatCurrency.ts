import { useCallback } from "react";

export function useFormatCurrency(
  defaultCurrency: string = "COP",
  defaultLocale: string = "es-CO"
) {
  return useCallback(
    (
      value: number,
      options?: {
        currency?: string;
        locale?: string;
        withDecimals?: boolean;
      }
    ) => {
      const {
        currency = defaultCurrency,
        locale = defaultLocale,
        withDecimals = false
      } = options || {};

      return new Intl.NumberFormat(locale, {
        style: "currency",
        currency,
        minimumFractionDigits: withDecimals ? 2 : 0
      }).format(value);
    },
    [defaultCurrency, defaultLocale]
  );
}

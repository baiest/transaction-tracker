"use client";

import { useFormatCurrency } from "@/hooks/useFormatCurrency";
import { useMovementsStore } from "@/infrastructure/store/movements";
import { useCallback, useEffect } from "react";

export default function Home() {
  const { movementsByYear, fetchMomentesByYear } = useMovementsStore();

  const format = useFormatCurrency();

  const calculateVerticalPercentage = useCallback((): number => {
    if (movementsByYear.totalIncome === 0) {
      return 0;
    }

    const percentage =
      (movementsByYear.balance / movementsByYear.totalIncome) * 100;

    return parseFloat(percentage.toFixed(2));
  }, [movementsByYear]);

  useEffect(() => {
    fetchMomentesByYear(2025);
  }, [fetchMomentesByYear]);

  return (
    <>
      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Income</h3>
        <p className="text-2xl font-semibold text-green-400">
          {format(movementsByYear.totalIncome)}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Primary</span>
          <span>Other</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Expenses</h3>
        <p className="text-2xl font-semibold text-red-500">
          {format(movementsByYear.totalOutcome)}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Fixed</span>
          <span>Variable</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Net Balance</h3>
        <p className="text-2xl font-semibold text-yellow-400">
          {format(movementsByYear.balance)}
        </p>
        <span className="text-xs text-gray-300 mt-1 block">
          {calculateVerticalPercentage()}% Expenses to Income Ratio
        </span>
      </div>

      <div className="md:col-span-3 dark:bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Income vs Expenses</h4>
        <div className="h-40 dark:bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Chart Placeholder
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded flex flex-col gap-4">
        <h4 className="text-sm text-gray-400 mb-2">Income by Category</h4>
        <div className="h-40 dark:bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Pie Chart Placeholder
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded flex flex-col gap-4">
        <h4 className="text-sm text-gray-400 mb-2">Expenses by Category</h4>
        <div className="h-40 dark:bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Bar Chart Placeholder
        </div>
      </div>

      <div className="flex flex-col gap-4 dark:bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Notes</h4>
        <textarea
          className="w-full h-40 dark:bg-gray-700 rounded p-2 text-white resize-none"
          placeholder="Add reminders, insights, or pending items about your finances here."
        />
        <div />
      </div>
    </>
  );
}

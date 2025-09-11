"use client";

import { useFormatCurrency } from "@/hooks/useFormatCurrency";
import { useMovementsStore } from "@/infrastructure/store/movements";
import { useMemo, useCallback, useEffect } from "react";
import LineChart from "@/ui/charts/LineChart";
import { MONTHS } from "@/utils/dates";

export default function Home() {
  const {
    movementsByYear,
    year,
    showAllYears,
    allYearsRaw,
    fetchMomentsByYear,
    fetchAllYearsData
  } = useMovementsStore();

  const values = useMemo(
    () => [
      {
        name: "Earned",
        data: movementsByYear.months?.map((m) => m.income) || [],
        color: "green"
      },
      {
        name: "Expense",
        data: movementsByYear.months?.map((m) => m.outcome) || [],
        color: "red"
      }
    ],
    [movementsByYear]
  );

  const allMonthsByYears = useCallback(
    (type: "income" | "outcome") =>
      allYearsRaw
        ?.map((m) => m.months)
        .flat()
        .map((m) => m[type]) || [],
    [allYearsRaw]
  );

  const allDataByYears = useCallback(
    (type: "totalIncome" | "totalOutcome" | "balance") =>
      allYearsRaw?.reduce((acc, m) => m[type] + acc, 0) || 0,
    [allYearsRaw]
  );

  const allYearsValues = useMemo(
    () => [
      {
        name: "Earned",
        data: allMonthsByYears("income"),
        color: "green"
      },
      {
        name: "Expense",
        data: allMonthsByYears("outcome"),
        color: "red"
      }
    ],
    [allYearsRaw, allMonthsByYears]
  );

  const format = useFormatCurrency();

  const calculateVerticalPercentage = useCallback((): number => {
    if (movementsByYear.totalIncome === 0) {
      return 0;
    }

    let balance = movementsByYear.balance;
    if (showAllYears) {
      balance = allDataByYears("balance");
    }

    const percentage = (balance / movementsByYear.totalIncome) * 100;

    return parseFloat(percentage.toFixed(2));
  }, [movementsByYear, showAllYears, allYearsRaw, allDataByYears]);

  useEffect(() => {
    fetchMomentsByYear(year);
  }, [fetchMomentsByYear, year]);

  useEffect(() => {
    fetchAllYearsData([2021, 2022, 2023, 2024, 2025]);
  }, [fetchAllYearsData, showAllYears]);

  return (
    <>
      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Income</h3>
        <p className="text-2xl font-semibold text-green-400">
          {format(
            showAllYears
              ? allDataByYears("totalIncome")
              : movementsByYear.totalIncome
          )}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Primary</span>
          <span>Other</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Expenses</h3>
        <p className="text-2xl font-semibold text-red-500">
          {format(
            showAllYears
              ? allDataByYears("totalOutcome")
              : movementsByYear.totalOutcome
          )}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Fixed</span>
          <span>Variable</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Net Balance</h3>
        <p className="text-2xl font-semibold text-yellow-400">
          {format(
            showAllYears ? allDataByYears("balance") : movementsByYear.balance
          )}
        </p>
        <span className="text-xs text-gray-300 mt-1 block">
          {calculateVerticalPercentage()}% Expenses to Income Ratio
        </span>
      </div>

      <div className="h-full md:col-span-3 dark:bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Income vs Expenses</h4>
        <div className="min-h-[300px] dark:bg-gray-700 rounded flex items-center justify-center text-gray-400">
          <LineChart
            xData={
              showAllYears
                ? allYearsRaw
                    ?.map((m) => m.months)
                    .flat()
                    .map(
                      (_, i) => `${MONTHS[i % 12]} ${2021 + Math.floor(i / 12)}`
                    )
                : MONTHS
            }
            series={showAllYears ? allYearsValues : values}
            height="calc(100dvh - 500px)"
          />
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

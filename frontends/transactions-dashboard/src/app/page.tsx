"use client";

import { useFormatCurrency } from "@/hooks/useFormatCurrency";
import { useMovementsStore } from "@/infrastructure/store/movements";
import { useMemo, useEffect } from "react";
import LineChart from "@/ui/charts/LineChart";
import { MONTHS } from "@/utils/dates";
import type {
  MovementByMonth,
  MovementByYear,
  MovementMonth,
  MovementYear
} from "@/core/entities/Movement";
import { cn } from "@/utils/styles";

type MovementsTotalTypes = "totalIncome" | "totalExpense" | "balance";
type MovementTypes = "income" | "outcome";

const getData = (data: MovementYear[] | MovementMonth[], type: MovementTypes) =>
  data?.map((m) => m[type]) || [];

const getTotal = (
  data: MovementByYear[] | MovementByMonth[],
  type: MovementsTotalTypes
) => data?.reduce((acc, m) => acc + m[type], 0) || 0;

export default function Home() {
  const {
    movementsByYear,
    movementsByMonth,
    year,
    month,
    timeSelected,
    institutionsSelected,
    allYearsRaw,
    fetchMomentsByYear,
    fetchMomentsByMonth,
    fetchAllYearsData,
    setInstitutionsSelected
  } = useMovementsStore();

  const format = useFormatCurrency();

  useEffect(() => {
    switch (timeSelected) {
      case "year":
        fetchMomentsByYear(year);
        break;
      case "month":
        fetchMomentsByMonth(year, month);
        break;
      case "all_years":
        fetchAllYearsData([2021, 2022, 2023, 2024, 2025]);
        break;
    }
  }, [
    timeSelected,
    year,
    month,
    fetchMomentsByYear,
    fetchMomentsByMonth,
    fetchAllYearsData,
    institutionsSelected
  ]);

  const chartData = useMemo(() => {
    let rawData: MovementYear[] | MovementMonth[] = [];
    let labels: string[] = [];
    let totalIncome = 0;
    let totalExpense = 0;
    let balance = 0;

    switch (timeSelected) {
      case "all_years":
        rawData = allYearsRaw?.flatMap((y) => y.months);
        labels = rawData?.map(
          (_, i) => `${MONTHS[i % 12]} ${2021 + Math.floor(i / 12)}`
        );
        totalIncome = getTotal(allYearsRaw, "totalIncome");
        totalExpense = getTotal(allYearsRaw, "totalExpense");
        balance = getTotal(allYearsRaw, "balance");
        break;
      case "year":
        const dataComplete = Array.from(
          {
            length: Math.max(...movementsByYear.months.map((m) => m.month ?? 0))
          },
          (_, index) => {
            const mov = movementsByYear.months.find(
              (m) => m.month === index + 1
            );

            return mov ? mov : { month: index + 1, income: 0, outcome: 0 };
          }
        );
        rawData = dataComplete;
        labels = MONTHS;
        totalIncome = movementsByYear.totalIncome;
        totalExpense = movementsByYear.totalExpense;
        balance = movementsByYear.balance;
        console.log(rawData);
        break;
      case "month":
        rawData = movementsByMonth.days;
        labels = Array.from({ length: 31 }, (_, index) => `${index + 1}`);
        totalIncome = movementsByMonth.totalIncome;
        totalExpense = movementsByMonth.totalExpense;
        balance = movementsByMonth.balance;
        break;
      default:
        rawData = [];
        labels = [];
    }

    const series = [
      { name: "Earned", data: getData(rawData, "income"), color: "green" },
      { name: "Expense", data: getData(rawData, "outcome"), color: "red" }
    ];

    const percentage =
      totalIncome > 0
        ? parseFloat(((balance / totalIncome) * 100).toFixed(2))
        : 0;

    return { series, labels, totalIncome, totalExpense, balance, percentage };
  }, [timeSelected, allYearsRaw, movementsByYear, movementsByMonth]);

  const { series, labels, totalIncome, totalExpense, balance, percentage } =
    chartData;

  return (
    <div className="grid md:grid-cols-3 grid-cols-1 justify-self-center  gap-8">
      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Total Income</h3>
        <p className="text-2xl font-semibold text-green-400">
          {format(totalIncome)}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Primary</span>
          <span>Other</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Total Expenses</h3>
        <p className="text-2xl font-semibold text-red-500">
          {format(totalExpense)}
        </p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Fixed</span>
          <span>Variable</span>
        </div>
      </div>

      <div className="dark:bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Net Balance</h3>
        <p className="text-2xl font-semibold text-yellow-400">
          {format(balance)}
        </p>
        <span className="text-xs text-gray-300 mt-1 block">
          {percentage}% Balance / Income Ratio
        </span>
      </div>

      <ul className="flex gap-2">
        {["davivienda", "manual"].map((i) => (
          <li
            key={i}
            className={cn(
              "flex items-center p-3 border border-gray-500 rounded-md text-gray-600 hover:bg-gray-200 dark:text-gray-300 dark:hover:bg-gray-700",
              institutionsSelected.includes(i)
                ? "flex items-center p-3 rounded-md bg-indigo-100 text-indigo-600 font-medium hover:bg-indigo-200 dark:bg-indigo-900 dark:text-indigo-400 dark:hover:bg-indigo-800"
                : ""
            )}
            onClick={() => {
              setInstitutionsSelected(
                institutionsSelected.includes(i)
                  ? institutionsSelected.filter((item) => item !== i)
                  : [...institutionsSelected, i]
              );
            }}
          >
            <span className="cursor-pointer">{i}</span>
          </li>
        ))}
      </ul>

      <div className="h-full md:col-span-3 dark:bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Income vs Expenses</h4>
        <div className="min-h-[300px] dark:bg-gray-700 rounded flex items-center justify-center text-gray-400">
          <LineChart
            xData={labels}
            series={series}
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
    </div>
  );
}

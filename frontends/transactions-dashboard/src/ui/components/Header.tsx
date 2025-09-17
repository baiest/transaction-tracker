"use client";

import { Time } from "@/infrastructure/store/models";
import { useMovementsStore } from "@/infrastructure/store/movements";
import { MONTHS } from "@/utils/dates";
import { usePathname } from "next/navigation";

export default function Header() {
  const { year, month, timeSelected, setYear, setMonth, setTimeSelected } =
    useMovementsStore();

  const pathname = usePathname();

  return (
    <>
      <header className="flex justify-between items-center">
        {pathname === "/" && (
          <div className="flex gap-4">
            <select
              className="dark:bg-gray-800 px-3 py-1 rounded"
              value={timeSelected}
              onChange={(e) => setTimeSelected(e.target.value as Time)}
            >
              <option value="all_years">All years</option>
              <option value="year">Year</option>
              <option value="month">Monthly</option>
            </select>
            {(timeSelected === "year" || timeSelected === "month") && (
              <select
                className="dark:bg-gray-800 px-3 py-1 rounded"
                value={year}
                onChange={(e) => setYear(Number(e.target.value))}
              >
                <option value="2025">2025</option>
                <option value="2024">2024</option>
                <option value="2023">2023</option>
                <option value="2022">2022</option>
                <option value="2021">2021</option>
                <option value="2021">2020</option>
                <option value="2021">2019</option>
              </select>
            )}

            {timeSelected === "month" && (
              <select
                className="dark:bg-gray-800 px-3 py-1 rounded"
                value={month}
                onChange={(e) => setMonth(Number(e.target.value))}
              >
                {MONTHS.map((m, i) => (
                  <option key={m} value={i}>
                    {m}
                  </option>
                ))}
              </select>
            )}
          </div>
        )}
        <div></div>
        <div className="flex gap-4">
          <button className="dark:bg-gray-800 px-3 py-1 rounded hover:bg-gray-700">
            Export
          </button>
          <button className="bg-green-600 px-4 py-1 rounded hover:bg-green-500">
            New Transaction
          </button>
        </div>
      </header>
    </>
  );
}

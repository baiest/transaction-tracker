"use client";

import { useMovementsStore } from "@/infrastructure/store/movements";

export default function Header() {
  const { year, showAllYears, setYear, setShowAllYears } = useMovementsStore();
  return (
    <>
      <h2 className="text-2xl">Transactions</h2>

      <header className="flex justify-between items-center">
        <div className="flex gap-4">
          <select
            className="dark:bg-gray-800 px-3 py-1 rounded"
            defaultValue="year"
            onChange={(e) => setShowAllYears(e.target.value === "all-years")}
          >
            <option value="all-years">All years</option>
            <option value="year">Year</option>
            <option value="month">Monthly</option>
          </select>
          {!showAllYears && (
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
        </div>
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

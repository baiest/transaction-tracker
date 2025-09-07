export default function Home() {
  return (
    <>
      <div className="bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Income</h3>
        <p className="text-2xl font-semibold text-green-400">$12,450</p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Primary</span>
          <span>Other</span>
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Expenses</h3>
        <p className="text-2xl font-semibold text-red-500">$8,320</p>
        <div className="flex gap-2 mt-2 text-xs text-gray-300">
          <span>Fixed</span>
          <span>Variable</span>
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded">
        <h3 className="text-sm text-gray-400">Net Balance</h3>
        <p className="text-2xl font-semibold text-yellow-400">$4,130</p>
        <span className="text-xs text-gray-300 mt-1 block">
          +6.2% vs last month
        </span>
      </div>

      <div className="md:col-span-3 bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Income vs Expenses</h4>
        <div className="h-40 bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Chart Placeholder
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded flex flex-col gap-4">
        <h4 className="text-sm text-gray-400 mb-2">Income by Category</h4>
        <div className="h-40 bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Pie Chart Placeholder
        </div>
      </div>

      <div className="bg-gray-800 p-4 rounded flex flex-col gap-4">
        <h4 className="text-sm text-gray-400 mb-2">Expenses by Category</h4>
        <div className="h-40 bg-gray-700 rounded flex items-center justify-center text-gray-400">
          Bar Chart Placeholder
        </div>
      </div>

      <div className="flex flex-col gap-4 bg-gray-800 p-4 rounded">
        <h4 className="text-sm text-gray-400 mb-2">Notes</h4>
        <textarea
          className="w-full h-40 bg-gray-700 rounded p-2 text-white resize-none"
          placeholder="Add reminders, insights, or pending items about your finances here."
        />
        <div />
      </div>
    </>
  );
}

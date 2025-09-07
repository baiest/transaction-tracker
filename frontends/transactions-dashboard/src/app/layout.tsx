import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";

import Sidebar from "@/components/Sidebar";

const poppins = Poppins({
  weight: ["400", "500", "600", "700"],
  subsets: ["latin"],
  variable: "--font-poppins"
});

export const metadata: Metadata = {
  title: "Track",
  description: "Transactions dashboard with information collected in emails"
};

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${poppins.variable} antialiased flex min-h-screen bg-white`}
      >
        <Sidebar />

        <div className="w-full min-h-screen bg-gray-900 text-white font-sans p-8 sm:p-5 grid grid-rows-[auto_auto_1fr] gap-8">
          <h2 className="text-2xl">Transactions</h2>
          <header className="flex justify-between items-center">
            <div className="flex gap-4">
              <select className="bg-gray-800 px-3 py-1 rounded">
                <option>Monthly</option>
                <option>Weekly</option>
              </select>
              <select className="bg-gray-800 px-3 py-1 rounded">
                <option>Q1 2025</option>
              </select>
              <select className="bg-gray-800 px-3 py-1 rounded">
                <option>All Accounts</option>
              </select>
            </div>
            <div className="flex gap-4">
              <button className="bg-gray-800 px-3 py-1 rounded hover:bg-gray-700">
                Export
              </button>
              <button className="bg-green-600 px-4 py-1 rounded hover:bg-green-500">
                New Transaction
              </button>
            </div>
          </header>
          <main className="w-full max-w-[1700px] grid grid-cols-1 justify-self-center md:grid-cols-3 gap-8">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}

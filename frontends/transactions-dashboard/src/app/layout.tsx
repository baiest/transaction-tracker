import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";

import Sidebar from "@/ui/components/Sidebar";
import Header from "@/ui/components/Header";

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

        <div className="w-full min-h-screen dark:bg-gray-900 text-white font-sans p-8 sm:p-5 grid grid-rows-[auto_auto_1fr] gap-8">
          <Header />
          <main className="w-full max-w-[1700px] grid grid-cols-1 justify-self-center md:grid-cols-3 gap-8">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}

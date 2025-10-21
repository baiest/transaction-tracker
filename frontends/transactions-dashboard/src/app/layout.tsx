import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";

import Sidebar from "@/ui/components/Sidebar";
import Header from "@/ui/components/Header";
import ThemeProvider from "@/ui/components/ThemeProvider";
import { Toaster } from "@/ui/components/Sonner";

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
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${poppins.variable} antialiased flex min-h-screen bg-white`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <Sidebar />

          <div className="w-full min-h-screen dark:bg-gray-900 text-white font-sans p-8 sm:p-5 grid grid-rows-[auto_auto_1fr] gap-8">
            <Header />
            <main className="w-full max-w-[1700px] mx-auto">{children}</main>
          </div>
          <Toaster richColors />
        </ThemeProvider>
      </body>
    </html>
  );
}

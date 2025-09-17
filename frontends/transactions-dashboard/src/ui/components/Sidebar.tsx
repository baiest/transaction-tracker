"use client";

import { cn } from "@/utils/styles";
import { usePathname } from "next/navigation";
import React from "react";
import Link from "next/link";

interface SidebarProps {
  className?: string;
}

const optionStyle =
  "flex items-center p-3 rounded-md text-gray-600 hover:bg-gray-200 dark:text-gray-300 dark:hover:bg-gray-700";

const optionSelectedStyle =
  "flex items-center p-3 rounded-md bg-indigo-100 text-indigo-600 font-medium hover:bg-indigo-200 dark:bg-indigo-900 dark:text-indigo-400 dark:hover:bg-indigo-800";

const options = [
  { href: "/", label: "Dashboard" },
  { href: "/movements", label: "Movimientos" },
  { href: "/accounts", label: "Cuentas" },
  { href: "/goals", label: "Metas" },
  { href: "/settings", label: "Ajustes" }
];
export default function Sidebar({ className }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "w-64 bg-gray-100 border-r border-gray-200 dark:bg-gray-800 dark:border-gray-700",
        className
      )}
    >
      <nav className="sticky top-0 flex-1 p-6">
        <div className="text-xl font-bold mb-8 text-black dark:text-white">
          Finanzas
        </div>
        <ul className="space-y-2">
          {options.map(({ href, label }) => (
            <li key={href}>
              <Link
                href={href}
                className={cn(
                  optionStyle,
                  pathname === href && optionSelectedStyle
                )}
              >
                {label}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}

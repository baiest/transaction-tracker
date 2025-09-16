"use client";

import { useFormatCurrency } from "@/hooks/useFormatCurrency";
import { useMovementsStore } from "@/infrastructure/store/movements";
import { useCallback, useEffect, useMemo } from "react";
import { cn } from "@/utils/styles";
import { TrendingDown, TrendingUp } from "lucide-react";
import { DataTable } from "@/ui/components/DataTable";
import type { ColumnDef } from "@tanstack/react-table";
import type { Movement } from "@/core/entities/Movement";
import { useRouter, useSearchParams } from "next/navigation";

export default function Movements() {
  const { movements, fetchMovements, totalPages } = useMovementsStore();
  const format = useFormatCurrency();
  const router = useRouter();
  const searchParams = useSearchParams();

  const pageFromUrl = Number(searchParams.get("page") ?? "1");
  const pageIndex = pageFromUrl - 1;

  const pageSize = 10;

  const columns = useMemo<ColumnDef<Movement>[]>(
    () => [
      { accessorKey: "detail", header: "Details" },
      { accessorKey: "date", header: "Date" },
      {
        accessorKey: "value",
        header: "Amount",
        cell: ({ row }) => {
          const amount: number = row.getValue("value");
          const isNegative = row.original.isNegative;

          return (
            <div
              className={cn(
                "flex gap-2",
                isNegative ? "text-red-300" : "text-green-300"
              )}
            >
              {isNegative ? <TrendingDown /> : <TrendingUp />} {format(amount)}
            </div>
          );
        }
      }
    ],
    [format]
  );

  const handlePageChange = useCallback(
    (p: number) => router.push(`/movements?page=${p + 1}`),
    [router]
  );

  useEffect(() => {
    fetchMovements(pageIndex + 1);
  }, [fetchMovements, pageIndex]);

  return (
    <DataTable
      columns={columns}
      data={Object.values(movements)}
      totalPages={totalPages}
      pageIndex={pageIndex}
      pageSize={pageSize}
      onPageChange={handlePageChange}
    />
  );
}

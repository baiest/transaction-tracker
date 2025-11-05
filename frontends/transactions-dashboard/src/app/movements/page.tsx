"use client";

import { Suspense } from "react";
import MovementsContent from "./MovementsContent";

export const dynamic = "force-dynamic";

export default function MovementsPage() {
  return (
    <Suspense fallback={<div>Loading movements...</div>}>
      <MovementsContent />
    </Suspense>
  );
}

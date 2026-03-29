"use client";

import ProtectedRoute, { useUser } from "@/components/ProtectedRoute";
import DashboardPage from "./Dashboard";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function DashboardPageWrapper() {
  const router = useRouter();

  const { user } = useUser();

  useEffect(() => {
    if (user === null) {
      router.replace("/login");
    }
  }, [user, router]);

  if (!user) return null;

  return (
    <ProtectedRoute>
      <DashboardPage />
    </ProtectedRoute>
  );
}

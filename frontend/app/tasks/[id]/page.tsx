"use client";

import ProtectedRoute, { useUser } from "@/components/ProtectedRoute";
import TaskDetailPage from "./TaskDetailPage";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function TaskDetailPageWrapper() {
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
      <TaskDetailPage />
    </ProtectedRoute>
  );
}

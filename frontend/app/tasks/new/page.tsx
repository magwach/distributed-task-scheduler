"use client";

import ProtectedRoute, { useUser } from "@/components/ProtectedRoute";
import CreateTaskPage from "./NewTask";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function NewTaskPageWrapper() {
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
      <CreateTaskPage />
    </ProtectedRoute>
  );
}
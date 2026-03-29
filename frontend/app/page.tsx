"use client";

import { useRouter } from "next/navigation";
import ProtectedRoute, { useUser } from "@/components/ProtectedRoute";
import { useEffect } from "react";

export default function Home() {
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
      <></>
    </ProtectedRoute>
  );
}

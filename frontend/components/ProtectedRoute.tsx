"use client";
import {
  useEffect,
  useState,
  ReactNode,
  createContext,
  useContext,
} from "react";
import { useRouter } from "next/navigation"; // ← Change this line
import { fetchUser } from "../services/api";
import { User } from "@/services/types";

interface ProtectedRouteProps {
  children: ReactNode;
}

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  loading: boolean;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function checkAuth() {
      try {
        const data = await fetchUser();
        setUser(data);
      } catch {
        router.replace("/login");
      } finally {
        setLoading(false);
      }
    }

    checkAuth();
  }, [router]);

  if (loading) {
    return (
      <div className="auth-layout">
        <div className="grid-bg" />
        <div
          style={{
            textAlign: "center",
            animation: "fadeInUp 0.4s ease both",
            marginLeft: "240px",
          }}
        >
          <div
            style={{
              width: 40,
              height: 40,
              border: "2px solid var(--border)",
              borderTop: "2px solid var(--accent-cyan)",
              borderRadius: "50%",
              animation: "spin 0.8s linear infinite",
              margin: "0 auto 20px",
            }}
          />
        </div>

        <style>{`
        @keyframes pulse-glow {
          0%, 100% { box-shadow: 0 0 32px var(--accent-cyan-glow); }
          50% { box-shadow: 0 0 48px rgba(0, 212, 255, 0.5); }
        }
        @keyframes progress-fill {
          from { width: 0%; }
          to { width: 100%; }
        }
      `}</style>
      </div>
    );
  }

  return (
    <UserContext.Provider value={{ user, setUser, loading }}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const context = useContext(UserContext);
  if (!context) throw new Error("useUser must be used within a UserProvider");
  return context;
}

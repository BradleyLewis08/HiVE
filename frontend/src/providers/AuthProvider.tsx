import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  useCallback,
  useMemo,
} from "react";
import User, { SystemRole } from "../models/User";
import { jwtDecode } from "jwt-decode";
import { AuthToken } from "../common/types";
import APIClient from "../api/APIClient";

interface AuthProviderProps {
  children: React.ReactNode;
}

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  handleLogout: () => void;
  handleRoleChange: (systemRole: SystemRole) => void;
  initializeAuth: (token: string) => void;
}

const initialAuthContextValue: AuthContextValue = {
  user: null,
  isAuthenticated: false,
  handleLogout: () => {},
  handleRoleChange: () => {},
  initializeAuth: () => {},
};

const AuthContext = createContext<AuthContextValue>(initialAuthContextValue);

interface AuthProviderProps {
  children: React.ReactNode;
}

async function getUser(netId: string): Promise<User> {
  const client = new APIClient();
  return await client.getUser(netId);
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isInitializing, setIsInitializing] = useState(true);

  // Memoize the refreshUserFromToken function since it's used in multiple callbacks
  const refreshUserFromToken = useCallback(async (token: string) => {
    const { data: tokenData } = jwtDecode(token) as AuthToken;
    const currentUser = await getUser(tokenData.netId);
    setUser(currentUser);
  }, []);

  const handleLogout = useCallback(() => {
    localStorage.removeItem("token");
    setUser(null);
  }, []);

  console.log(user);

  // Memoize initializeAuth
  const initializeAuth = useCallback(
    async (token: string) => {
      try {
        localStorage.setItem("token", token);
        await refreshUserFromToken(token);
      } catch (error) {
        console.error("Failed to initialize auth with token:", error);
        handleLogout();
        throw error;
      }
    },
    [refreshUserFromToken, handleLogout]
  );

  // Memoize handleRoleChange
  const handleRoleChange = useCallback(
    async (systemRole: SystemRole) => {
      if (!user) return;

      try {
        const client = new APIClient();
        const newUser = await client.updateUserPermissions(
          user.netId,
          systemRole
        );
        setUser(newUser);
      } catch (error) {
        console.error("Error updating user role:", error);
        const token = localStorage.getItem("token");
        if (token) {
          await refreshUserFromToken(token);
        }
      }
    },
    [user, refreshUserFromToken]
  );

  useEffect(() => {
    async function loadInitialAuth() {
      const token = localStorage.getItem("token");
      if (!token) {
        setIsInitializing(false);
        return;
      }

      try {
        await refreshUserFromToken(token);
      } catch (error) {
        console.error("Failed to initialize auth:", error);
        handleLogout();
      } finally {
        setIsInitializing(false);
      }
    }

    loadInitialAuth();
  }, [refreshUserFromToken, handleLogout]); // Now these dependencies are stable

  const value = useMemo(
    () => ({
      user,
      isAuthenticated: user !== null,
      handleLogout,
      handleRoleChange,
      initializeAuth,
    }),
    [user, handleLogout, handleRoleChange, initializeAuth]
  );

  if (isInitializing) {
    return <div>Loading...</div>;
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}


import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { api } from "@/lib/api-client";
import { toast } from "@/components/ui/use-toast";

export interface User {
  id: string;
  name: string;
  username: string;
  email: string;
  bio?: string;
  avatar?: string;
  location?: string;
  website?: string;
  createdAt: string;
  updatedAt: string;
  followers?: number;
  following?: number;
  isFollowed?: boolean;
}

interface AuthResponse {
  token: string;
  user: User;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  updateUser: (data: UpdateUserData) => Promise<void>;
}

interface AuthProviderProps {
  children: ReactNode;
}

interface RegisterData {
  name: string;
  username: string;
  email: string;
  password: string;
}

interface UpdateUserData {
  name?: string;
  bio?: string;
  avatar?: string;
  location?: string;
  website?: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  
  // Check if user is logged in on mount
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      fetchCurrentUser();
    } else {
      setIsLoading(false);
    }
  }, []);
  
  const fetchCurrentUser = async () => {
    try {
      const userData = await api.get<User>("/users/me");
      setUser(userData);
    } catch (error) {
      // If token is invalid, remove it
      localStorage.removeItem("token");
    } finally {
      setIsLoading(false);
    }
  };
  
  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const data = await api.post<AuthResponse>("/auth/login", { email, password });
      localStorage.setItem("token", data.token);
      setUser(data.user);
      toast({
        title: "Welcome back!",
        description: `Logged in as ${data.user.name}`,
      });
    } finally {
      setIsLoading(false);
    }
  };
  
  const register = async (data: RegisterData) => {
    setIsLoading(true);
    try {
      const response = await api.post<AuthResponse>("/auth/register", data);
      localStorage.setItem("token", response.token);
      setUser(response.user);
      toast({
        title: "Account created!",
        description: "You're now logged in to your new account.",
      });
    } finally {
      setIsLoading(false);
    }
  };
  
  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    toast({
      title: "Logged out",
      description: "You've been successfully logged out.",
    });
  };
  
  const updateUser = async (data: UpdateUserData) => {
    if (!user) return;
    
    try {
      const updatedUser = await api.put<User>(`/users/${user.id}`, data);
      setUser(updatedUser);
      toast({
        title: "Profile updated",
        description: "Your profile has been updated successfully.",
      });
    } catch (error) {
      console.error("Failed to update user:", error);
    }
  };
  
  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};


import { createContext, useContext, useState, ReactNode, useEffect } from "react";
import { api } from "@/lib/api-client";
import { toast } from "@/components/ui/use-toast";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export interface User {
  id: string;
  name: string;
  username: string;
  email: string;
  bio?: string;
  avatar?: string;
  cover?: string;
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
  changePassword: (currentPassword: string, newPassword: string) => Promise<void>;
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
  const queryClient = useQueryClient()
  const { data: user, isLoading: isFetchingUser } = useQuery<User>({
    queryKey: ['auth'],
    queryFn: () => api.get<User>("/users/me"),
    enabled: !!localStorage.getItem('token')
  });
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setIsLoading(isFetchingUser)
  }, [isFetchingUser])

  const login = async (email: string, password: string) => {
    try {
      const data = await api.post<AuthResponse>("/auth/login", { email, password });
      localStorage.setItem("token", data.token);
      queryClient.setQueryData(['auth'], () => {
        return data.user
      });
      toast({
        title: "Welcome back!",
        description: `Logged in as ${data.user.name}`,
      });
    } finally {
      setIsLoading(false);
    }
  }

  const register = async (data: RegisterData) => {
    setIsLoading(true);
    try {
      const response = await api.post<AuthResponse>("/auth/register", data);
      localStorage.setItem("token", response.token);
      queryClient.setQueryData(['auth'], () => {
        return response.user
      });
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
    queryClient.setQueryData(['auth'], () => {
      return null
    });
    toast({
      title: "Logged out",
      description: "You've been successfully logged out.",
    });
  };

  const updateUser = async (data: UpdateUserData) => {
    if (!user) return;

    try {
      const updatedUser = await api.put<User>(`/users/${user.id}`, data);
      queryClient.setQueryData(['auth'], () => {
        return updatedUser
      });
      toast({
        title: "Profile updated",
        description: "Your profile has been updated successfully.",
      });
    } catch (error) {
      console.error("Failed to update user:", error);
    }
  };

  const changePassword = async (currentPassword: string, newPassword: string) => {
    if (!user) return;

    try {
      await api.put('/auth/change-password', {
        current_password: currentPassword,
        new_password: newPassword
      });

      toast({
        title: "Password updated",
        description: "Your password has been changed successfully.",
      });
    } catch (error) {
      console.error("Failed to change password:", error);
      throw error;
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
        changePassword,
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

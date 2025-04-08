
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { User } from "@/contexts/AuthContext";
import { toast } from "@/components/ui/use-toast";

export const useUser = (userId: string) => {
  return useQuery<User>({
    queryKey: ["user", userId],
    queryFn: () => api.get<User>(`/users/${userId}`),
    enabled: !!userId,
  });
};

export const useUserByUsername = (username: string) => {
  return useQuery<User>({
    queryKey: ["user", "username", username],
    queryFn: () => api.get<User>(`/users/username/${username}`),
    enabled: !!username,
  });
};

export const useFollowUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation<void, Error, string>({
    mutationFn: (userId: string) => api.post<void>(`/users/follow/${userId}`),
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: ["user", userId] });
      queryClient.invalidateQueries({ queryKey: ["users"] });
      queryClient.invalidateQueries({ queryKey: ["followers"] });
      queryClient.invalidateQueries({ queryKey: ["following"] });
      
      toast({
        title: "Success",
        description: "You are now following this user",
      });
    },
  });
};

export const useUnfollowUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation<void, Error, string>({
    mutationFn: (userId: string) => api.delete<void>(`/users/follow/${userId}`),
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: ["user", userId] });
      queryClient.invalidateQueries({ queryKey: ["users"] });
      queryClient.invalidateQueries({ queryKey: ["followers"] });
      queryClient.invalidateQueries({ queryKey: ["following"] });
      
      toast({
        title: "Success",
        description: "You have unfollowed this user",
      });
    },
  });
};

export const useFollowers = (limit = 10) => {
  return useQuery<User[]>({
    queryKey: ["followers", limit],
    queryFn: () => api.get<User[]>(`/users/followers?limit=${limit}`),
  });
};

export const useFollowing = (limit = 10) => {
  return useQuery<User[]>({
    queryKey: ["following", limit],
    queryFn: () => api.get<User[]>(`/users/following?limit=${limit}`),
  });
};

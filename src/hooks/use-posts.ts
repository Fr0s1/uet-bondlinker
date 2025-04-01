
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { toast } from "@/components/ui/use-toast";
import { User } from "@/contexts/AuthContext";

export interface Post {
  id: string;
  user_id: string;
  author?: User;
  content: string;
  image?: string;
  created_at: string;
  updated_at: string;
  likes: number;
  comments: number;
  is_liked?: boolean;
}

export interface CreatePostData {
  content: string;
  image?: string;
}

export const usePosts = (limit = 10) => {
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();
  
  const { data, isLoading, error } = useQuery({
    queryKey: ["posts", page, limit],
    queryFn: () => api.get<Post[]>(`/posts?limit=${limit}&offset=${(page - 1) * limit}`),
  });
  
  const createPost = useMutation({
    mutationFn: (newPost: CreatePostData) => api.post<Post>("/posts", newPost),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      toast({
        title: "Post created",
        description: "Your post has been published successfully!",
      });
    },
  });
  
  const likePost = useMutation({
    mutationFn: (postId: string) => api.post<void>(`/posts/${postId}/like`),
    onSuccess: (_, postId) => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["post", postId] });
    },
  });
  
  const unlikePost = useMutation({
    mutationFn: (postId: string) => api.delete<void>(`/posts/${postId}/like`),
    onSuccess: (_, postId) => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["post", postId] });
    },
  });
  
  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
    createPost,
    likePost,
    unlikePost,
  };
};

export const usePost = (postId: string) => {
  return useQuery({
    queryKey: ["post", postId],
    queryFn: () => api.get<Post>(`/posts/${postId}`),
    enabled: !!postId,
  });
};

export const useFeed = (limit = 10) => {
  const [page, setPage] = useState(1);
  
  const { data, isLoading, error } = useQuery({
    queryKey: ["feed", page, limit],
    queryFn: () => api.get<Post[]>(`/posts/feed?limit=${limit}&offset=${(page - 1) * limit}`),
  });
  
  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

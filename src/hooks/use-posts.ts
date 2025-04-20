
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { toast } from "@/components/ui/use-toast";

export interface Post {
  id: string;
  userId: string;
  author?: {
    id: string;
    name: string;
    username: string;
    avatar?: string;
  };
  content: string;
  image?: string;
  createdAt: string;
  updatedAt: string;
  likes: number;
  comments: number;
  shares: number;
  isLiked?: boolean;
  sharedPostId?: string;
  sharedPost?: Post;
}

export interface CreatePostData {
  content: string;
  image?: string;
}

export interface SharePostData {
  content: string;
}

export const usePosts = (userId?: string, limit = 10) => {
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();

  // Calculate offset based on page and limit
  const offset = (page - 1) * limit;

  // Construct query key based on parameters
  const queryKey = userId
    ? ["posts", userId, page, limit]
    : ["posts", page, limit];

  // Construct endpoint URL based on parameters
  const endpoint = userId
    ? `/posts?userId=${userId}&limit=${limit}&offset=${offset}`
    : `/posts?limit=${limit}&offset=${offset}`;

  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: queryKey,
    queryFn: () => api.get<Post[]>(endpoint),
  });

  const createPost = useMutation<Post, Error, CreatePostData>({
    mutationFn: (postData: CreatePostData) => api.post<Post>('/posts', postData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["feed"] });
      toast({
        title: "Post created",
        description: "Your post has been created successfully!",
      });
    },
  });

  const likePost = useMutation<{ likes: number }, Error, string>({
    mutationFn: (postId: string) => api.post<{ likes: number }>(`/posts/${postId}/like`),
    onSuccess: (data, postId) => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["feed"] });
      queryClient.invalidateQueries({ queryKey: ["trending"] });
    },
  });

  const unlikePost = useMutation<{ likes: number }, Error, string>({
    mutationFn: (postId: string) => api.delete<{ likes: number }>(`/posts/${postId}/like`),
    onSuccess: (data, postId) => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["feed"] });
      queryClient.invalidateQueries({ queryKey: ["trending"] });
    },
  });

  const sharePost = useMutation<Post, Error, { postId: string, content: string }>({
    mutationFn: ({ postId, content }) => api.post<Post>(`/posts/${postId}/share`, { content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["feed"] });
      queryClient.invalidateQueries({ queryKey: ["trending"] });
      toast({
        title: "Post shared",
        description: "The post has been shared successfully!",
      });
    },
  });

  const deletePost = useMutation<void, Error, string>({
    mutationFn: (postId: string) => api.delete<void>(`/posts/${postId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["feed"] });
      queryClient.invalidateQueries({ queryKey: ["trending"] });
      toast({
        title: "Post deleted",
        description: "Your post has been deleted successfully!",
      });
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
    sharePost,
    deletePost,
  };
};

export const useFeed = (limit = 10) => {
  const [page, setPage] = useState(1);
  const offset = (page - 1) * limit;

  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["feed", page, limit],
    queryFn: () => api.get<Post[]>(`/posts/feed?limit=${limit}&offset=${offset}`),
  });

  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useTrendingPosts = (limit = 10) => {
  const [page, setPage] = useState(1);
  const offset = (page - 1) * limit;

  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["trending", page, limit],
    queryFn: () => api.get<Post[]>(`/posts/trending?limit=${limit}&offset=${offset}`),
  });

  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useSuggestedPosts = (limit = 10) => {
  const [page, setPage] = useState(1);
  const offset = (page - 1) * limit;

  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["suggested", page, limit],
    queryFn: () => api.get<Post[]>(`/posts/suggested?limit=${limit}&offset=${offset}`),
  });

  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const usePost = (postId: string) => {
  return useQuery<Post>({
    queryKey: ["post", postId],
    queryFn: () => api.get<Post>(`/posts/${postId}`),
    enabled: !!postId,
  });
};

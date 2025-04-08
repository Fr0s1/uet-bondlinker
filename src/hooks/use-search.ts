
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { User } from "@/contexts/AuthContext";
import { Post } from "./use-posts";

export interface SearchResults {
  users: User[];
  posts: Post[];
}

export const useSearch = (query: string, limit = 10) => {
  const [page, setPage] = useState(1);
  
  const { data, isLoading, error } = useQuery<SearchResults>({
    queryKey: ["search", query, page, limit],
    queryFn: () => api.get<SearchResults>(`/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: query.length > 0,
  });
  
  return {
    results: data || { users: [], posts: [] },
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useSearchUsers = (query: string, limit = 10) => {
  const [page, setPage] = useState(1);
  
  const { data, isLoading, error } = useQuery<User[]>({
    queryKey: ["search", "users", query, page, limit],
    queryFn: () => api.get<User[]>(`/search/users?q=${encodeURIComponent(query)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: query.length > 0,
  });
  
  return {
    users: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useSearchPosts = (query: string, limit = 10) => {
  const [page, setPage] = useState(1);
  
  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["search", "posts", query, page, limit],
    queryFn: () => api.get<Post[]>(`/search/posts?q=${encodeURIComponent(query)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: query.length > 0,
  });
  
  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

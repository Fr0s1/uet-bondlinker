
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { User } from "@/contexts/AuthContext";
import { Post } from "./use-posts";
import { useDebounce } from "./use-debounce";

export interface SearchResults {
  users: User[];
  posts: Post[];
}

export const useSearch = (query: string, limit = 10) => {
  const [page, setPage] = useState(1);
  const debouncedQuery = useDebounce(query, 300);
  
  const { data, isLoading, error } = useQuery<SearchResults>({
    queryKey: ["search", debouncedQuery, page, limit],
    queryFn: () => api.get<SearchResults>(`/search?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0,
    staleTime: 60000, // 1 minute
    keepPreviousData: true, // Keep old data while fetching new
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
  const debouncedQuery = useDebounce(query, 300);
  
  const { data, isLoading, error } = useQuery<User[]>({
    queryKey: ["search", "users", debouncedQuery, page, limit],
    queryFn: () => api.get<User[]>(`/search/users?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0,
    staleTime: 60000,
    keepPreviousData: true,
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
  const debouncedQuery = useDebounce(query, 300);
  
  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["search", "posts", debouncedQuery, page, limit],
    queryFn: () => api.get<Post[]>(`/search/posts?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0,
    staleTime: 60000,
    keepPreviousData: true,
  });
  
  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

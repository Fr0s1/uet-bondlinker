
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api, ApiResponse } from "@/lib/api-client";
import { User } from "@/contexts/AuthContext";
import { Post } from "./use-posts";
import { useDebounce } from "./use-debounce";

export interface SearchResults {
  users: User[];
  posts: Post[];
}

export const useSearch = ({ query, limit = 10, enabled }: { query: string, enabled: boolean, limit?: number }) => {
  const [page, setPage] = useState(1);
  const debouncedQuery = useDebounce(query, 300);

  const { data, isLoading, error } = useQuery<SearchResults>({
    queryKey: ["search", debouncedQuery, page, limit],
    queryFn: () => api.get<SearchResults>(`/search?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0 && enabled,
    staleTime: 60000, // 1 minute
    placeholderData: { users: [], posts: [] }, // Use placeholderData instead of keepPreviousData
  });

  return {
    results: data || { users: [], posts: [] },
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useSearchUsers = ({ query, limit = 10, enabled }: { query: string, enabled: boolean, limit?: number }) => {
  const [page, setPage] = useState(1);
  const debouncedQuery = useDebounce(query, 300);

  const { data, isLoading, error } = useQuery<User[]>({
    queryKey: ["search", "users", debouncedQuery, page, limit],
    queryFn: () => api.get<User[]>(`/search/users?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0 && enabled,
    staleTime: 60000,
    placeholderData: [], // Use placeholderData instead of keepPreviousData
  });

  return {
    users: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

export const useSearchPosts = ({ query, limit = 10, enabled }: { query: string, enabled: boolean, limit?: number }) => {
  const [page, setPage] = useState(1);
  const debouncedQuery = useDebounce(query, 300);

  const { data, isLoading, error } = useQuery<Post[]>({
    queryKey: ["search", "posts", debouncedQuery, page, limit],
    queryFn: () => api.get<Post[]>(`/search/posts?q=${encodeURIComponent(debouncedQuery)}&limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: debouncedQuery.length > 0 && enabled,
    staleTime: 60000,
    placeholderData: [], // Use placeholderData instead of keepPreviousData
  });

  return {
    posts: data || [],
    isLoading,
    error,
    page,
    setPage,
  };
};

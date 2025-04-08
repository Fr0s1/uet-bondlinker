
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { User } from "@/contexts/AuthContext";
import { toast } from "@/components/ui/use-toast";

export interface Comment {
  id: string;
  user_id: string;
  author?: User;
  post_id: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCommentData {
  content: string;
}

export const useComments = (postId: string, limit = 10) => {
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();
  
  const { data, isLoading, error } = useQuery<Comment[]>({
    queryKey: ["comments", postId, page, limit],
    queryFn: () => api.get<Comment[]>(`/posts/${postId}/comments?limit=${limit}&offset=${(page - 1) * limit}`),
    enabled: !!postId,
  });
  
  const createComment = useMutation<Comment, Error, CreateCommentData>({
    mutationFn: (comment: CreateCommentData) => api.post<Comment>(`/posts/${postId}/comments`, comment),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
      queryClient.invalidateQueries({ queryKey: ["posts"] }); // To update comment count
      toast({
        title: "Comment added",
        description: "Your comment has been added successfully!",
      });
    },
  });
  
  interface UpdateCommentParams {
    commentId: string;
    content: string;
  }
  
  const updateComment = useMutation<Comment, Error, UpdateCommentParams>({
    mutationFn: ({ commentId, content }: UpdateCommentParams) => 
      api.put<Comment>(`/posts/comments/${commentId}`, { content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
      toast({
        title: "Comment updated",
        description: "Your comment has been updated successfully!",
      });
    },
  });
  
  const deleteComment = useMutation<void, Error, string>({
    mutationFn: (commentId: string) => api.delete<void>(`/posts/comments/${commentId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
      queryClient.invalidateQueries({ queryKey: ["posts"] }); // To update comment count
      toast({
        title: "Comment deleted",
        description: "Your comment has been deleted successfully!",
      });
    },
  });
  
  return {
    comments: data || [],
    isLoading,
    error,
    page,
    setPage,
    createComment,
    updateComment,
    deleteComment,
  };
};

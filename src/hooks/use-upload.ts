
import { useMutation } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { toast } from '@/components/ui/use-toast';

interface UploadResponse {
  url: string;
  filename: string;
}

export const useFileUpload = () => {
  return useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('file', file);
      
      try {
        toast({
          title: "Uploading...",
          description: "Your file is being uploaded to the server.",
        });
        
        const response = await api.post<UploadResponse>('/uploads', formData, true);
        
        toast({
          title: "Upload complete",
          description: "Your file has been uploaded successfully.",
        });
        
        return response;
      } catch (error) {
        toast({
          title: "Upload failed",
          description: "Failed to upload your file. Please try again.",
          variant: "destructive",
        });
        throw error;
      }
    }
  });
};

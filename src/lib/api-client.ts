
import { toast } from "@/components/ui/use-toast";

const API_URL = "http://localhost:8080/api/v1";

// Error handling utility
const handleError = (error: unknown) => {
  console.error("API Error:", error);
  
  if (error instanceof Response) {
    return error.json().then(data => {
      toast({
        title: "Error",
        description: data.error || "Something went wrong",
        variant: "destructive",
      });
      throw data;
    });
  }
  
  toast({
    title: "Error",
    description: "Network error. Please try again later.",
    variant: "destructive",
  });
  throw error;
};

// Generic fetch function with authorization header
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {},
  isFormData: boolean = false
): Promise<T> {
  const token = localStorage.getItem("token");
  
  const headers: HeadersInit = {
    ...options.headers,
  };
  
  if (!isFormData) {
    headers["Content-Type"] = "application/json";
  }
  
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  
  try {
    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    });
    
    if (!response.ok) {
      return handleError(response) as Promise<T>;
    }
    
    // Handle empty responses (like for DELETE operations)
    if (response.status === 204) {
      return {} as T;
    }
    
    return await response.json();
  } catch (error) {
    return handleError(error) as Promise<T>;
  }
}

// API client object with methods for common operations
export const api = {
  get: <T>(endpoint: string) => fetchApi<T>(endpoint, { method: "GET" }),
  
  post: <T>(endpoint: string, data?: any, isFormData: boolean = false) =>
    fetchApi<T>(
      endpoint, 
      {
        method: "POST",
        body: isFormData ? data : data ? JSON.stringify(data) : undefined,
      },
      isFormData
    ),
  
  put: <T>(endpoint: string, data: any, isFormData: boolean = false) =>
    fetchApi<T>(
      endpoint, 
      {
        method: "PUT",
        body: isFormData ? data : JSON.stringify(data),
      },
      isFormData
    ),
  
  delete: <T>(endpoint: string) =>
    fetchApi<T>(endpoint, { method: "DELETE" }),
};

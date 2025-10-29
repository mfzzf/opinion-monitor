import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://165.154.98.129:8080';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors and token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, add to queue
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return api(originalRequest);
        }).catch((err) => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = localStorage.getItem('refresh_token');

      if (refreshToken) {
        try {
          const response = await axios.post(`${API_URL}/api/auth/refresh`, {
            refresh_token: refreshToken,
          });

          const { access_token, refresh_token: newRefreshToken } = response.data;

          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', newRefreshToken);

          processQueue(null, access_token);

          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return api(originalRequest);
        } catch (refreshError) {
          processQueue(refreshError, null);

          // Refresh failed, clear tokens and redirect to login
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');

          if (typeof window !== 'undefined') {
            window.location.href = '/login';
          }

          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      } else {
        // No refresh token, clear everything and redirect
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');

        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
      }
    }

    return Promise.reject(error);
  }
);

export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface Video {
  id: number;
  user_id: number;
  original_filename: string;
  file_path: string;
  cover_path?: string;
  file_size: number;
  duration: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
}

export interface Report {
  id: number;
  video_id: number;
  cover_text: string;
  transcript_text?: string;
  sentiment_score: number;
  sentiment_label: string;
  key_topics: string;
  risk_level: string;
  detailed_analysis: string;
  recommendations: string;
  processing_time: number;
  created_at: string;
}

export interface Job {
  id: number;
  video_id: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  retry_count: number;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

// Auth APIs
export const authAPI = {
  register: (data: { username: string; email: string; password: string }) =>
    api.post('/api/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/api/auth/login', data),
  logout: () => api.post('/api/auth/logout'),
  refresh: (refreshToken: string) =>
    axios.post(`${API_URL}/api/auth/refresh`, { refresh_token: refreshToken }),
  me: () => api.get('/api/auth/me'),
};

// Video APIs
export const videoAPI = {
  upload: (formData: FormData, onUploadProgress?: (progressEvent: any) => void) =>
    api.post('/api/videos/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress,
    }),
  list: (params?: { page?: number; page_size?: number; status?: string }) =>
    api.get('/api/videos', { params }),
  get: (id: number) => api.get(`/api/videos/${id}`),
  delete: (id: number) => api.delete(`/api/videos/${id}`),
};

// Report APIs
export const reportAPI = {
  getByVideoId: (videoId: number) => api.get(`/api/reports/${videoId}`),
  list: (params?: { page?: number; page_size?: number }) =>
    api.get('/api/reports', { params }),
};

// Job APIs
export const jobAPI = {
  getStatus: (id: number) => api.get(`/api/jobs/${id}/status`),
  list: (params?: { page?: number; page_size?: number }) =>
    api.get('/api/jobs', { params }),
};


// src/utils/request.js
import axios from "axios";
import { generateRandomString } from "./index";

// API基础URL
const BASE_URL = import.meta.env.VITE_IS_DOCKER ? "" : "http://localhost:8080";


// 创建Axios实例
const instance = axios.create({
  baseURL: BASE_URL, // 使用配置的API基础URL
  timeout: 30000, // 请求超时时间
  headers: {
    "Content-Type": "application/json",
    "X-Request-ID": `${generateRandomString(12)}`,
  },
});


instance.interceptors.request.use(
  (config) => {
    // 添加JWT token认证
    const token = localStorage.getItem('weknora_token');
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    
    config.headers["X-Request-ID"] = `${generateRandomString(12)}`;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Token刷新标志，防止多个请求同时刷新token
let isRefreshing = false;
let failedQueue: Array<{ resolve: Function; reject: Function }> = [];
let hasRedirectedOn401 = false;

// 处理队列中的请求
const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else {
      resolve(token);
    }
  });
  
  failedQueue = [];
};

instance.interceptors.response.use(
  (response) => {
    // 根据业务状态码处理逻辑
    const { status, data } = response;
    if (status === 200 || status === 201) {
      return data;
    } else {
      return Promise.reject(data);
    }
  },
  async (error: any) => {
    const originalRequest = error.config;
    
    if (!error.response) {
      return Promise.reject({ message: "网络错误，请检查您的网络连接" });
    }
    
    // 如果是登录接口的401，直接返回错误以便页面展示toast，不做跳转
    if (error.response.status === 401 && originalRequest?.url?.includes('/auth/login')) {
      const { status, data } = error.response;
      return Promise.reject({ status, message: (typeof data === 'object' ? data?.message : data) || '用户名或密码错误' });
    }

    // 如果是401错误且不是刷新token的请求，尝试刷新token
    if (error.response.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh')) {
      if (isRefreshing) {
        // 如果正在刷新token，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers['Authorization'] = 'Bearer ' + token;
          return instance(originalRequest);
        }).catch(err => {
          return Promise.reject(err);
        });
      }
      
      originalRequest._retry = true;
      isRefreshing = true;
      
      const refreshToken = localStorage.getItem('weknora_refresh_token');
      
      if (refreshToken) {
        try {
          // 动态导入refresh token API
          const { refreshToken: refreshTokenAPI } = await import('../api/auth/index');
          const response = await refreshTokenAPI(refreshToken);
          
          if (response.success && response.data) {
            const { token, refreshToken: newRefreshToken } = response.data;
            
            // 更新localStorage中的token
            localStorage.setItem('weknora_token', token);
            localStorage.setItem('weknora_refresh_token', newRefreshToken);
            
            // 更新请求头
            originalRequest.headers['Authorization'] = 'Bearer ' + token;
            
            // 处理队列中的请求
            processQueue(null, token);
            
            return instance(originalRequest);
          } else {
            throw new Error(response.message || 'Token刷新失败');
          }
        } catch (refreshError) {
          // 刷新失败，清除所有token并跳转到登录页
          localStorage.removeItem('weknora_token');
          localStorage.removeItem('weknora_refresh_token');
          localStorage.removeItem('weknora_user');
          localStorage.removeItem('weknora_tenant');
          
          processQueue(refreshError, null);
          
          // 跳转到登录页
          if (!hasRedirectedOn401 && typeof window !== 'undefined') {
            hasRedirectedOn401 = true;
            window.location.href = '/login';
          }
          
          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      } else {
        // 没有refresh token，直接跳转到登录页
        localStorage.removeItem('weknora_token');
        localStorage.removeItem('weknora_user');
        localStorage.removeItem('weknora_tenant');
        
        if (!hasRedirectedOn401 && typeof window !== 'undefined') {
          hasRedirectedOn401 = true;
          window.location.href = '/login';
        }
        
        return Promise.reject({ message: '请重新登录' });
      }
    }
    
    const { status, data } = error.response;
    // 将HTTP状态码一并抛出，方便上层判断401等场景
    return Promise.reject({ status, ...(typeof data === 'object' ? data : { message: data }) });
  }
);

export function get(url: string) {
  return instance.get(url);
}

export async function getDown(url: string) {
  let res = await instance.get(url, {
    responseType: "blob",
  });
  return res
}

export function postUpload(url: string, data = {}) {
  return instance.post(url, data, {
    headers: {
      "Content-Type": "multipart/form-data",
      "X-Request-ID": `${generateRandomString(12)}`,
    },
  });
}

export function postChat(url: string, data = {}) {
  return instance.post(url, data, {
    headers: {
      "Content-Type": "text/event-stream;charset=utf-8",
      "X-Request-ID": `${generateRandomString(12)}`,
    },
  });
}

export function post(url: string, data = {}) {
  return instance.post(url, data);
}

export function put(url: string, data = {}) {
  return instance.put(url, data);
}

export function del(url: string) {
  return instance.delete(url);
}

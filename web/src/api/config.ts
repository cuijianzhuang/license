import axios, {AxiosRequestConfig} from 'axios';
import {message} from 'antd';
import {ApiResponse} from '../types';

// 创建 axios 实例
const axiosInstance = axios.create({
  baseURL: '/api', // 在实际部署时，可以通过环境变量或配置文件设置
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

const skipUrls= [
    "/jetbrains/generate",
    "/jetbrains/licenseServerRule",
]

// 提供一个通用的API请求函数
const api = {
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axiosInstance.get(url, config);

      let length = skipUrls.filter(item => url.includes(item)).length;
      if (length > 0) {
        // 直接返回接口响应内容
        return response.data as T;
      }
      
      const apiResponse = response.data as ApiResponse<T>;
      
      if (apiResponse.code !== 200) {
        message.error(apiResponse.message || '请求失败');
        throw new Error(apiResponse.message || '请求失败');
      }
      
      return apiResponse.data;
    } catch (error: any) {
      // 处理错误
      const errorMsg = error.response?.data?.message || '网络错误，请稍后重试';
      message.error(errorMsg);
      throw error;
    }
  },
  
  async post<T = any>(url: string, data?: any, config?: any): Promise<T> {
    try {
      const response = await axiosInstance.post(url, data, config);

      let length = skipUrls.filter(item => url.includes(item)).length;
      if (length > 0) {
        // 直接返回接口响应内容
        return response.data as T;
      }

      const apiResponse = response.data as ApiResponse<T>;
      
      if (apiResponse.code !== 200) {
        message.error(apiResponse.message || '请求失败');
        throw new Error(apiResponse.message || '请求失败');
      }
      
      return apiResponse.data;
    } catch (error: any) {
      // 处理错误
      const errorMsg = error.response?.data?.message || '网络错误，请稍后重试';
      message.error(errorMsg);
      throw error;
    }
  }
};

export default api; 
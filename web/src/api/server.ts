import axios from 'axios';

interface ServerVersion {
  version: string;
}

/**
 * 获取服务器版本信息
 * @returns Promise<string> 包含版本号的Promise
 */
export const getVersion = async (): Promise<string> => {
  try {
    const response = await axios.get<ServerVersion>('/api/server/version');
    return response.data.version;
  } catch (error) {
    console.error('Failed to fetch server version:', error);
    return '';
  }
}; 
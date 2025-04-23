import axios from 'axios';

interface ServerVersionResponse {
  version: string;
  needUpdate: boolean;
  latestVersion?: string;
}

/**
 * 版本信息接口
 */
export interface VersionInfo {
  version: string;
  needUpdate: boolean;
  latestVersion?: string;
}

/**
 * 获取服务器版本信息
 * @returns Promise<VersionInfo> 包含版本信息的Promise
 */
export const getVersion = async (): Promise<VersionInfo> => {
  try {
    const response = await axios.get<ServerVersionResponse>('/api/server/version');
    return {
      version: response.data.version,
      needUpdate: response.data.needUpdate,
      latestVersion: response.data.latestVersion
    };
  } catch (error) {
    console.error('Failed to fetch server version:', error);
    return {
      version: '',
      needUpdate: false
    };
  }
}; 
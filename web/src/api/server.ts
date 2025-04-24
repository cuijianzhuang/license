import api from "./config";

interface ServerVersionResponse {
  version: string;
  hash: string;
  arch: string;
  needUpdate: boolean;
  latestVersion?: string;
}

/**
 * 版本信息接口
 */
export interface VersionInfo {
  version: string;
  hash: string;
  arch: string;
  needUpdate: boolean;
  latestVersion?: string;
}

/**
 * 获取服务器版本信息
 * @param queryParam 可选的查询参数，用于避免缓存
 * @returns Promise<VersionInfo> 包含版本信息的Promise
 */
export const getVersion = async (queryParam?: string): Promise<VersionInfo> => {
  try {
    const url = `/server/version${queryParam || ''}`;
    const response =  await api.get<ServerVersionResponse>(url, {
      headers: {
        'Cache-Control': 'no-cache, no-store',
        'Pragma': 'no-cache',
        'Expires': '0'
      }
    });
    return {
      version: response.version,
      hash: response.hash,
      arch: response.arch,
      needUpdate: response.needUpdate,
      latestVersion: response.latestVersion
    };
  } catch (error) {
    console.error('Failed to fetch server version:', error);
    return {
      version: '',
      hash: '',
      arch: '',
      needUpdate: false
    };
  }
}; 
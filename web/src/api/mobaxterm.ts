import api from './config';
import { MobaXtermLicense } from '../types';

/**
 * 生成MobaXterm许可证
 * @param username 用户名
 * @param version 版本号
 */
export const generateLicense = async (
  username: string, 
  version: string
): Promise<MobaXtermLicense> => {
  return api.post<MobaXtermLicense>('/mobaxterm/generate', { 
    username,
    version
  });
}; 
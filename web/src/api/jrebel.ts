import api from './config';
import { JRebelLicense } from '../types';

/**
 * 生成JRebel许可证
 * @param username 用户名
 * @param email 邮箱
 * @param teamName 团队名称
 */
export const generateLicense = async (
  username: string,
  email: string,
  teamName: string
): Promise<JRebelLicense> => {
  // 注意：JRebel后端API可能与其他不同，需要根据实际API路径和参数调整
  return api.post<JRebelLicense>('/jrebel/leases', {
    username,
    email,
    teamName
  });
}; 
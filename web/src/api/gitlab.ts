import api from './config';
import { GitLabLicense } from '../types';

/**
 * 生成GitLab许可证
 * @param company 公司名称
 * @param email 邮箱
 * @param userCount 用户数量
 * @param expiresAt 过期时间
 */
export const generateLicense = async (
  company: string, 
  email: string, 
  userCount: number, 
  expiresAt: string
): Promise<GitLabLicense> => {
  return api.post<GitLabLicense>('/gitlab/generate', {
    company,
    email,
    userCount,
    expiresAt
  });
}; 
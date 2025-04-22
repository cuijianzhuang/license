import api from './config';
import { FinalShellLicense } from '../types';

/**
 * 生成FinalShell许可证
 * @param username 用户名
 */
export const generateLicense = async (username: string): Promise<FinalShellLicense> => {
  return api.post<FinalShellLicense>('/final-shell/generateLicense', { username });
}; 
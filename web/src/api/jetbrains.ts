import api from './config';
import {ApiResponse, JetBrainsLicense} from '../types';

/**
 * 生成JetBrains激活码
 * @param licenseeName 授权用户名
 * @param effectiveDate 有效日期，可选，格式为: YYYY-MM-DD HH:MM:SS
 * @param codes 产品代码，多个代码用逗号分隔
 */
export const generateLicense = async (
  licenseeName?: string, 
  effectiveDate?: string,
  codes?: string
): Promise<ApiResponse<JetBrainsLicense>> => {
  const params: any = {};
  
  if (licenseeName) params.licenseeName = licenseeName;
  if (effectiveDate) params.effectiveDate = effectiveDate;
  if (codes) params.codes = codes;
  
  return api.get<ApiResponse<JetBrainsLicense>>('/jetbrains/generate', {params: params});
};

/**
 * 获取JetBrains服务器规则
 */
export const getLicenseServerRule = async (): Promise<string> => {
  return api.get<string>('/jetbrains/licenseServerRule');
}; 
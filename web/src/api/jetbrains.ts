import api from './config';
import { JetBrainsLicense } from '../types';

/**
 * 获取JetBrains产品列表
 * 注意：这个接口会触发后端更新，应谨慎调用
 */
export const fetchProductList = async (): Promise<string[]> => {
  return api.get<string[]>('/jetbrains/product/fetchLatest');
};

/**
 * 获取JetBrains插件列表
 * 注意：这个接口会触发后端更新，应谨慎调用
 */
export const fetchPluginList = async (): Promise<string[]> => {
  return api.get<string[]>('/jetbrains/plugin/fetchLatest');
};

/**
 * 手动触发更新JetBrains产品和插件列表
 * 注意：我们尝试调用这两个接口来触发更新，它们在后端执行异步更新
 */
export const updateProductsAndPlugins = async (): Promise<{message: string}> => {
  // 同时触发产品和插件的更新
  await Promise.all([
    api.get<string[]>('/jetbrains/product/fetchLatest'),
    api.get<string[]>('/jetbrains/plugin/fetchLatest')
  ]);
  
  return { message: "更新请求已发送，正在后台处理" };
};

/**
 * 手动触发仅更新JetBrains插件列表
 * 注意：这个接口会触发后端更新，应谨慎调用
 */
export const updatePluginsOnly = async (): Promise<{message: string}> => {
  // 只触发插件的更新
  await api.get<string[]>('/jetbrains/plugin/fetchLatest');
  
  return { message: "插件更新请求已发送，正在后台处理" };
};

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
): Promise<JetBrainsLicense> => {
  const params: any = {};
  
  if (licenseeName) params.licenseeName = licenseeName;
  if (effectiveDate) params.effectiveDate = effectiveDate;
  if (codes) params.codes = codes;
  
  return api.get<JetBrainsLicense>('/jetbrains/generate', params);
};

/**
 * 获取JetBrains服务器规则
 */
export const getLicenseServerRule = async (): Promise<string> => {
  return api.get<string>('/jetbrains/licenseServerRule');
}; 
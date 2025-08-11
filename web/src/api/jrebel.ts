import {v4} from 'uuid';

/**
 * 生成随机GUID
 */
export const generateGuid = (): string => {
  return v4();
}; 
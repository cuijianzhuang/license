import api from './config';
import { MobaXtermLicense } from '../types';

/**
 * Generate MobaXterm license
 * @param formData Form data containing username, version, count
 * @returns File stream for the license
 */
export const generateLicense = async (
  formData: FormData
): Promise<Blob> => {
  return api.post('/mobaxterm/generate', formData, {
    responseType: 'blob',
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}; 
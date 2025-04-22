import api from './config';
import axios from 'axios';

/**
 * 生成GitLab许可证
 * @param name 姓名
 * @param email 邮箱
 * @param company 公司名称
 * @param expireTime 过期时间
 * @returns 返回一个Promise，成功时为true
 */
export const generateLicense = async (
  name: string, 
  email: string, 
  company: string, 
  expireTime: string
): Promise<boolean> => {
  try {
    // 使用FormData格式发送请求
    const formData = new FormData();
    formData.append('Name', name);
    formData.append('Email', email);
    formData.append('Company', company);
    formData.append('ExpireTime', expireTime);

    // 使用axios直接请求，以便可以设置responseType为blob
    const response = await axios({
      method: 'post',
      url: '/api/gitlab/generate',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      responseType: 'blob' // 指定响应类型为blob
    });

    // 验证响应是否包含有效的blob数据
    if (response.data && response.data.size > 0) {
      // 创建Blob并生成下载链接
      const blob = new Blob([response.data], { 
        type: response.headers['content-type'] || 'application/octet-stream' 
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      
      // 从Content-Disposition获取文件名，如果没有则使用默认名称
      const contentDisposition = response.headers['content-disposition'];
      let filename = 'gitlab-license.key';
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
        if (filenameMatch && filenameMatch[1]) {
          filename = filenameMatch[1];
        }
      }
      
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      return true;
    } else {
      throw new Error('返回的文件数据无效');
    }
  } catch (error) {
    console.error('生成GitLab许可证失败:', error);
    throw error;
  }
}; 
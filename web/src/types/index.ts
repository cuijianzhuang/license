// 通用响应类型
export interface ApiResponse<T = any> {
  status: number;
  message: string;
  data: T;
}

// JetBrains 激活码类型
export interface JetBrainsLicense {
  id?: number;
  licenseId?: string;
  code: string;
  product?: string;
  createdAt?: string;
  serverConfig?: string;
  activationCode?: string;
  powerConf?: string;
}

// GitLab 许可证类型
export interface GitLabLicense {
  id?: number;
  license: string;
  company?: string;
  email?: string;
  expiresAt?: string;
  userCount?: number;
  createdAt?: string;
}

// FinalShell 许可证类型
export interface FinalShellLicense {
  id?: number;
  license: string;
  username?: string;
  createdAt?: string;
}

// MobaXterm 许可证类型
export interface MobaXtermLicense {
  id?: number;
  license: string;
  username?: string;
  version?: string;
  createdAt?: string;
}

// JRebel 许可证类型
export interface JRebelLicense {
  id?: number;
  username?: string;
  email?: string;
  teamName?: string;
  validKey?: string;
  createdAt?: string;
}

// 分页请求参数
export interface PaginationParams {
  page: number;
  limit: number;
}

// 分页响应
export interface PaginatedResponse<T> {
  total: number;
  items: T[];
} 
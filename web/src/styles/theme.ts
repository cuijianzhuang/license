import { ThemeConfig } from 'antd';

// 定义全局主题
export const theme: ThemeConfig = {
  token: {
    colorPrimary: '#2563eb',
    borderRadius: 8,
    colorBgBase: '#ffffff',
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
  },
  components: {
    Layout: {
      bodyBg: '#f5f7fa',
      headerBg: '#ffffff',
    },
    Menu: {
      itemHoverBg: '#f0f5ff',
      itemSelectedBg: '#e6f4ff',
      itemSelectedColor: '#2563eb',
    },
    Card: {
      colorBorderSecondary: '#e5e7eb',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05)',
    },
    Button: {
      colorPrimaryHover: '#3a75eb',
      borderRadius: 6,
    },
    Select: {
      borderRadius: 6,
    },
    Input: {
      borderRadius: 6,
    },
    DatePicker: {
      borderRadius: 6,
    },
  },
};

// 媒体查询断点
export const breakpoints = {
  xs: '480px',
  sm: '576px',
  md: '768px',
  lg: '992px',
  xl: '1200px',
  xxl: '1600px',
};

// 定义一个响应式函数，用于在不同屏幕尺寸应用不同样式
export const responsive = {
  isMobile: () => window.innerWidth < 768,
  isTablet: () => window.innerWidth >= 768 && window.innerWidth < 992,
  isDesktop: () => window.innerWidth >= 992,
}; 
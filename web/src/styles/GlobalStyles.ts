import {createGlobalStyle} from 'styled-components';
import {breakpoints} from './theme';

const GlobalStyles = createGlobalStyle`
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  html, body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
      'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
      'Noto Color Emoji';
    font-size: 14px;
    line-height: 1.6;
    color: #1f2937;
    background-color: #f5f7fa;
  }

  #root {
    min-height: 100vh;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 16px;
  }

  h1, h2, h3, h4, h5, h6 {
    margin-bottom: 0.5em;
    font-weight: 600;
    line-height: 1.4;
    color: #111827;
  }

  p {
    margin-bottom: 1em;
    color: #374151;
  }

  a {
    text-decoration: none;
    color: #2563eb;
    transition: color 0.2s;
    
    &:hover {
      color: #3a75eb;
    }
  }

  .ant-form-item-label > label {
    font-weight: 500;
    color: #4b5563;
  }

  .ant-alert {
    border-radius: 8px;
  }

  .ant-input::placeholder, 
  .ant-select-selection-placeholder {
    color: #9ca3af;
  }

  .ant-typography {
    color: #374151;
  }

  /* 响应式样式 */
  @media (max-width: ${breakpoints.md}) {
    html, body {
      font-size: 14px;
    }
  }

  @media (max-width: ${breakpoints.sm}) {
    html, body {
      font-size: 13px;
    }
  }

  /* 自定义滚动条 */
  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }

  ::-webkit-scrollbar-track {
    background: #f1f1f1;
  }

  ::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
  }
`;

export default GlobalStyles; 
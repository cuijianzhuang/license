import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import reportWebVitals from './reportWebVitals';
// Import i18n configuration
import './i18n';

// 在生产环境下禁用所有控制台日志
if (process.env.NODE_ENV === 'production') {
  // 保存原始控制台方法的引用，以便必要时可以恢复
  const originalConsole = {
    log: console.log,
    info: console.info,
    warn: console.warn,
    error: console.error,
    debug: console.debug,
    trace: console.trace
  };

  // 创建一个函数来拦截并识别日志来源
  const createConsoleProxy = (method: string, originalMethod: (...args: any[]) => void) => {
    return function(...args: any[]) {
      // 在生产环境中完全禁用日志输出
      // 临时注释掉下面的 return 可以帮助你识别日志来源
      return;

      // 下面的代码在调试时可以取消注释，帮助识别日志来源
      // try {
      //   throw new Error('Console trace');
      // } catch (e) {
      //   originalMethod.apply(console, [
      //     '拦截到日志调用:',
      //     ...args,
      //     '\n来源堆栈:',
      //     e.stack
      //   ]);
      // }
    };
  };

  // 使用代理替换控制台方法
  console.log = createConsoleProxy('log', originalConsole.log);
  console.info = createConsoleProxy('info', originalConsole.info);
  console.warn = createConsoleProxy('warn', originalConsole.warn);
  console.debug = createConsoleProxy('debug', originalConsole.debug);
  console.trace = createConsoleProxy('trace', originalConsole.trace);
  
  // 保留错误日志，通常这些对于问题排查很重要
  // 如果需要也禁用错误日志，取消下面一行的注释
  // console.error = createConsoleProxy('error', originalConsole.error);
  
  // 控制台的其他方法如 table, time, timeEnd 等也可以在此禁用
  
  // 防止第三方库使用 window.console
  if (typeof window !== 'undefined') {
    window.console = console;
  }
}

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

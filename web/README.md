# 许可证管理服务前端

这是许可证管理服务的前端项目，基于React开发，提供各类软件的激活码生成界面。

## 功能特点

- 现代化、响应式界面设计，支持PC和移动端
- 支持多种软件的激活码生成：
  - JetBrains系列产品
  - GitLab企业版
  - FinalShell
  - MobaXterm
  - JRebel
- 简洁直观的用户界面
- 一键复制和下载激活信息

## 技术栈

- React 18
- TypeScript
- React Router
- Ant Design
- Styled Components
- Axios
- Node.js 22.14.0+
- npm 10.9.0+

## 安装和使用

### 开发环境

1. 安装依赖：

```bash
npm install
```

2. 启动开发服务器：

```bash
npm start
```

应用将在 [http://localhost:3000](http://localhost:3000) 启动。

### 生产环境构建

```bash
npm run build
```

构建后的文件将生成在 `build` 目录中。

## 项目结构

```
src/
├── api/          # API服务和请求处理
├── assets/       # 静态资源文件
├── components/   # 共用组件
├── layouts/      # 布局组件
├── pages/        # 页面组件
├── styles/       # 全局样式和主题
├── types/        # TypeScript类型定义
├── utils/        # 工具函数
├── App.tsx       # 应用入口
└── index.tsx     # 渲染入口
```

## API代理配置

开发环境下，API请求会通过代理转发到后端服务。代理配置在 `src/setupProxy.js` 文件中：

```js
const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
      pathRewrite: {
        '^/api': '',
      },
    })
  );
};
```

如需修改后端服务地址，请更新 `target` 配置。

## 部署

1. 执行生产环境构建：

```bash
npm run build
```

2. 将 `build` 目录中的文件部署到Web服务器

可以使用Nginx或Apache等Web服务器部署，并配置合适的反向代理以转发API请求到后端服务。

## 贡献

如有问题或建议，请提交Issue或Pull Request。

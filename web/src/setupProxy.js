const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:15000',
      changeOrigin: true,
      pathRewrite: {
        '^/api': '', // 将 /api 替换为空字符串
      },
    })
  );
}; 
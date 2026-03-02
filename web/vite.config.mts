import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const target = env.VITE_API_PROXY_TARGET || `http://localhost:${env.HTTP_PORT || '15000'}`;
  const port = Number(env.PORT || 3000);

  return {
    plugins: [react()],
    server: {
      host: true,
      port,
      open: false,
      proxy: {
        '/api': {
          target,
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: 'build',
      sourcemap: false,
    },
  };
});

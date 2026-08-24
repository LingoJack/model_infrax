import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// 构建产物输出到 ../dist（go:embed 嵌入二进制）
// 开发模式：npm run dev 后访问 vite 端口，/api 代理到本地 jen ui 服务
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://127.0.0.1:8899',
    },
  },
})

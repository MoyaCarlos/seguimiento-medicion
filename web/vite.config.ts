import react from '@vitejs/plugin-react'
import { defineConfig, type ProxyOptions } from 'vite'

const apiProxy: ProxyOptions = {
  target: 'http://localhost:8080',
  bypass(request) {
    if (request.headers.accept?.includes('text/html')) return '/index.html'
  },
}

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/backlog': apiProxy,
      '/projects': apiProxy,
      '/sprints': apiProxy,
      '/demo': apiProxy,
    },
  },
})

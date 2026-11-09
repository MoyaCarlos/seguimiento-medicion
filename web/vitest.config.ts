import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
    // Vitest no detiene una prueba vencida: sigue actuando sobre la pantalla de la siguiente.
    testTimeout: 30_000,
    css: { include: /features\/proyecto\/estilos\.css/ },
  },
})

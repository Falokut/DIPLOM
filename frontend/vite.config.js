import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    host: true,
    port: 8080,
    strictPort: true,
    cors: false,
    hmr: false,
  },
  preview: {
    host: process.env.SERVER_HOST,
    port: process.env.SERVER_PORT,
    strictPort: true,
    cors: process.env.SERVER_ENABLE_CORS,
    hmr: false,
  }
})
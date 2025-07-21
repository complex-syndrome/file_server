import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, '../'); // <-- load from ../.env
	return {
		define: {
			'import.meta.env.VITE_BACKEND_PORT': env.VITE_BACKEND_PORT
		},
		plugins: [tailwindcss(), sveltekit()],
		server: {
			proxy: {
				'/ws': {
					target: `ws://localhost:${env.VITE_BACKEND_PORT}`,
					ws: true,
					changeOrigin: true
				},
				'/api': {
					target: `http://localhost:${env.VITE_BACKEND_PORT}`,
					changeOrigin: true
				}
			}
		}
	};
});

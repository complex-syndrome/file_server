import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import type { HttpProxy } from 'vite';

// Proxy with header to get real ip of request
function config(proxy: HttpProxy.Server) {
	proxy.on('proxyReq', (proxyReq, req) => {
		proxyReq.removeHeader('X-Forwarded-For');
		proxyReq.setHeader('X-Forwarded-For', req.socket.remoteAddress || '');
	});
}

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, '../'); // <-- load from ../.env

	return {
		define: {
			'import.meta.env.VITE_CUSTOM_VALUE': JSON.stringify(env.VITE_CUSTOM_VALUE)
		},
		plugins: [tailwindcss(), sveltekit()],
		server: {
			// Proxy to backend
			proxy: {
				'/ws': {
					target: `ws://localhost:${env.VITE_BACKEND_PORT}`,
					ws: true,
					changeOrigin: true,
					configure: config
				},
				'/api': {
					target: `http://localhost:${env.VITE_BACKEND_PORT}`,
					changeOrigin: true,
					configure: config
				}
			}
		}
	};
});

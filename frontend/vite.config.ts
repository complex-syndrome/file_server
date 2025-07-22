import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import type { HttpProxy, ProxyOptions } from 'vite';

function config(proxy: HttpProxy.Server, _options: ProxyOptions) {
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

import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, '../'); // <-- load from ../.env
	return {
		define: {
			'import.meta.env.VITE_WS_URL': JSON.stringify(env.VITE_WS_URL),
			'import.meta.env.VITE_API_URL': JSON.stringify(env.VITE_API_URL)
		},
		plugins: [tailwindcss(), sveltekit()]
	};
});

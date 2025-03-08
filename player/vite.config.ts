import fs from 'fs';
import tailwindcss from '@tailwindcss/vite';
import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd());
	let allowedHosts: string[] = [];
	if (env.VITE_ALLOWED_SERVERS) {
		allowedHosts = env.VITE_ALLOWED_SERVERS.split(',');
	}
	console.log(`Allowed hosts: ${allowedHosts}`);

	return {
		plugins: [sveltekit(), tailwindcss()],

		test: {
			workspace: [
				{
					extends: './vite.config.ts',
					plugins: [svelteTesting()],

					test: {
						name: 'client',
						environment: 'jsdom',
						clearMocks: true,
						include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
						exclude: ['src/lib/server/**'],
						setupFiles: ['./vitest-setup-client.ts']
					}
				},
				{
					extends: './vite.config.ts',

					test: {
						name: 'server',
						environment: 'node',
						include: ['src/**/*.{test,spec}.{js,ts}'],
						exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
					}
				}
			]
		},

		server: {
			https: {
				key: fs.readFileSync(`./certs/${env.VITE_DOMAIN}.key.pem`),
				cert: fs.readFileSync(`./certs/${env.VITE_DOMAIN}.pem`)
			},
			allowedHosts,
			hmr: {
				clientPort: 5173
			},
			host: '0.0.0.0',
			port: 5173
		}
	};
});

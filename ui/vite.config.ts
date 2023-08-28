import {defineConfig} from 'vite';
import {svelte} from '@sveltejs/vite-plugin-svelte';
import path from 'path';

export default defineConfig({
    plugins: [
        svelte(),
    ],
    build: {
        sourcemap: true,
    },
    server: {
        proxy: {
            '^/service/.*': {
                target: 'http://localhost:8080',
            },
        },
    },
});

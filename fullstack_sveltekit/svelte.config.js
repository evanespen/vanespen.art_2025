import adapter from '@sveltejs/adapter-node';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: preprocess({}),
    vitePlugin: {
        experimental: {
            useVitePreprocess: true,
        },
    },
    kit: {
        adapter: adapter({out: 'build'}),

        // // Override http methods in the Todo forms
        // methodOverride: {
        //     allowed: ['PATCH', 'DELETE']
        // },

        alias: {
            '$src': './src/',
            '$lib': './src/lib/',
        },
    }
};

export default config;

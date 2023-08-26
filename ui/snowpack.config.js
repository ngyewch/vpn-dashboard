module.exports = {
    mount: {
        public: '/',
        src: '/_dist_',
    },
    plugins: [
        '@snowpack/plugin-svelte',
        '@snowpack/plugin-dotenv',
        [
            '@snowpack/plugin-run-script',
            {cmd: 'svelte-check --output human', watch: '$1 --watch', output: 'stream'},
        ],
    ],
    devOptions: {
        port: 5000,
        open: "chrome",
    },
    proxy: {
        "/service/": "http://localhost:8080/service/"
    },
};

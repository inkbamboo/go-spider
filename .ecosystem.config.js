const os = process.env.OS ? process.env.OS : "darwin";
const watch = process.env.WATCH ? process.env.WATCH : false;
const env = process.env.ENV ? process.env.ENV : 'dev';
const go_env = process.env.ENV ? process.env.ENV : 'dev';
const configs = {
    instances: 1,
    max_memory_restart: "1000M",
    autorestart: false,
    max_restarts: 500,
    exec_mode: "fork",
    interpreter: "none",
    log_file: '../../logs/go-spider/go-spider.log',
    combine_logs: true,
    watch,
    env_uat: {
        env: 'dev',
        GO_ENV: go_env
    },
    env_stage: {
        env: 'tm',
        GO_ENV: go_env
    },
    env_prod: {
        env: 'prod',
        GO_ENV: go_env
    }
}
module.exports = {
    apps: [
        {
            name: "ke-sjz-area-spider",
            cwd: `./dist/housespider_${os}_amd64_v1`,
            script: "housespider",
            args:`--env ${env} --platform ke --city sjz --spider area hs`,
            ...configs,
            max_restarts: 10,
        },
        {
            name: "ke-sjz-ershou-spider",
            cwd: `./dist/housespider_${os}_amd64_v1`,
            script: "housespider",
            args:`--env ${env} --platform ke --city sjz --spider ershou hs`,
            ...configs,
            max_restarts: 10,
        },
        {
            name: "ke-sjz-chengjiao-spider",
            cwd: `./dist/housespider_${os}_amd64_v1`,
            script: "housespider",
            args:`--env ${env} --platform ke --city sjz --spider chengjiao hs`,
            ...configs,
            max_restarts: 10,
        },
        {
            name: "zhsc-poetry-spider",
            cwd: `./dist/poetryspider_${os}_amd64_v1`,
            script: "poetryspider",
            args:`--env ${env} --platform zhsc --spider poetry  ps`,
            ...configs,
            max_restarts: 10,
        }
    ],
};
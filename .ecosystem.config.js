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
        env: 'test',
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
            cwd: `./dist/kespider_${os}_amd64_v1`,
            script: "kespider",
            args:`--env ${env} --city sjz --spider area ke`,
            ...configs,
            max_restarts: 10,
        },
        {
            name: "ke-sjz-ershou-spider",
            cwd: `./dist/kespider_${os}_amd64_v1`,
            script: "kespider",
            args:`--env ${env}  --city sjz --spider ershou ke`,
            ...configs,
            max_restarts: 10,
        },
        {
            name: "ke-sjz-chengjiao-spider",
            cwd: `./dist/kespider_${os}_amd64_v1`,
            script: "kespider",
            args:`--env ${env} --city sjz --spider chengjiao ke`,
            ...configs,
            max_restarts: 10,
        }
    ],
};
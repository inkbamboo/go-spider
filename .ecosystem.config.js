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
    log_file: '../../logs/boxserver/boxserver.log',
    combine_logs: true,
    watch,
    env_uat: {
        env: 'uat',
        GO_ENV: go_env
    },
    env_stage: {
        env: 'stage',
        GO_ENV: go_env
    },
    env_sit: {
        env: 'sit',
        GO_ENV: go_env
    },
    env_stress: {
        env: 'stress',
        GO_ENV: go_env
    },
    env_prod: {
        env: 'prod',
        GO_ENV: go_env
    },
    env_preprod: {
        env: 'preprod',
        GO_ENV: go_env
    }
}
module.exports = {
    apps: [
        {
            name: "boxconfidence-box-server",
            cwd: `./build/box-server_${os}_amd64`,
            script: "boxconfidence-box-server",
            args:`--env ${env} s`,
            ...configs,
            max_restarts: 10,
        }
    ],
};
alias d := dev
alias s := start

# 开发服务器
dev:
    export IS_DEV=true && watchexec -e go -r "go run main.go"

# 启动服务器(不会监听文件变化重启服务器)
start:
    export IS_DEV=true && go run main.go

# 数据库迁移 & 数据库填充
migrate:
    rm -rf ./database.sqlite && go run main.go --migrate && go run main.go --seed

# 重新生成 swagger 文档
docs:
    rm -rf ./docs && swag init

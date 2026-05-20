## 介绍

学习和示例 fiber + viper + zap + swagger + gorm + jwt + asynq 的一个示例练习项目

请注意安装 Go 开发环境和 [just](https://github.com/casey/just)

## 启动

```sh
git clone https://github.com/liaohui5/fiber-curd-api
cd fiber-curd-api
go mod tidy

just dev #启动开发服务器
```

打开浏览器并访问: `http://127.0.0.1:3000/swagger/docs`

## 数据库

- 数据库迁移
- 数据库假数据填充

```sh
just migrate

# 或者
go run main.go --migrate && go run main.go --seed
```

## swagger 文档

- 参考: https://docs.gofiber.cn/recipes/tableflip/

```sh
just docs

# 或者
swag init
```

## 异步队列

默认情况下, 没有启动这个服务, 需要手动修改 `main.go`

- [asynq](https://github.com/hibiken/asynq)
- [学习笔记](https://golang.liaohui5.cn/libs/9.%E5%BC%82%E6%AD%A5%E9%98%9F%E5%88%97-asynq.html)
- 注: 如果要使用异步队列, 需要先启动一个 redis 服务器

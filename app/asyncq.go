package app

import (
	"context"
	"fiber_curd_api/models"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func InitConnectRedisOpts() asynq.RedisClientOpt {
	host := Config.Get("redis.host").(string)
	port := Config.Get("redis.port").(int64)
	user := Config.Get("redis.user").(string)
	pass := Config.Get("redis.pass").(string)
	addr := fmt.Sprintf("%s:%d", host, port)

	fmt.Printf("[InitConnectRedisOpts] Addr = %s \n", addr)
	return asynq.RedisClientOpt{
		Addr:     addr,
		Username: user,
		Password: pass,
	}
}

// 全局 asynq client 实例
var asynqClient *asynq.Client

// InitAsynqClient 初始化 asynq 客户端
func InitAsynqClient() *asynq.Client {
	if asynqClient != nil {
		return asynqClient
	}
	asynqClient = asynq.NewClient(InitConnectRedisOpts())
	return asynqClient
}

// CloseAsynqClient 关闭 asynq 客户端
func CloseAsynqClient() {
	if asynqClient != nil {
		asynqClient.Close()
	}
}

// 全局 asynq server 实例
var asyncServer *asynq.Server

// 停止 asynq 服务端
func CloseAsynqServer() {
	if asyncServer != nil {
		asyncServer.Shutdown()
	}
}

// 启动 asynq 服务端
func StartAsynqServer() {
	asyncServer = asynq.NewServer(
		InitConnectRedisOpts(),
		asynq.Config{
			Concurrency: 4, // 并发数量
		},
	)

	server := asynq.NewServeMux()

	// 处理默认的异步任务队列
	server.HandleFunc("default", DefaultHandler)

	if err := asyncServer.Run(server); err != nil {
		log.Fatalf("[error]asynqServer 异常: %v", err)
	}
}

// DefaultHandler 默认队列处理器
func DefaultHandler(ctx context.Context, t *asynq.Task) error {
	taskId := string(t.Payload())

	// 查询数据库并输出
	var task models.Task
	result := ConnectDB().First(&task, taskId)
	if result.Error != nil {
		return fmt.Errorf("Canot find data by id=%s in database; %v", taskId, result.Error)
	}

	// 处理中: 更新数据库状态
	task.Status = models.TaskStatusProcessing
	ConnectDB().Save(&task)

	// 模拟处理过程
	time.Sleep(10 * time.Second)

	// 处理完成: 更新数据库状态
	task.Status = models.TaskStatusFinished
	ConnectDB().Save(&task)

	return nil
}

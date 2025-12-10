package main

import (
	"context"   // Go标准库：提供上下文(context)功能，用于控制goroutine的生命周期、取消操作和超时控制
	"flag"      // Go标准库：命令行参数解析包，用于解析程序启动时传入的命令行参数
	"fmt"       // Go标准库：格式化I/O包，提供格式化输入输出功能，如Printf、Sprintf等
	"log"       // Go标准库：简单日志包，提供基本的日志记录功能
	"net/http"  // Go标准库：HTTP客户端和服务器实现，提供HTTP协议相关功能
	"os"        // Go标准库：操作系统功能包，提供与操作系统交互的功能，如文件操作、环境变量等
	"os/signal" // Go标准库：信号处理包，用于处理系统信号，如Ctrl+C终止信号
	"syscall"   // Go标准库：系统调用包，包含系统相关的常量和类型，如信号类型
	"time"      // Go标准库：时间包，提供时间相关功能，如获取当前时间、时间格式化、定时器等

	// 内部包导入（项目内部模块）
	"github.com/MGter/xStreamTool_go/internal/api"    // API处理层：包含HTTP处理器和路由配置
	"github.com/MGter/xStreamTool_go/internal/config" // 配置管理：负责应用配置的加载和保存
	"github.com/MGter/xStreamTool_go/internal/store"  // 数据存储层：提供数据存储接口和内存存储实现
)

func main() {
	// 解析命令行参数
	port := flag.String("port", "8080", "服务器端口") // 定义port命令行参数，默认值"8080"，描述"服务器端口"
	debug := flag.Bool("debug", false, "启用调试模式") // 定义debug命令行参数，默认值false，描述"启用调试模式"
	flag.Parse()                                 // 解析命令行参数，将命令行参数值赋给对应的变量

	fmt.Println("🚀 xStreamTool Go HTTP 服务器启动中...") // 打印启动信息

	// 加载配置
	cfg := config.LoadConfig() // 调用配置模块的LoadConfig函数加载配置文件
	cfg.Server.Port = *port    // 用命令行参数覆盖配置中的端口设置（*是取指针值）
	cfg.Server.Debug = *debug  // 用命令行参数覆盖配置中的调试模式设置

	// 初始化存储
	todoStore := store.NewMemoryStore() // 创建内存存储实例，用于数据持久化

	// 初始化 API 处理器
	handler := api.NewHandler(todoStore) // 创建API处理器，传入存储实例作为依赖

	// 设置路由
	router := api.SetupRoutes(handler) // 设置所有HTTP路由，返回配置好的路由器

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port), // 服务器监听地址，格式为":端口号"
		Handler:      router,                              // 使用上面设置的路由器处理请求
		ReadTimeout:  15 * time.Second,                    // 读取请求超时时间
		WriteTimeout: 15 * time.Second,                    // 写入响应超时时间
		IdleTimeout:  60 * time.Second,                    // 空闲连接超时时间
	}

	// 优雅关闭 - 创建信号通道用于接收系统信号
	quit := make(chan os.Signal, 1)                      // 创建带缓冲区的信号通道，容量为1
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 注册信号监听，监听SIGINT(Ctrl+C)和SIGTERM(终止信号)

	// 启动服务器（在新的goroutine中）
	go func() {
		// 打印服务器信息
		log.Printf("📡 服务器监听地址: http://localhost:%s", cfg.Server.Port)          // 打印服务器访问地址
		log.Printf("📊 API 文档: http://localhost:%s/api/docs", cfg.Server.Port)  // 打印API文档地址
		log.Printf("🖥️  管理界面: http://localhost:%s/dashboard", cfg.Server.Port) // 打印管理界面地址
		log.Printf("🔧 调试模式: %v", cfg.Server.Debug)                             // 打印调试模式状态
		log.Println("🛑 按 Ctrl+C 停止服务器")                                        // 提示如何停止服务器

		// 启动HTTP服务器
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 如果启动失败且不是因为服务器已关闭，则记录致命错误
			log.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号（主goroutine阻塞在此处）
	<-quit                      // 阻塞等待直到收到信号，从quit通道接收到信号
	log.Println("🛑 正在关闭服务器...") // 打印正在关闭服务器的提示

	// 设置关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 创建30秒超时的上下文
	defer cancel()                                                           // 确保在函数返回时取消上下文，释放资源

	// 关闭 HTTP 服务器（优雅关闭）
	if err := server.Shutdown(ctx); err != nil {
		// 如果关闭失败，记录致命错误
		log.Fatalf("❌ 服务器关闭失败: %v", err)
	}

	log.Println("✅ 服务器已安全关闭") // 打印服务器已安全关闭的信息
}

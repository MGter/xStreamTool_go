package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger 自定义日志器
// 这是一个简单的日志工具，提供不同级别的日志记录功能
// 支持调试模式控制，调试信息只在debug模式下输出
type Logger struct {
	debug bool // 调试模式标志，为true时输出debug日志
}

// NewLogger 创建新的日志器
// 参数：debug - 是否启用调试模式
// 返回值：Logger实例指针
func NewLogger(debug bool) *Logger {
	return &Logger{debug: debug}
}

// Info 记录信息日志
// 用于记录一般的信息性消息
// 参数：format - 格式化字符串，v - 可变参数列表
func (l *Logger) Info(format string, v ...interface{}) {
	// 使用标准库log.Printf记录日志，前缀为[INFO]
	log.Printf("[INFO] "+format, v...)
}

// Error 记录错误日志
// 用于记录错误信息
// 参数：format - 格式化字符串，v - 可变参数列表
func (l *Logger) Error(format string, v ...interface{}) {
	// 使用标准库log.Printf记录日志，前缀为[ERROR]
	log.Printf("[ERROR] "+format, v...)
}

// Debug 记录调试日志
// 只在调试模式下输出，用于开发阶段的调试信息
// 参数：format - 格式化字符串，v - 可变参数列表
func (l *Logger) Debug(format string, v ...interface{}) {
	// 检查是否启用调试模式
	if l.debug {
		// 只在debug为true时记录调试日志，前缀为[DEBUG]
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Warn 记录警告日志
// 用于记录警告信息
// 参数：format - 格式化字符串，v - 可变参数列表
func (l *Logger) Warn(format string, v ...interface{}) {
	// 使用标准库log.Printf记录日志，前缀为[WARN]
	log.Printf("[WARN] "+format, v...)
}

// Fatal 记录致命错误并退出
// 记录致命错误后程序会立即退出，返回状态码1
// 参数：format - 格式化字符串，v - 可变参数列表
func (l *Logger) Fatal(format string, v ...interface{}) {
	// 使用标准库log.Fatalf记录日志并退出程序，前缀为[FATAL]
	log.Fatalf("[FATAL] "+format, v...)
}

// FileLogger 文件日志器
// 将日志写入文件的日志器
type FileLogger struct {
	file *os.File // 日志文件句柄
}

// NewFileLogger 创建文件日志器
// 参数：filename - 日志文件路径
// 返回值：FileLogger实例指针和可能的错误
func NewFileLogger(filename string) (*FileLogger, error) {
	// 打开或创建日志文件
	// os.O_APPEND - 以追加模式打开文件
	// os.O_CREATE - 如果文件不存在则创建
	// os.O_WRONLY - 只写模式
	// 0644 - 文件权限：所有者可读写，其他人只读
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err // 如果打开文件失败，返回错误
	}

	// 创建并返回FileLogger实例
	return &FileLogger{file: file}, nil
}

// Log 写入日志
// 将日志条目写入文件，包含时间戳和日志级别
// 参数：level - 日志级别，message - 日志消息
func (l *FileLogger) Log(level, message string) {
	// 生成当前时间戳，格式为：2006-01-02 15:04:05
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 格式化日志条目：[时间戳] 级别: 消息
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)

	// 将日志条目写入文件
	l.file.WriteString(logEntry)
}

// Close 关闭日志文件
// 关闭文件句柄，释放资源
func (l *FileLogger) Close() {
	l.file.Close()
}

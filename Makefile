.PHONY: run build clean test

run:
	@echo "🚀 启动 xStreamTool Go..."
	@go run ./cmd/xstream/main.go

build:
	@echo "📦 构建项目..."
	@go build -o bin/xstream ./cmd/xstream/main.go

clean:
	@echo "🧹 清理文件..."
	@rm -rf bin/
	@go clean

test:
	@echo "🧪 运行测试..."
	@go test ./...

dev:
	@go run ./cmd/xstream/main.go --debug=true

deps:
	@go mod tidy
	@go mod download

help:
	@echo "可用命令:"
	@echo "  make run    - 运行服务器"
	@echo "  make dev    - 开发模式运行"
	@echo "  make build  - 构建项目"
	@echo "  make clean  - 清理文件"
	@echo "  make test   - 运行测试"
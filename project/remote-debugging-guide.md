# 远程调试镜像构建指南

本文档介绍如何为 New API 项目构建支持远程调试的 Docker 镜像。

## 目录

- [方案概述](#方案概述)
- [方案一：Delve 调试器（推荐）](#方案一delve-调试器推荐)
- [方案二：Docker Compose 调试环境](#方案二docker-compose-调试环境)
- [IDE 远程调试配置](#ide-远程调试配置)
- [常见问题](#常见问题)

## 方案概述

远程调试适用于以下场景：
- 在 Docker 容器中运行的应用出现难以复现的问题
- 需要在类生产环境中调试
- 团队协作调试

**核心工具**: [Delve](https://github.com/go-delve/delve) - Go 语言官方调试器

## 方案一：Delve 调试器（推荐）

### 1. 创建调试版 Dockerfile

创建 `Dockerfile.debug` 文件：

```dockerfile
# ===========================================
# 构建阶段
# ===========================================
FROM golang:1.25.1-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git make

# 安装 Delve 调试器
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译带调试信息的二进制文件
# -gcflags="all=-N -l" 参数说明：
#   -N: 禁用优化
#   -l: 禁用内联
# 这样可以保证调试体验，但会降低性能
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -gcflags="all=-N -l" \
    -ldflags="-X 'main.Version=debug'" \
    -o /app/new-api \
    main.go

# ===========================================
# 运行阶段
# ===========================================
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 复制二进制文件和调试器
COPY --from=builder /app/new-api /app/new-api
COPY --from=builder /go/bin/dlv /app/dlv

# 复制配置文件示例
COPY .env.example /app/.env.example

# 创建必要的目录
RUN mkdir -p /app/data /app/logs

# 暴露端口
# 3000: 应用端口
# 2345: Delve 调试端口
EXPOSE 3000 2345

# 使用 Delve headless 模式启动应用
# --listen=:2345: 监听调试连接
# --headless=true: 无终端模式
# --api-version=2: 使用 DLV API v2
# --accept-multiclient: 允许多个调试客户端连接
CMD ["/app/dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/new-api", "--", "-port=3000"]
```

### 2. 构建调试镜像

```bash
# 构建调试镜像
docker build -f Dockerfile.debug -t new-api:debug .

# 查看镜像大小
docker images | grep new-api
```

### 3. 运行调试容器

```bash
# 方式一：直接运行
docker run -d \
  --name new-api-debug \
  -p 3000:3000 \
  -p 2345:2345 \
  -v $(pwd)/data:/app/data \
  -e SQL_DSN="postgresql://user:password@host:5432/new-api" \
  -e REDIS_CONN_STRING="redis://host:6379" \
  -e SESSION_SECRET="your-session-secret" \
  new-api:debug

# 方式二：使用 docker-compose（推荐，见下一节）
```

### 4. 验证调试端口

```bash
# 检查容器是否正常运行
docker ps | grep new-api-debug

# 检查调试端口是否监听
nc -zv localhost 2345

# 查看容器日志
docker logs -f new-api-debug
```

## 方案二：Docker Compose 调试环境

### 1. 创建 docker-compose.debug.yml

```yaml
version: '3.8'

services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:16-alpine
    container_name: new-api-postgres-debug
    environment:
      POSTGRES_DB: new-api
      POSTGRES_USER: new-api
      POSTGRES_PASSWORD: new-api-password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - new-api-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U new-api"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: new-api-redis-debug
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    networks:
      - new-api-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  # New API 调试版
  new-api-debug:
    build:
      context: .
      dockerfile: Dockerfile.debug
    container_name: new-api-debug
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      # 数据库配置
      SQL_DSN: "postgresql://new-api:new-api-password@postgres:5432/new-api?sslmode=disable"

      # Redis 配置
      REDIS_CONN_STRING: "redis://redis:6379"

      # 必需的安全配置
      SESSION_SECRET: "debug-session-secret-change-in-production"
      CRYPTO_SECRET: "debug-crypto-secret-change-in-production"

      # 其他配置
      TZ: Asia/Shanghai
      LOG_LEVEL: debug

      # 调试相关
      STREAMING_TIMEOUT: 600
    ports:
      - "3000:3000"   # 应用端口
      - "2345:2345"   # Delve 调试端口
    volumes:
      # 挂载源代码，实现热重载（可选）
      - ./:/app/src
      - ./data:/app/data
      - ./logs:/app/logs
    networks:
      - new-api-network
    restart: unless-stopped

volumes:
  postgres-data:
  redis-data:

networks:
  new-api-network:
    driver: bridge
```

### 2. 启动调试环境

```bash
# 构建并启动所有服务
docker-compose -f docker-compose.debug.yml up -d

# 查看服务状态
docker-compose -f docker-compose.debug.yml ps

# 查看日志
docker-compose -f docker-compose.debug.yml logs -f new-api-debug

# 停止服务
docker-compose -f docker-compose.debug.yml down

# 停止服务并删除数据卷
docker-compose -f docker-compose.debug.yml down -v
```

### 3. 重新构建镜像

```bash
# 强制重新构建
docker-compose -f docker-compose.debug.yml build --no-cache

# 构建特定服务
docker-compose -f docker-compose.debug.yml build new-api-debug
```

## IDE 远程调试配置

### VS Code 配置

#### 1. 安装必要扩展

```bash
# 安装 Go 扩展
code --install-extension golang.go

# 安装 Docker 扩展（可选）
code --install-extension ms-azuretools.vscode-docker
```

#### 2. 创建调试配置

在项目根目录创建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Docker Remote Debug",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app/new-api",
      "port": 2345,
      "host": "localhost",
      "showLog": true,
      "trace": "log",
      "apiVersion": 2,
      "preLaunchTask": "docker-debug-build"
    },
    {
      "name": "Docker Remote Debug (Hot Reload)",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app/src",
      "port": 2345,
      "host": "localhost",
      "showLog": true,
      "trace": "verbose"
    }
  ]
}
```

#### 3. 创建构建任务

在项目根目录创建 `.vscode/tasks.json`：

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "docker-debug-build",
      "type": "shell",
      "command": "docker-compose -f docker-compose.debug.yml build new-api-debug",
      "group": "build",
      "problemMatcher": []
    },
    {
      "label": "docker-debug-restart",
      "type": "shell",
      "command": "docker-compose -f docker-compose.debug.yml restart new-api-debug",
      "group": "build",
      "problemMatcher": []
    },
    {
      "label": "docker-debug-logs",
      "type": "shell",
      "command": "docker-compose -f docker-compose.debug.yml logs -f new-api-debug",
      "group": "build",
      "problemMatcher": [],
      "isBackground": true
    }
  ]
}
```

#### 4. 开始调试

1. 启动调试容器：
```bash
docker-compose -f docker-compose.debug.yml up -d
```

2. 在 VS Code 中按 `F5` 或点击 "Run and Debug"

3. 设置断点并触发 API 请求

### GoLand / IntelliJ IDEA 配置

#### 1. 创建远程调试配置

1. 打开 `Run → Edit Configurations...`
2. 点击 `+` → 选择 `Go Remote`
3. 配置如下：
   - **Name**: `Docker Remote Debug`
   - **Host**: `localhost`
   - **Port**: `2345`
4. 保存配置

#### 2. 开始调试

1. 启动调试容器：
```bash
docker-compose -f docker-compose.debug.yml up -d
```

2. 在 IDEA 中选择 `Docker Remote Debug` 配置

3. 点击 Debug 按钮（虫子图标）

4. 设置断点并触发请求

### Vim / Neovim 配置

使用 [vim-delve](https://github.com/brentlt/vim-delve) 或 [nvim-dap](https://github.com/mfussenegger/nvim-dap)：

```lua
-- init.lua
local dap = require('dap')

dap.adapters.go = {
  type = 'server',
  host = '127.0.0.1',
  port = 2345,
}

dap.configurations.go = {
  {
    type = 'go',
    name = 'Docker Remote Debug',
    request = 'attach',
    mode = 'remote',
    remotePath = '/app/new-api',
    port = 2345,
    host = '127.0.0.1',
  }
}
```

## 常见问题

### Q1: 连接调试端口失败

**问题**: `could not connect to remote debugger`

**解决方案**:
```bash
# 检查端口是否监听
netstat -tuln | grep 2345

# 检查容器是否运行
docker ps | grep new-api-debug

# 检查端口映射
docker port new-api-debug

# 重启容器
docker-compose -f docker-compose.debug.yml restart new-api-debug
```

### Q2: 断点不起作用

**问题**: 设置断点但程序没有停止

**解决方案**:
1. 确保使用 `-gcflags="all=-N -l"` 编译
2. 检查 `remotePath` 配置是否正确
3. 确认代码版本与容器中运行的代码一致

### Q3: 性能太慢

**问题**: 调试版本性能下降明显

**解决方案**:
- 这是正常现象，调试版本禁用了优化
- 生产环境使用优化版本（去掉 `-gcflags`）
- 或者使用 `perf` 工具进行性能分析

### Q4: 如何调试特定的 Goroutine

在 Delve 中：
```go
// 列出所有 goroutine
(goroutines)

// 切换到指定 goroutine
(goroutine 123)

// 查看调用栈
(stack)

// 查看局部变量
(locals)
```

### Q5: 多容器调试

如果需要同时调试多个服务：

```yaml
# docker-compose.debug.yml
services:
  new-api-debug-1:
    ports:
      - "3000:3000"
      - "2345:2345"

  new-api-debug-2:
    ports:
      - "3001:3000"
      - "2346:2345"
```

然后在 IDE 中配置多个调试配置，使用不同的端口。

### Q6: 生产环境可以使用吗？

**强烈不建议**！原因：

1. 调试端口暴露安全风险
2. 性能下降严重（可能 10x+）
3. 可能泄露敏感信息

生产环境应该：
- 使用优化构建（去掉调试标志）
- 移除 Delve 调试器
- 限制端口暴露

## 调试技巧

### 1. 条件断点

在 VS Code 中：
```go
// 只在特定条件下触发断点
// 在行号右侧点击设置条件
if userId == 12345 {
    // 这里设置断点
}
```

### 2. 日志断点

不停止程序，只输出日志：
```
// 在断点设置中添加日志消息
User {userId} requested endpoint {endpoint}
```

### 3. 表达式求值

调试时在 Debug Console 中执行：
```go
// 查看变量值
print user.Token

// 调用函数
call utils.FormatTime(time.Now())

// 修改变量值
set retryCount = 0
```

### 4. 性能分析

```bash
# 进入容器
docker exec -it new-api-debug sh

# 使用 Go pprof 工具
wget http://localhost:3000/debug/pprof/heap
go tool pprof heap

# 或者使用 Delve 的 profiling 功能
dlv attach 1
(dlv) profile cpu.zip
# ... 等待一段时间 ...
(dlv) profile -output=cpu.zip
```

## 参考资源

- [Delve 官方文档](https://github.com/go-delve/delve/tree/master/Documentation)
- [Go 调试指南](https://go.dev/doc/gdb)
- [VS Code Go 调试](https://github.com/golang/vscode-go/wiki/debugging)
- [Docker 最佳实践](https://docs.docker.com/develop/dev-best-practices/)

## 许可证

本文档遵循项目的 AGPLv3 许可证。

# go-note-api

## 项目说明

go-note-api 是一个轻量级的笔记管理 Web API 服务，适用于个人知识管理或团队协作场景，
支持快速记录和检索信息，数据持久化存储不丢失。

服务基于 Go 标准库实现，使用文件作为持久化存储（默认 `notes.json`），无需依赖任何外部
服务（如 MySQL、Redis、MQ 等）。启动后即可独立运行。

### 功能

- 创建笔记（POST `/api/notes`）
- 获取单个笔记（GET `/api/notes/{id}`）
- 列出全部笔记（GET `/api/notes`）
- 更新笔记（PUT `/api/notes/{id}`）
- 删除笔记（DELETE `/api/notes/{id}`）
- 搜索笔记，按标题或内容匹配（GET `/api/notes/search?q=关键词`）
- 获取最近创建的笔记（GET `/api/notes/recent?limit=N`）

### 请求/响应示例

创建笔记：

```bash
curl -X POST http://localhost:8080/api/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"第一条笔记","content":"Hello go-note-api","tags":["demo"]}'
```

返回：

```json
{"id":"1","title":"第一条笔记","content":"Hello go-note-api","tags":["demo"],"created_at":"...","updated_at":"..."}
```

## 项目结构

```
go-note-api/
├── cmd/main.go              # 程序入口
├── internal/
│   ├── model/note.go        # Note 领域模型与校验
│   ├── store/store.go       # 文件持久化存储层
│   ├── service/service.go   # 业务逻辑层
│   ├── handler/handler.go   # HTTP 处理层
│   └── server/server.go     # HTTP 服务器生命周期
├── go.mod
├── benzhi.Dockerfile
└── benzhi.build_docker.sh
```

## 标准命令

```bash
go build ./...     # 编译
go test ./...      # 测试
go run ./cmd       # 启动（默认监听 :8080，持久化文件 notes.json）
```

启动时可通过参数自定义监听地址与持久化文件：

```bash
go run ./cmd -addr=:9090 -db=/path/to/notes.json
```

## Docker 命令

```bash
docker build -t go-note-api .                                  # 构建镜像
docker run -it go-note-api:latest                              # 运行容器
docker build --platform linux/arm64 -t go-note-api .           # 构建 arm64
```

也可使用构建脚本：

```bash
./benzhi.build_docker.sh go-note-api linux/amd64
./benzhi.build_docker.sh go-note-api linux/arm64
```

容器内可直接使用 `go build ./...`、`go test ./...`、`go vet ./...` 等命令。

# tasks.md · myblog

> 可勾选的执行 checklist(CLAUDE.md 强制)。每完成一步立刻勾选。

## M0 · 视觉预览(已完成)
- [x] 建立 `docs/preview/` 目录
- [x] 第一版莫兰迪风(被否)
- [x] 第二版紫色科技感(已确认)
- [x] 五张预览页:index / tech / photography / about / blog-detail

## M1 · 后端骨架(进行中)

### 工程
- [x] 目录结构:cmd/server, internal/{config,bootstrap,handler,middleware,service,repository,model,dto,router,pkg}, migrations, deploy, docs, test
- [x] `go.mod` 初始化(module `github.com/yinyin/myblog`)
- [x] `Makefile`:build / run / test / tidy / lint / fmt
- [x] `.gitignore`、`.env.example`、`configs/config.example.yaml`

### 基础 pkg
- [x] `internal/pkg/logger`:zap JSON,带 trace_id 字段
- [x] `internal/pkg/response`:统一 `{code,message,data}` 封装 + helper
- [x] `internal/pkg/errcode`:业务错误码表(预留文件结构,M1 只放 OK 和若干通用码)

### config & bootstrap
- [x] `internal/config/config.go`:viper 加载 yaml + env 覆盖,结构体定义 Server / Log
- [x] `internal/bootstrap/app.go`:装配 logger、Gin engine、路由、优雅关停

### middleware
- [x] `internal/middleware/trace.go`:生成 / 透传 trace_id
- [x] `internal/middleware/access_log.go`:结构化访问日志
- [x] `internal/middleware/recover.go`:recover 后返回统一错误响应
- [x] `internal/middleware/cors.go`:开发环境宽松 CORS

### handler / router
- [x] `internal/handler/public/health.go`:GET /health,返回统一响应
- [x] `internal/router/router.go`:`/api/v1/public/health`

### 入口
- [x] `cmd/server/main.go`:调用 bootstrap.Run(),处理信号优雅关停

### 测试(TDD)
- [x] `test/health_test.go`:先写测试,预期 200 + `code=0` + `data.status=up`
- [x] `go test ./...` 通过

### 验证
- [x] `make run` 启动成功(用 MYBLOG_SERVER_PORT=18080 避开本地占用)
- [x] `curl /api/v1/public/health` 返回期望格式
- [x] 访问日志包含 `trace_id`
- [x] 主动触发 panic 能被 recover 捕获(由 TestRecover_PanicReturnsUnified500 覆盖)

### 收尾
- [ ] 更新 `context.md` 记录 M1 完成状态
- [ ] git 首次提交(待用户确认 `git init`)

## M2 · 数据层 & 登录(已完成)

### 环境
- [x] 启动本地 Redis(`redis-server --daemonize yes`,6379 端口)
- [x] 创建 myblog 数据库 + myblog 业务账号(localhost/127.0.0.1 均授权)

### 配置
- [x] `config` 新增 MySQL / Redis / JWT
- [x] `configs/config.yaml`、`config.example.yaml`、`.env.example` 同步

### 数据模型
- [x] `migrations/000001_init_schema.{up,down}.sql`:users/categories/tags/articles/article_tags/profiles
- [x] `internal/model/{base,user}.go`:BaseModel + User

### pkg
- [x] `internal/pkg/jwt`:HS256 access+refresh 签发/校验(单测覆盖:issue/verify/过期/类型错/签名错)
- [x] `internal/pkg/password`:bcrypt cost=10 哈希/校验(单测覆盖:正确/错误/空密码)
- [x] `internal/pkg/errcode`:新增 3xxxx 用户错误码

### bootstrap / 基础设施
- [x] `bootstrap/db.go`:GORM + MySQL 连接池 + GORM 日志桥接到 zap
- [x] `bootstrap/redis.go`:go-redis v9 + 健康检查
- [x] `bootstrap/app.go`:装配 db/redis/signer/repo/service 注入 router

### 业务
- [x] `repository.UserRepo`:Create / GetByID / GetByUsername
- [x] `service.AuthService.Login`:密码校验 + token 签发(不泄漏"用户存不存在")
- [x] `handler/admin/auth.go`:POST /api/v1/admin/auth/login,映射 service 错误到 errcode
- [x] `handler/admin/ping.go`:GET /api/v1/admin/ping(受保护)
- [x] `middleware/auth.go`:JWTAuth 中间件 + `UserIDFrom` helper
- [x] `router/router.go`:挂载 /admin/auth/login(公开)+ /admin/ping(JWTAuth 保护)

### 运维脚手架
- [x] `cmd/seed/main.go`:交互式创建/重置管理员
- [x] `Makefile`:`migrate-up` / `migrate-down` / `seed-admin`

### 验证
- [x] `go vet ./...` 通过
- [x] `go test -race ./...` 全绿(含 jwt 4 个 case + password 2 个 case + health 4 个 case)
- [x] `make migrate-up` 成功建 6 张表
- [x] `make seed-admin PW='Admin@123'` 创建 admin 账号
- [x] 端到端:health / 错误密码 / 正确密码 / 无 token / 坏 token / 有 token / 参数校验 全部符合预期

### 收尾
- [x] 更新 `context.md` / `tasks.md`

## M3 ~ M5
> 详情见 `~/.claude/plans/10-curious-flurry.md` 的里程碑定义。

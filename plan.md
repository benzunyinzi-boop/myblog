# plan.md · myblog

> 任务级 Plan 文档(CLAUDE.md 强制)。全局设计方案见 `~/.claude/plans/10-curious-flurry.md`。

## 当前任务

**M1 · 后端骨架(Week 1)** — 起一个能跑的 Go/Gin 最小服务,打通配置、日志、中间件、统一响应和 health 探活。

## 做什么

1. 目录骨架按 CLAUDE.md 严格划分:`cmd/server`, `internal/{config,bootstrap,handler,middleware,service,repository,model,dto,router,pkg}`, `migrations`, `deploy`, `docs`, `test`
2. 基础 pkg:`logger`(zap JSON)、`response`(统一 `{code,message,data}`)、`errcode`(业务错误码)
3. `config` 用 viper 加载 `configs/config.yaml`,支持环境变量覆盖
4. `bootstrap` 装配 logger / Gin engine / 路由 / 优雅关停
5. 四个中间件:`trace_id` / `access_log` / `recover` / `cors`
6. 一个 `GET /api/v1/public/health` handler 作为探活,返回统一响应
7. `cmd/server/main.go` 入口拼装 bootstrap
8. Makefile(build / run / test / lint / tidy)
9. 集成测试(TDD):先写 `test/health_test.go`,用 httptest 起完整路由跑通

## 怎么做(顺序)

```
三件套 → 目录/go.mod/Makefile → pkg/{logger,response,errcode}
     → config(viper) → bootstrap → middleware → health handler
     → router → main.go → 先写 test(TDD) → go test → make run 验证
```

## 不做什么(本迭代范围外)

- 不接 MySQL / Redis(M1 先不连,避免本地环境阻塞。M2 再接)
- 不实现 JWT / 登录(放 M1 尾部或 M2 开头)
- 不实现业务 handler / service / repository
- 不写数据库迁移文件(留给接入 DB 时一起做)

## 验收标准

- [x] 目录结构符合 CLAUDE.md
- [x] `go test ./...` 通过,包含 `/health` 集成测试
- [x] `make run` 启动服务,`curl http://localhost:8080/api/v1/public/health` 返回
      ```json
      {"code":0,"message":"ok","data":{"status":"up","service":"myblog","ts":...}}
      ```
- [x] 每条请求日志(JSON)包含 `trace_id`/`method`/`path`/`status`/`latency`
- [x] 触发 panic 能被 recover 中间件捕获并返回统一错误响应
- [x] `make lint` 通过(或 golangci-lint 缺失时在文档注明)

## 风险与对策

| 风险 | 对策 |
|---|---|
| 本地没装 golangci-lint | Makefile 里做 check,缺失时打印提示而非 fail |
| 8080 端口被占 | 配置可覆盖;`make run` 前检查端口 |
| viper 读不到配置文件 | 允许配置缺失,使用内置默认值,仅在严重字段缺失时 fail-fast |

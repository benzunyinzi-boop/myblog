# tasks.md · myblog

> 可勾选的执行 checklist(CLAUDE.md 强制)。每完成一步立刻勾选。

## M0 · 视觉预览(已完成)
- [x] 建立 `docs/preview/` 目录
- [x] 第一版莫兰迪风(被否)
- [x] 第二版紫色科技感(已确认)
- [x] 五张预览页:index / tech / photography / about / blog-detail

## M1 · 后端骨架(已完成)

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
- [x] 更新 `context.md` 记录 M1 完成状态
- [x] git 首次提交(afeb196 chore: init project skeleton (M1 + M2))

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

## M3 · 核心 API(已完成)

### 分类管理
- [x] `model/category.go`:Category 模型
- [x] `repository/category_repo.go`:CRUD + FirstOrCreateByName
- [x] `service/category_service.go`:业务逻辑 + 错误映射
- [x] `handler/admin/category.go`:POST/PUT/DELETE/GET(需 JWT)
- [x] `handler/public/category.go`:GET(公开)
- [x] 路由挂载 + bootstrap 装配

### 标签管理
- [x] `model/tag.go`:Tag 模型
- [x] `repository/tag_repo.go`:CRUD + FirstOrCreateByName + ListByIDs
- [x] `service/tag_service.go`:业务逻辑(无 Update)
- [x] `handler/admin/tag.go`:POST/DELETE/GET(需 JWT)
- [x] `handler/public/tag.go`:GET(公开)
- [x] 路由挂载 + bootstrap 装配

### 文章管理
- [x] `model/article.go`:Article 模型 + 状态机(draft/published)
- [x] `model/article_tag.go`:ArticleTag 关联表
- [x] `repository/article_repo.go`:CRUD + 事务 + 过滤分页 + ReplaceTags
- [x] `service/article_service.go`:CRUD + Publish/Unpublish + IncrView
- [x] `dto/article.go`:Create/Update/List 请求 + Summary/Detail/List 响应
- [x] `handler/admin/article.go`:完整 CRUD + publish/unpublish(需 JWT)
- [x] `handler/public/article.go`:列表(强制 published)+ 详情(slug + 浏览计数)
- [x] 路由挂载 + bootstrap 装配

### 个人资料
- [x] `model/profile.go`:Profile 模型(单条记录,id=1)
- [x] `repository/profile_repo.go`:Get + Upsert(ON DUPLICATE KEY UPDATE)
- [x] `service/profile_service.go`:Get(不存在返回空)+ Update
- [x] `dto/profile.go`:ProfileResp + ProfileUpdateReq
- [x] `handler/admin/profile.go`:PUT(需 JWT)
- [x] `handler/public/profile.go`:GET(公开)
- [x] 路由挂载 + bootstrap 装配

### 文件上传
- [x] `pkg/uploader/uploader.go`:Uploader 接口 + LocalUploader 实现
- [x] `handler/admin/upload.go`:POST /uploads(需 JWT,限 10MB)
- [x] 静态文件服务:`r.Static("/uploads", "./uploads")`
- [x] 路由挂载 + bootstrap 装配

### 辅助 pkg
- [x] `pkg/validation/validation.go`:自定义 slug validator(小写字母数字短横线)
- [x] `pkg/dberr/dberr.go`:MySQL 1062 唯一键冲突判断
- [x] `service/errcode.go`:统一管理 service 层错误

### migration
- [x] `migrations/000002_profile_schema_rewrite.{up,down}.sql`:修正 Profile schema

### 文档
- [x] `docs/api.md`:完整 API 文档(50+ 接口)
- [x] `docs/M3-article-verification.md`:Article CRUD 验证报告
- [x] `docs/M3-verification-plan.md`:完整验证计划(9 大类 50+ 验证点)
- [x] `docs/M3-e2e-verification-result.md`:端到端验证执行报告(25/25 全绿)

### 测试
- [x] `test/e2e-verify.sh`:端到端验证脚本(7 模块 25 断言)
- [x] 修复 Bug 1:浏览计数 ctx 泄漏(25eb86b)
- [x] 修复 Bug 2:Profile schema 不匹配(9d6ef1d)
- [x] 修复 Bug 3:GORM snake_case 命名(55bf31e)
- [x] 修复 Bug 4:Upsert 零值时间戳(55bf31e)
- [x] 最终验证:25 通过 / 0 失败 ✓

### 收尾
- [x] 更新 `context.md` 记录 M3 完成状态 + 关键决策
- [x] git 推送所有 commit 到 GitHub(benzunyinzi-boop/myblog)

## M4 · 前端开发(待开始)

### 脚手架
- [ ] Vite + Vue 3 + TypeScript 初始化
- [ ] 安装 Naive UI + 配置主题
- [ ] 目录结构:`web/src/{views,components,api,stores,router,styles,utils}`
- [ ] 迁移 `docs/preview/styles.css` 到 `web/src/styles/theme.css`(紫色科技感)

### 路由 & 布局
- [ ] Vue Router 配置:/ / /tech / /photography / /about / /blog/:slug
- [ ] 全局布局:导航栏(四项)+ 主题切换 + 页脚
- [ ] 响应式适配(移动端 / 桌面端)

### 公开页面
- [ ] 首页:最新文章列表 + 分类导航
- [ ] 技术栏:按分类过滤文章(Golang/MySQL/Redis/MQ/运维)
- [ ] 摄影栏:图片瀑布流(预留,M4 可用占位)
- [ ] 关于我:Profile 展示 + 社交链接
- [ ] 文章详情:Markdown 渲染 + 目录 + 浏览计数 + 标签

### API 对接
- [ ] `api/auth.ts`:login
- [ ] `api/article.ts`:列表 / 详情(slug)
- [ ] `api/category.ts`:列表
- [ ] `api/tag.ts`:列表
- [ ] `api/profile.ts`:GET
- [ ] axios 封装:统一错误处理 + token 注入

### 管理后台
- [ ] 登录页:`/admin/login`
- [ ] 后台布局:侧边栏(文章/分类/标签/关于/上传)+ 顶栏(用户信息/退出)
- [ ] 文章管理:列表(分页+过滤)+ 新建 + 编辑 + 删除 + 发布/下线
- [ ] Markdown 编辑器:集成 vditor 或 bytemd
- [ ] 分类管理:列表 + 新建 + 编辑 + 删除
- [ ] 标签管理:列表 + 新建 + 删除
- [ ] 关于我:表单编辑 Profile
- [ ] 图片上传:拖拽上传 + 预览 + 插入 Markdown

### 状态管理
- [ ] Pinia stores:user(登录状态)/ article / category / tag
- [ ] localStorage 持久化 token

### 验证
- [ ] 本地开发:`npm run dev` 能访问所有页面
- [ ] 前后端联调:登录 → 文章列表 → 详情 → 管理后台 CRUD
- [ ] 响应式测试:移动端 / 平板 / 桌面
- [ ] 构建:`npm run build` 生成 dist/

### 收尾
- [ ] 更新 `context.md` / `tasks.md`
- [ ] git commit + push

## M5 · 部署(待开始)

### Docker 化
- [ ] `Dockerfile`:多阶段构建(builder + runtime)
- [ ] `docker-compose.yml`:mysql / redis / myblog 三容器
- [ ] `.dockerignore`:排除 node_modules / .git / tmp
- [ ] 本地验证:`docker-compose up` 能完整启动

### Nginx
- [ ] `deploy/nginx.conf`:反向代理 /api → :8080,/ → 静态文件
- [ ] 静态文件服务:`/uploads` 映射到容器内 `/app/uploads`
- [ ] gzip 压缩 + 缓存策略

### 云服务器
- [ ] 购买 VPS(阿里云 / 腾讯云,1C2G 起步)
- [ ] 安装 Docker + docker-compose
- [ ] 配置防火墙:开放 80 / 443 / 22
- [ ] 上传 docker-compose.yml + 配置文件

### 域名 & HTTPS
- [ ] 域名解析:A 记录指向 VPS IP
- [ ] Let's Encrypt 证书:certbot 自动续期
- [ ] Nginx 配置 HTTPS + HTTP → HTTPS 重定向

### CI/CD(可选)
- [ ] GitHub Actions:push main → 自动构建镜像 → 推送到 Docker Hub
- [ ] 服务器 webhook:拉取新镜像 → docker-compose up -d

### 监控 & 日志
- [ ] 日志收集:容器日志 → 宿主机 /var/log/myblog
- [ ] 健康检查:定时 curl /api/v1/public/health
- [ ] 告警(可选):Prometheus + Grafana

### 验证
- [ ] 通过域名访问博客首页
- [ ] HTTPS 证书有效
- [ ] 管理后台能登录并 CRUD
- [ ] 图片上传 + 访问正常
- [ ] 服务重启后数据不丢失

### 收尾
- [ ] 更新 `context.md` / `tasks.md`
- [ ] 编写 `README.md`:项目介绍 + 部署指南
- [ ] git commit + push

## 下一步

**当前状态**:M3 已完成,25/25 端到端验证全绿,所有 commit 已推送到 GitHub。

**明天继续**:
1. 如果说"继续",默认进入 **M4 · 前端开发**
2. 如果想先部署,说"跳到 M5"
3. 如果想补充测试 / 优化性能,说明具体需求

**关键文档**:
- `docs/api.md`:完整 API 文档
- `docs/M3-e2e-verification-result.md`:验证报告
- `context.md`:决策记录 + 环境约定
- `test/e2e-verify.sh`:一键验证脚本

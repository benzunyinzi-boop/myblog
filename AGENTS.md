# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## 项目概述

- 项目名称：myblog
- 业务描述：博客微服务
- 服务类型：Go 微服务
- 主要协议：HTTP (Gin) + gRPC

## 技术栈

- 语言：Go 1.22+
- Web 框架：Gin
- ORM：GORM
- 数据库：MySQL (主库)
- 缓存：Redis
- 日志：zap


你是我的 Go 后端协作工程师，请严格遵守以下规范：

1) 目录结构必须使用：
cmd/server, internal/{config,bootstrap,handler,middleware,service,repository,model,dto,router}, migrations, deploy, docs, test。

2) 分层职责：
- handler: 参数解析/校验/返回
- service: 业务逻辑
- repository: DB/Redis访问
不得跨层直接调用（handler不能直接查库）。

3) 代码风格：
- 必须可通过 gofmt/goimports/golangci-lint
- context放第一参数
- 错误必须wrap: fmt.Errorf("...: %w", err)
- 禁止 panic 处理业务错误

4) 命名规范：
- 导出 PascalCase，非导出 camelCase
- 缩写用 ID/URL/HTTP
- DTO 命名：XxxReq/XxxResp
- 方法名动词开头：Create/List/Get/Update/Delete/Approve/Reject

5) 输出要求：
每次改动必须给出文件清单、关键代码片段、运行命令、测试命令。


## 常用命令

```bash
make build          # 编译
make test           # 运行测试
make lint           # 代码检查 (golangci-lint)
make proto          # 生成 protobuf 代码
make migrate-up     # 执行数据库迁移
make migrate-down   # 回滚数据库迁移
make run            # 本地运行
go test ./...       # 运行全部测试
go test ./internal/service/ -run TestGetUser  # 运行单个测试
```

## 项目结构

```
.
├── cmd/server/main.go       # 服务入口
├── internal/
│   ├── handler/             # HTTP handler (Gin)，参数校验、请求/响应转换
│   ├── service/             # 核心业务逻辑，事务编排
│   ├── repository/          # 数据访问层 (GORM)
│   ├── model/               # 数据模型 / 领域对象
│   ├── middleware/           # 中间件
│   └── pkg/                 # 项目内部公共包（含 errcode/ 业务错误码）
├── api/                     # protobuf 定义 / OpenAPI spec
├── pkg/                     # 可被外部引用的公共包
├── configs/                 # 配置文件
├── migrations/              # 数据库迁移文件
└── test/                    # 集成测试 / E2E 测试
```

## 分层架构

handler → service → repository 三层，依赖方向只能向下，禁止跨层调用和循环依赖。

- handler：参数校验、请求/响应转换，不含业务逻辑
- service：核心业务逻辑，事务编排，调用 repository 和外部服务
- repository：纯数据访问，GORM 操作封装，不含业务判断

## 编码规范

### 命名
- 文件名：snake_case（user_handler.go）
- 包名：小写单词，不用下划线（userservice 而非 user_service）
- 接口：动词或 -er 后缀（Reader, UserService）

### 错误处理
- 始终检查 error，不用 _ 忽略
- 使用 fmt.Errorf + %w 包装错误，保留错误链
- 业务错误码统一定义在 internal/pkg/errcode/
- HTTP 响应统一格式：`{"code": 0, "message": "ok", "data": {}}`
- gRPC 使用 status.Error 返回标准错误码
- 不要 panic，除非不可恢复的程序错误

### 并发
- goroutine 必须有明确退出机制（context cancel / done channel）
- 使用 errgroup 管理并发任务组

## 日志与可观测性规则
1. 后端日志统一使用 Zap JSON 格式。
2. 所有请求日志必须包含：`trace_id`、path、method、status、latency。
3. 错误日志必须包含：错误原因、关键参数、`trace_id`。
4. 数据库慢查询需记录阈值与 SQL 摘要（避免泄露敏感信息）。
5. 不要在循环中打日志


## 数据库规范

- Model 嵌入 gorm.Model 或自定义 BaseModel
- 事务使用 db.Transaction() 闭包形式
- 批量操作使用 CreateInBatches / 分批查询
- 默认软删除，物理删除需注释说明原因
- migration 文件按时间戳命名：{timestamp}_{description}.sql，必须有 UP 和 DOWN

## API 设计

### HTTP (RESTful)
- URL：kebab-case，复数资源名，版本号在 URL 中（/api/v1/users）
- 分页：page + page_size，响应包含 total
- 请求/响应 DTO 与数据库 Model 分离


## 测试

- 表驱动测试为默认模式
- Mock 使用 gomock 或 testify/mock
- repository 层测试连真实数据库（测试库），不 mock
- 测试命名：Test{函数名}_{场景}_{预期结果}
- 单元测试命名清晰，覆盖正常路径和边界路径。
- 修复 bug 必须补回归测试。
- 未说明原因不得跳过失败测试直接提交。
- 对外 API 需包含基础集成测试或契约测试。

## Git 工作流

- 分支：main ← develop ← feature/xxx, fix/xxx, refactor/xxx
- Commit message：`{type}({scope}): {description}`（feat / fix / refactor / docs / test / chore / perf）

## Git 管理规则（强制）

1. 每次代码改动都必须通过 Git 管理，不允许“裸改”不留痕。
2. 每完成一个最小可验证步骤就提交一次（小步提交，便于回滚）。
3. 提交前必须：
- 运行相关测试
- 更新 `tasks.md` 与 `context.md`
4. 提交信息格式：
- `feat:` 新功能
- `fix:` 修复
- `refactor:` 重构
- `test:` 测试
- `docs:` 文档
5. 禁止一次提交混入无关改动。
6. 需要回滚时，优先使用 `git revert`，避免破坏历史。



## 工作流规则（强制）
每个任务都必须创建并维护以下 3 个文档：

1. `plan.md`
- 本任务要做什么
- 具体怎么做
- 风险点与验收标准

2. `context.md`
- 涉及的关键文件/模块
- 决策记录（为什么这样做）
- 当前进度与下一步交接说明

3. `tasks.md`
- 可勾选 checklist（执行清单）
- 每完成一步必须立刻更新

## 执行顺序（必须 TDD）
始终遵循以下顺序：

1. 先写/先改测试用例
2. 运行测试（在未实现前允许失败）
3. 编写实现代码
4. 再次运行测试（应通过）
5. 更新 `tasks.md` 和 `context.md`

除非用户明确批准，否则不得跳过“先测试后实现”。

## 会话交接规则
当上下文快满时，在结束当前会话前必须：

1. 更新 `context.md`，写清：
- 已完成内容
- 当前状态
- 明确的下一步动作

2. 更新 `tasks.md`：
- 勾选已完成项
- 补充剩余可执行项

下一个会话中，用户只需说：`继续`。  
助手必须基于文档无缝接续，不从零重新规划。

## 完成定义（Definition of Done）
一个任务完成必须同时满足：

1. 功能代码完成
2. 测试已新增/更新且全部通过
3. `tasks.md` 已按实际进度更新
4. `context.md` 已记录总结与下一步（如有）
5. 关键技术决策有文档记录

## 代码规范
1. 改动保持最小且聚焦当前任务
2. 命名清晰，函数尽量小而单一职责
3. 不做无关重构
4. 仅在必要处写简洁注释
5. 行为变更要同步更新文档

## 沟通规范
1. 回复简洁、面向执行
2. 按小步反馈进展
3. 遇阻塞时给出可选方案与权衡
4. 没有测试依据时，不得声称“已完成”

## Codex 交互偏好
- 用中文沟通
- 写代码前先说明设计思路和分层考虑
- 涉及架构决策时给出至少 2 个方案并对比优劣
- 修改现有代码前先阅读理解，说明改动影响范围
- 不要做超出需求范围的"优化"或"改进"
- 遇到不确定的业务逻辑，先问清楚再实现

## 配置与密钥规则
1. 禁止提交任何密钥、Token、密码到仓库。
2. 新增配置必须同步更新 `.env.example`。
3. 区分开发/测试/生产配置，禁止硬编码环境参数。

## API 变更规则
1. API 结构变更必须更新 `docs/api.md`。
2. 涉及兼容性变更时，优先保留旧字段/旧接口过渡。
3. 错误码与错误信息保持稳定、可检索。

## 发布与回滚规则
1. 每个可上线节点打 tag（如 `v0.1.0`）。
2. 发布失败优先回滚到上一个稳定 tag。
3. 回滚后必须在 `context.md` 记录原因、影响、后续动作。

## 禁止事项
1. 未经确认不得删除生产数据。
2. 不得执行高风险破坏性命令（如强制重置历史）。
3. 未记录决策，不做跨模块大改。


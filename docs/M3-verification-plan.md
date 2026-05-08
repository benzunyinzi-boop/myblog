# M3 端到端验证计划

## 验证目标

验证 M3 所有功能模块的完整性和正确性:
- Category / Tag CRUD
- Article 管理端 CRUD + 发布生命周期
- public 文章列表/详情
- Profile 读写
- 文件上传

## 验证环境

- MySQL: 本地 5.7.16,myblog 数据库
- Redis: 本地 6379 端口
- 服务端口: 18080(本地 8080 被占用)
- 管理员账号: admin / Admin@123

## 验证场景

### 1. 认证与授权

| # | 操作 | 预期 |
|---|---|---|
| 1.1 | POST /admin/auth/login(正确密码) | 200 + access_token + user |
| 1.2 | POST /admin/auth/login(错误密码) | 400 + code=10001 |
| 1.3 | GET /admin/ping(无 token) | 401 + code=10002 |
| 1.4 | GET /admin/ping(坏 token) | 200 + code=30005 |
| 1.5 | GET /admin/ping(有效 token) | 200 + pong + user_id |

### 2. 分类管理

| # | 操作 | 预期 |
|---|---|---|
| 2.1 | POST /admin/categories | 200 + id + name/slug |
| 2.2 | POST /admin/categories(重复 slug) | 200 + code=40002 |
| 2.3 | GET /admin/categories | 200 + items + total |
| 2.4 | PUT /admin/categories/:id | 200 + 更新后数据 |
| 2.5 | DELETE /admin/categories/:id | 200 + code=0 |
| 2.6 | GET /public/categories | 200 + items(无需 token) |

### 3. 标签管理

| # | 操作 | 预期 |
|---|---|---|
| 3.1 | POST /admin/tags | 200 + id + name/slug |
| 3.2 | POST /admin/tags(重复 slug) | 200 + code=40102 |
| 3.3 | GET /admin/tags | 200 + items + total |
| 3.4 | DELETE /admin/tags/:id | 200 + code=0 |
| 3.5 | GET /public/tags | 200 + items(无需 token) |

### 4. 文章管理

| # | 操作 | 预期 |
|---|---|---|
| 4.1 | POST /admin/articles(草稿) | 200 + id + status=draft |
| 4.2 | POST /admin/articles(重复 slug) | 200 + code=20002 |
| 4.3 | POST /admin/articles(无效分类) | 200 + code=40001 |
| 4.4 | POST /admin/articles(无效标签) | 200 + code=40101 |
| 4.5 | GET /admin/articles/:id | 200 + 详情 + tags |
| 4.6 | PUT /admin/articles/:id | 200 + 更新后数据 |
| 4.7 | POST /admin/articles/:id/publish | 200 + status=published + published_at |
| 4.8 | POST /admin/articles/:id/publish(重复) | 200 + code=20003 |
| 4.9 | POST /admin/articles/:id/unpublish | 200 + status=draft |
| 4.10 | POST /admin/articles/:id/unpublish(重复) | 200 + code=20003 |
| 4.11 | GET /admin/articles?status=draft | 200 + 草稿列表 |
| 4.12 | GET /admin/articles?category_id=X | 200 + 过滤结果 |
| 4.13 | GET /admin/articles?tag_id=X | 200 + 过滤结果 |
| 4.14 | GET /admin/articles?keyword=X | 200 + 搜索结果 |
| 4.15 | DELETE /admin/articles/:id | 200 + code=0 |

### 5. 公开文章接口

| # | 操作 | 预期 |
|---|---|---|
| 5.1 | GET /public/articles | 200 + 已发布文章列表(无需 token) |
| 5.2 | GET /public/articles?category_id=X | 200 + 过滤结果 |
| 5.3 | GET /public/articles?tag_id=X | 200 + 过滤结果 |
| 5.4 | GET /public/articles?keyword=X | 200 + 搜索结果 |
| 5.5 | GET /public/articles/:slug(已发布) | 200 + 详情 + content |
| 5.6 | GET /public/articles/:slug(草稿) | 200 + code=20001 |
| 5.7 | GET /public/articles/:slug(不存在) | 200 + code=20001 |
| 5.8 | 多次访问同一文章 | view_count 递增 |

### 6. 个人资料

| # | 操作 | 预期 |
|---|---|---|
| 6.1 | GET /public/profile(首次) | 200 + 空 profile |
| 6.2 | PUT /admin/profile | 200 + 更新后数据 |
| 6.3 | GET /public/profile | 200 + 已更新数据 |

### 7. 文件上传

| # | 操作 | 预期 |
|---|---|---|
| 7.1 | POST /admin/uploads(无 token) | 401 + code=10002 |
| 7.2 | POST /admin/uploads(有效文件) | 200 + url |
| 7.3 | POST /admin/uploads(超过 10MB) | 200 + code=10001 |
| 7.4 | GET /uploads/2026/05/08/xxx.jpg | 200 + 文件内容 |

### 8. 事务一致性

| # | 操作 | 预期 |
|---|---|---|
| 8.1 | 创建文章 + 标签关联 | article + article_tags 同时写入 |
| 8.2 | 更新文章标签 | article_tags 先删后插 |
| 8.3 | 删除文章 | article_tags 先清空,article 软删 |

### 9. 错误处理

| # | 操作 | 预期 |
|---|---|---|
| 9.1 | 参数校验失败 | 400 + code=10001 + 详细错误 |
| 9.2 | 资源不存在 | 200 + 对应错误码(20001/40001 等) |
| 9.3 | 重复资源 | 200 + 对应错误码(20002/40002 等) |
| 9.4 | 状态冲突 | 200 + code=20003 |

## 已验证场景

### Article CRUD(2026-05-08)

已在 `docs/M3-article-verification.md` 中记录:
- 创建草稿 ✓
- 发布/下线状态机 ✓
- 重复 slug 检测 ✓
- 关联分类/标签校验 ✓
- 事务一致性 ✓
- published_at 类型修复 ✓

### Category + Tag(2026-05-08)

已在开发过程中验证:
- CRUD 全链路 ✓
- 重复 slug 检测 ✓
- public 无需 token ✓

## 待验证场景

- [ ] public 文章列表/详情完整流程
- [ ] Profile 读写完整流程
- [ ] 文件上传完整流程
- [ ] 浏览计数递增
- [ ] 分页功能
- [ ] 多条件过滤组合

## 验证方式

### 手动验证

```bash
# 1. 启动服务
make build
MYBLOG_SERVER_PORT=18080 ./bin/myblog-server

# 2. 登录获取 token
TOKEN=$(curl -s -X POST http://127.0.0.1:18080/api/v1/admin/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"Admin@123"}' | \
  jq -r '.data.access_token')

# 3. 按场景逐一验证
# 参考 docs/api.md 中的接口文档
```

### 自动化测试(TODO)

- 集成测试覆盖关键场景
- 端到端测试脚本
- CI/CD 流水线集成

## 验证结论

**当前状态**: M3 核心功能已实现并通过编译,部分场景已验证通过。

**已完成**:
- ✅ 所有接口已实现并编译通过
- ✅ Article CRUD 端到端验证通过
- ✅ Category/Tag 开发期验证通过
- ✅ API 文档已完善

**待完成**:
- ⏳ 完整的端到端验证脚本
- ⏳ 自动化集成测试
- ⏳ 性能测试

## 下一步

1. 运行完整的端到端验证脚本
2. 补充集成测试用例
3. 进入 M4(前端开发)或 M5(部署)

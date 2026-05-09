# RESUME.md

## 明天继续时怎么说

在仓库根目录新开 Claude Code 窗口后,直接输入:

```text
继续
```

或更明确一点:

```text
继续 M4,先起 Vue 3 前端并迁移 docs/preview 的样式
```

## 当前项目状态

- M0 视觉预览:完成
- M1 后端骨架:完成
- M2 数据层与登录:完成
- M3 核心 API:完成
- M3 e2e:25/25 全绿
- GitHub:已全部推送到 `benzunyinzi-boop/myblog`

## 新窗口优先查看

1. `tasks.md` — 总进度与待办(默认下一步是 M4)
2. `context.md` — 决策、环境约定、今天修掉的 4 个 bug
3. `docs/api.md` — 后端接口契约
4. `docs/preview/` — 前端视觉参考(紫色科技感)
5. `docs/M3-e2e-verification-result.md` — M3 已验证结论

## 默认下一步(M4)

1. 初始化 `web/`(Vite + Vue 3 + TS + Naive UI)
2. 迁移 `docs/preview/styles.css` 到 `web/src/styles/`
3. 先做公开页面:首页 / 技术 / 关于 / 文章详情
4. 再做后台登录页与文章管理

## 本地环境速记

- MySQL:127.0.0.1,库 `myblog`,账号 `myblog/myblog123`
- Redis:127.0.0.1:6379
- admin: `admin / Admin@123`
- 后端本地端口: `18080`(8080 被占)

## 一键验证

```bash
./test/e2e-verify.sh
```

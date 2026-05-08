#!/bin/bash
# M3 端到端验证脚本
# 测试所有已实现的 API 接口

set -e

BASE="http://127.0.0.1:18080/api/v1"
PASS=0
FAIL=0

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log_pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((PASS++))
}

log_fail() {
    echo -e "${RED}✗${NC} $1"
    ((FAIL++))
}

echo "=== M3 端到端验证 ==="
echo

# 1. 登录获取 token
echo "1. 认证与授权"
TOKEN=$(curl -s -X POST $BASE/admin/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"Admin@123"}' | \
  python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["access_token"])' 2>/dev/null)

if [ -n "$TOKEN" ]; then
    log_pass "登录成功,获取 token"
else
    log_fail "登录失败"
    exit 1
fi

H_AUTH="Authorization: Bearer $TOKEN"
H_JSON="Content-Type: application/json"

# 2. 创建分类
echo
echo "2. 分类管理"
CAT=$(curl -s -X POST $BASE/admin/categories -H "$H_AUTH" -H "$H_JSON" \
  -d '{"name":"Golang","slug":"golang","description":"Go 语言","sort_order":1}')
CAT_ID=$(echo "$CAT" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["id"])' 2>/dev/null)

if [ -n "$CAT_ID" ]; then
    log_pass "创建分类成功 (id=$CAT_ID)"
else
    log_fail "创建分类失败"
fi

# 3. 创建标签
echo
echo "3. 标签管理"
TAG=$(curl -s -X POST $BASE/admin/tags -H "$H_AUTH" -H "$H_JSON" \
  -d '{"name":"并发","slug":"concurrency"}')
TAG_ID=$(echo "$TAG" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["id"])' 2>/dev/null)

if [ -n "$TAG_ID" ]; then
    log_pass "创建标签成功 (id=$TAG_ID)"
else
    log_fail "创建标签失败"
fi

# 4. 创建草稿文章
echo
echo "4. 文章管理"
ART=$(curl -s -X POST $BASE/admin/articles -H "$H_AUTH" -H "$H_JSON" \
  -d "{\"title\":\"Go 并发模式\",\"slug\":\"go-concurrency\",\"summary\":\"errgroup + context\",\"content\":\"正文内容\",\"category_id\":$CAT_ID,\"tag_ids\":[$TAG_ID]}")
ART_ID=$(echo "$ART" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["id"])' 2>/dev/null)

if [ -n "$ART_ID" ]; then
    log_pass "创建草稿文章成功 (id=$ART_ID)"
else
    log_fail "创建草稿文章失败"
fi

# 5. 发布文章
PUB=$(curl -s -X POST $BASE/admin/articles/$ART_ID/publish -H "$H_AUTH")
PUB_STATUS=$(echo "$PUB" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["status"])' 2>/dev/null)

if [ "$PUB_STATUS" = "published" ]; then
    log_pass "发布文章成功"
else
    log_fail "发布文章失败"
fi

# 6. public 文章列表
echo
echo "5. 公开接口"
PUB_LIST=$(curl -s "$BASE/public/articles")
PUB_TOTAL=$(echo "$PUB_LIST" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["total"])' 2>/dev/null)

if [ "$PUB_TOTAL" -ge 1 ]; then
    log_pass "public 文章列表返回 $PUB_TOTAL 篇"
else
    log_fail "public 文章列表为空"
fi

# 7. public 文章详情
PUB_DETAIL=$(curl -s "$BASE/public/articles/go-concurrency")
PUB_TITLE=$(echo "$PUB_DETAIL" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["title"])' 2>/dev/null)

if [ "$PUB_TITLE" = "Go 并发模式" ]; then
    log_pass "public 文章详情正确"
else
    log_fail "public 文章详情错误"
fi

# 8. 浏览计数
sleep 1
PUB_DETAIL2=$(curl -s "$BASE/public/articles/go-concurrency")
VIEW_COUNT=$(echo "$PUB_DETAIL2" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["view_count"])' 2>/dev/null)

if [ "$VIEW_COUNT" -ge 2 ]; then
    log_pass "浏览计数递增 (view_count=$VIEW_COUNT)"
else
    log_fail "浏览计数未递增"
fi

# 9. Profile 更新
echo
echo "6. 个人资料"
PROF=$(curl -s -X PUT $BASE/admin/profile -H "$H_AUTH" -H "$H_JSON" \
  -d '{"name":"测试用户","bio":"全栈开发","avatar":"","email":"test@example.com","github":"","twitter":"","linkedin":"","website":""}')
PROF_NAME=$(echo "$PROF" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["name"])' 2>/dev/null)

if [ "$PROF_NAME" = "测试用户" ]; then
    log_pass "更新 Profile 成功"
else
    log_fail "更新 Profile 失败"
fi

# 10. public Profile
PUB_PROF=$(curl -s "$BASE/public/profile")
PUB_PROF_NAME=$(echo "$PUB_PROF" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["name"])' 2>/dev/null)

if [ "$PUB_PROF_NAME" = "测试用户" ]; then
    log_pass "public Profile 正确"
else
    log_fail "public Profile 错误"
fi

# 11. 下线文章
echo
echo "7. 文章生命周期"
UNPUB=$(curl -s -X POST $BASE/admin/articles/$ART_ID/unpublish -H "$H_AUTH")
UNPUB_STATUS=$(echo "$UNPUB" | python3 -c 'import json,sys; print(json.load(sys.stdin)["data"]["status"])' 2>/dev/null)

if [ "$UNPUB_STATUS" = "draft" ]; then
    log_pass "下线文章成功"
else
    log_fail "下线文章失败"
fi

# 12. 删除文章
DEL=$(curl -s -X DELETE $BASE/admin/articles/$ART_ID -H "$H_AUTH")
DEL_CODE=$(echo "$DEL" | python3 -c 'import json,sys; print(json.load(sys.stdin)["code"])' 2>/dev/null)

if [ "$DEL_CODE" = "0" ]; then
    log_pass "删除文章成功"
else
    log_fail "删除文章失败"
fi

# 总结
echo
echo "=== 验证结果 ==="
echo "通过: $PASS"
echo "失败: $FAIL"

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过${NC}"
    exit 0
else
    echo -e "${RED}✗ 有 $FAIL 个测试失败${NC}"
    exit 1
fi

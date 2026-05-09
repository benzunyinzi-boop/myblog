#!/bin/bash
# M3 端到端验证脚本
# 用法:
#   1. 启动服务:MYBLOG_SERVER_PORT=18080 ./bin/myblog-server &
#   2. 运行:./test/e2e-verify.sh

BASE="http://127.0.0.1:18080/api/v1"
PASS=0
FAIL=0

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

log_pass() { echo -e "${GREEN}✓${NC} $1"; PASS=$((PASS + 1)); }
log_fail() { echo -e "${RED}✗${NC} $1 ${YELLOW}[$2]${NC}"; FAIL=$((FAIL + 1)); }

# json_get <path>:从 stdin 读 JSON,按点号路径取值。失败返回空。
json_get() {
    python3 -c "
import json, sys
try:
    d = json.loads(sys.stdin.read())
    v = d
    for k in '$1'.split('.'):
        v = v[int(k)] if k.isdigit() else v[k]
    print(v)
except Exception:
    pass
" 2>/dev/null
}

echo "=== M3 端到端验证 ==="
echo

# ---- 1. 认证 ----
echo "[1/7] 认证与授权"
LOGIN=$(curl -s -X POST $BASE/admin/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"username":"admin","password":"Admin@123"}')
TOKEN=$(echo "$LOGIN" | json_get data.access_token)
if [ -n "$TOKEN" ]; then
    log_pass "登录成功"
else
    log_fail "登录失败" "$LOGIN"
    echo "终止:无 token 后续测试无法进行"
    exit 1
fi

# 错误密码
WRONG=$(curl -s -X POST $BASE/admin/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"username":"admin","password":"wrong1"}')
CODE=$(echo "$WRONG" | json_get code)
if [ "$CODE" = "30002" ] || [ "$CODE" = "10001" ]; then
    log_pass "错误密码被拒绝 (code=$CODE)"
else
    log_fail "错误密码验证失败" "code=$CODE"
fi

# 无 token 访问受保护接口
NOAUTH=$(curl -s $BASE/admin/ping)
CODE=$(echo "$NOAUTH" | json_get code)
if [ "$CODE" = "10002" ]; then
    log_pass "无 token 被拒绝 (code=10002)"
else
    log_fail "无 token 验证失败" "code=$CODE"
fi

H_AUTH="Authorization: Bearer $TOKEN"
H_JSON="Content-Type: application/json"

# ---- 2. 分类 ----
echo
echo "[2/7] 分类管理"
CAT=$(curl -s -X POST $BASE/admin/categories -H "$H_AUTH" -H "$H_JSON" \
    -d '{"name":"Golang","slug":"golang","description":"Go 语言","sort_order":1}')
CAT_ID=$(echo "$CAT" | json_get data.id)
if [ -n "$CAT_ID" ]; then
    log_pass "创建分类成功 (id=$CAT_ID)"
else
    log_fail "创建分类失败" "$CAT"
fi

# 重复 slug
DUP_CAT=$(curl -s -X POST $BASE/admin/categories -H "$H_AUTH" -H "$H_JSON" \
    -d '{"name":"Go2","slug":"golang","sort_order":2}')
CODE=$(echo "$DUP_CAT" | json_get code)
if [ "$CODE" = "40002" ]; then
    log_pass "重复 slug 被拒绝 (code=40002)"
else
    log_fail "重复 slug 验证失败" "code=$CODE"
fi

# public 访问(无需 token)
PUB_CATS=$(curl -s $BASE/public/categories)
TOTAL=$(echo "$PUB_CATS" | json_get data.total)
if [ "$TOTAL" = "1" ]; then
    log_pass "public 分类列表正确 (total=1)"
else
    log_fail "public 分类列表错误" "total=$TOTAL"
fi

# ---- 3. 标签 ----
echo
echo "[3/7] 标签管理"
TAG=$(curl -s -X POST $BASE/admin/tags -H "$H_AUTH" -H "$H_JSON" \
    -d '{"name":"并发","slug":"concurrency"}')
TAG_ID=$(echo "$TAG" | json_get data.id)
if [ -n "$TAG_ID" ]; then
    log_pass "创建标签成功 (id=$TAG_ID)"
else
    log_fail "创建标签失败" "$TAG"
fi

TAG2=$(curl -s -X POST $BASE/admin/tags -H "$H_AUTH" -H "$H_JSON" \
    -d '{"name":"性能","slug":"performance"}')
TAG2_ID=$(echo "$TAG2" | json_get data.id)
if [ -n "$TAG2_ID" ]; then
    log_pass "创建第二个标签成功 (id=$TAG2_ID)"
else
    log_fail "创建第二个标签失败" "$TAG2"
fi

# ---- 4. 文章 ----
echo
echo "[4/7] 文章管理"
ART=$(curl -s -X POST $BASE/admin/articles -H "$H_AUTH" -H "$H_JSON" \
    -d "{\"title\":\"Go 并发模式\",\"slug\":\"go-concurrency\",\"summary\":\"errgroup + context\",\"content\":\"正文\",\"category_id\":$CAT_ID,\"tag_ids\":[$TAG_ID,$TAG2_ID]}")
ART_ID=$(echo "$ART" | json_get data.id)
if [ -n "$ART_ID" ]; then
    log_pass "创建草稿文章成功 (id=$ART_ID)"
else
    log_fail "创建文章失败" "$ART"
fi

# 无效分类
BAD_CAT=$(curl -s -X POST $BASE/admin/articles -H "$H_AUTH" -H "$H_JSON" \
    -d '{"title":"x","slug":"bad-cat","content":"x","category_id":99999}')
CODE=$(echo "$BAD_CAT" | json_get code)
if [ "$CODE" = "40001" ]; then
    log_pass "无效分类被拒绝 (code=40001)"
else
    log_fail "无效分类验证失败" "code=$CODE"
fi

# 重复 slug
DUP_ART=$(curl -s -X POST $BASE/admin/articles -H "$H_AUTH" -H "$H_JSON" \
    -d "{\"title\":\"dup\",\"slug\":\"go-concurrency\",\"content\":\"x\",\"category_id\":$CAT_ID}")
CODE=$(echo "$DUP_ART" | json_get code)
if [ "$CODE" = "20002" ]; then
    log_pass "重复文章 slug 被拒绝 (code=20002)"
else
    log_fail "重复文章 slug 验证失败" "code=$CODE"
fi

# 发布
PUB=$(curl -s -X POST $BASE/admin/articles/$ART_ID/publish -H "$H_AUTH")
STATUS=$(echo "$PUB" | json_get data.status)
if [ "$STATUS" = "published" ]; then
    log_pass "发布文章成功 (status=published)"
else
    log_fail "发布文章失败" "status=$STATUS"
fi

# 重复发布
REPUB=$(curl -s -X POST $BASE/admin/articles/$ART_ID/publish -H "$H_AUTH")
CODE=$(echo "$REPUB" | json_get code)
if [ "$CODE" = "20003" ]; then
    log_pass "重复发布被拒绝 (code=20003)"
else
    log_fail "重复发布验证失败" "code=$CODE"
fi

# ---- 5. public 文章 ----
echo
echo "[5/7] 公开文章接口"
PUB_LIST=$(curl -s "$BASE/public/articles")
TOTAL=$(echo "$PUB_LIST" | json_get data.total)
if [ "$TOTAL" = "1" ]; then
    log_pass "public 文章列表正确 (total=1)"
else
    log_fail "public 文章列表错误" "total=$TOTAL"
fi

# 按分类过滤
FILTER_CAT=$(curl -s "$BASE/public/articles?category_id=$CAT_ID")
TOTAL=$(echo "$FILTER_CAT" | json_get data.total)
if [ "$TOTAL" = "1" ]; then
    log_pass "按分类过滤正确 (total=1)"
else
    log_fail "按分类过滤错误" "total=$TOTAL"
fi

# 按标签过滤
FILTER_TAG=$(curl -s "$BASE/public/articles?tag_id=$TAG_ID")
TOTAL=$(echo "$FILTER_TAG" | json_get data.total)
if [ "$TOTAL" = "1" ]; then
    log_pass "按标签过滤正确 (total=1)"
else
    log_fail "按标签过滤错误" "total=$TOTAL"
fi

# 关键词
FILTER_KW=$(curl -s "$BASE/public/articles?keyword=errgroup")
TOTAL=$(echo "$FILTER_KW" | json_get data.total)
if [ "$TOTAL" = "1" ]; then
    log_pass "关键词搜索正确 (total=1)"
else
    log_fail "关键词搜索错误" "total=$TOTAL"
fi

# slug 详情
DETAIL=$(curl -s "$BASE/public/articles/go-concurrency")
TITLE=$(echo "$DETAIL" | json_get data.title)
if [ "$TITLE" = "Go 并发模式" ]; then
    log_pass "public 文章详情正确"
else
    log_fail "public 文章详情错误" "title=$TITLE"
fi

# 浏览计数(异步,等一下)
sleep 1
DETAIL2=$(curl -s "$BASE/public/articles/go-concurrency")
sleep 1
DETAIL3=$(curl -s "$BASE/public/articles/go-concurrency")
VIEW=$(echo "$DETAIL3" | json_get data.view_count)
if [ -n "$VIEW" ] && [ "$VIEW" -ge 1 ]; then
    log_pass "浏览计数递增 (view_count=$VIEW)"
else
    log_fail "浏览计数未递增" "view_count=$VIEW"
fi

# 不存在的 slug
NOT_FOUND=$(curl -s "$BASE/public/articles/not-exist-xxx")
CODE=$(echo "$NOT_FOUND" | json_get code)
if [ "$CODE" = "20001" ]; then
    log_pass "不存在文章返回 404 (code=20001)"
else
    log_fail "不存在文章验证失败" "code=$CODE"
fi

# ---- 6. Profile ----
echo
echo "[6/7] 个人资料"
PROF=$(curl -s -X PUT $BASE/admin/profile -H "$H_AUTH" -H "$H_JSON" \
    -d '{"name":"Yinyin","bio":"十年 Go 后端","avatar":"","email":"me@example.com","github":"https://github.com/yy","twitter":"","linkedin":"","website":""}')
NAME=$(echo "$PROF" | json_get data.name)
if [ "$NAME" = "Yinyin" ]; then
    log_pass "更新 Profile 成功"
else
    log_fail "更新 Profile 失败" "$PROF"
fi

PUB_PROF=$(curl -s "$BASE/public/profile")
NAME=$(echo "$PUB_PROF" | json_get data.name)
if [ "$NAME" = "Yinyin" ]; then
    log_pass "public Profile 正确"
else
    log_fail "public Profile 错误" "name=$NAME"
fi

# ---- 7. 生命周期 & 清理 ----
echo
echo "[7/7] 文章生命周期"
UNPUB=$(curl -s -X POST $BASE/admin/articles/$ART_ID/unpublish -H "$H_AUTH")
STATUS=$(echo "$UNPUB" | json_get data.status)
if [ "$STATUS" = "draft" ]; then
    log_pass "下线文章成功 (status=draft)"
else
    log_fail "下线文章失败" "status=$STATUS"
fi

DEL=$(curl -s -X DELETE $BASE/admin/articles/$ART_ID -H "$H_AUTH")
CODE=$(echo "$DEL" | json_get code)
if [ "$CODE" = "0" ]; then
    log_pass "删除文章成功"
else
    log_fail "删除文章失败" "code=$CODE"
fi

# 删除后再查询
GONE=$(curl -s "$BASE/admin/articles/$ART_ID" -H "$H_AUTH")
CODE=$(echo "$GONE" | json_get code)
if [ "$CODE" = "20001" ]; then
    log_pass "已删除文章返回 404 (code=20001)"
else
    log_fail "已删除文章验证失败" "code=$CODE"
fi

# ---- 总结 ----
echo
echo "=== 验证结果 ==="
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过${NC}"
    exit 0
else
    echo -e "${RED}✗ 有 $FAIL 个测试失败${NC}"
    exit 1
fi

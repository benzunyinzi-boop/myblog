<template>
  <section class="section-gap container tech-page">
    <div class="section-headline stack-headline">
      <div>
        <span class="section-chip">taxonomy / tech</span>
        <h1 class="page-title">技术地图</h1>
        <p class="page-copy">按分类浏览 Go、MySQL、Redis、MQ、AI 与部署实践，每个分类最多展示 10 篇最新文章。</p>
      </div>
    </div>

    <div v-if="loading" class="tech-loading">
      <div v-for="i in 3" :key="i" class="skeleton-card glass-card" />
    </div>

    <div v-else-if="!groups.length" class="state-box glass-card">暂时还没有已发布的分类文章。</div>

    <div v-else class="cat-groups">
      <section v-for="group in groups" :key="group.id" class="cat-group">
        <header class="cat-group-head">
          <div class="cat-icons">
            <span
              v-for="(emoji, i) in group.icons"
              :key="i"
              class="cat-icon"
              :style="{ '--i': i }"
            >{{ emoji }}</span>
          </div>
          <div class="cat-group-title-wrap">
            <span class="cat-group-slug">{{ group.slug }}</span>
            <h2 class="cat-group-title">{{ group.name }}</h2>
          </div>
          <span class="cat-group-count">{{ group.total }} 篇</span>
        </header>

        <p v-if="group.description" class="cat-group-desc">{{ group.description }}</p>

        <div v-if="group.items.length" class="cat-articles">
          <OrganicCard
            v-for="article in group.items.slice(0, 10)"
            :key="article.id"
            :article="article"
          />
        </div>
        <div v-else class="state-box glass-card subtle">这个分类下还没有已发布文章。</div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchCategories, type Category } from '../../api/meta'
import { fetchPublicArticles, type ArticleSummary } from '../../api/article'
import OrganicCard from '../../components/OrganicCard.vue'

type CategoryGroup = Category & {
  items: ArticleSummary[]
  total: number
  icons: string[]
}

const ICON_MAP: Record<string, string[]> = {
  go: ['🐹', '💨', '⚡'],
  golang: ['🐹', '💨', '⚡'],
  mysql: ['🐬', '🌊', '📊'],
  database: ['🐬', '🗄️', '⚙️'],
  db: ['🐬', '🗄️', '⚙️'],
  redis: ['🟥', '⚡', '🔥'],
  cache: ['🟥', '💾', '⚡'],
  mq: ['🪶', '📡', '📮'],
  kafka: ['🪶', '📡', '🌊'],
  rabbitmq: ['🐰', '📡', '🪶'],
  ai: ['✨', '🧠', '🤖'],
  llm: ['✨', '🧠', '🤖'],
  cloud: ['☁️', '🌩️', '⛅'],
  devops: ['🔧', '🚀', '⚙️'],
  deploy: ['🚀', '📦', '⚙️'],
  arch: ['🏛️', '🧭', '🔱'],
  architecture: ['🏛️', '🧭', '🔱']
}

function pickIcons(slug: string): string[] {
  return ICON_MAP[slug?.toLowerCase()] ?? ['◆', '◇', '◈']
}

const groups = ref<CategoryGroup[]>([])
const loading = ref(true)

async function loadAll() {
  loading.value = true
  try {
    const catResp = await fetchCategories()
    const cats = catResp.data.items ?? []
    if (!cats.length) {
      groups.value = []
      return
    }
    const results = await Promise.all(
      cats.map(async (c) => {
        try {
          const r = await fetchPublicArticles({ category_id: c.id, page_size: 10 })
          return {
            ...c,
            items: r.data.items ?? [],
            total: r.data.total ?? 0,
            icons: pickIcons(c.slug)
          } as CategoryGroup
        } catch {
          return { ...c, items: [], total: 0, icons: pickIcons(c.slug) } as CategoryGroup
        }
      })
    )
    groups.value = results
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.tech-page {
  position: relative;
}

.tech-loading {
  display: grid;
  gap: 20px;
}

.cat-groups {
  display: grid;
  gap: 64px;
  margin-top: 40px;
}

.cat-group {
  display: grid;
  gap: 18px;
}

.cat-group-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 14px;
  border-bottom: 1px dashed rgba(196, 181, 253, 0.18);
}

.cat-icons {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  flex-shrink: 0;
}

.cat-icon {
  font-size: 28px;
  line-height: 1;
  display: inline-block;
  animation: cat-bob 2.8s ease-in-out infinite;
  animation-delay: calc(var(--i) * 220ms);
  filter: drop-shadow(0 0 8px rgba(139, 92, 246, 0.4));
}

.cat-group:hover .cat-icon {
  animation-duration: 1.4s;
}

@keyframes cat-bob {
  0%, 100% { transform: translateY(0); }
  50%      { transform: translateY(-6px); }
}

.cat-group-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  margin-left: 6px;
}

.cat-group-slug {
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.cat-group-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 30px;
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

.cat-group-count {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--neon-cyan);
  letter-spacing: 0.08em;
  padding-bottom: 6px;
}

.cat-group-desc {
  margin: 0;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.75;
}

.cat-articles {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 22px;
}

@media (max-width: 1080px) {
  .cat-articles { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 720px) {
  .cat-articles { grid-template-columns: 1fr; }
  .cat-group-head { flex-wrap: wrap; }
  .cat-group-title { font-size: 24px; }
  .cat-icon { font-size: 22px; }
}

.state-box.subtle {
  padding: 18px 20px;
  font-size: 13px;
  color: var(--text-muted);
}
</style>

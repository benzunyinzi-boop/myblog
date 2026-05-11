<template>
  <section class="hero section-gap">
    <div class="hero-glow" />
    <div class="container hero-inner">
      <div class="hero-left">
        <TypingTerminal />

        <div class="hero-actions">
          <RouterLink class="primary-button" to="/tech">浏览技术栏目</RouterLink>
          <RouterLink class="ghost-button" to="/about">关于我</RouterLink>
        </div>

        <div class="hero-metrics">
          <div class="metric-card glass-card">
            <span>published</span>
            <strong>{{ items.length }}</strong>
          </div>
          <div class="metric-card glass-card">
            <span>categories</span>
            <strong>{{ categories.length }}</strong>
          </div>
          <div class="metric-card glass-card">
            <span>status</span>
            <strong>stable</strong>
          </div>
        </div>
      </div>

      <div class="hero-orbit">
        <!-- 静态虚线轨道 -->
        <div class="orbit-ring orbit-ring-lg" />
        <div class="orbit-ring orbit-ring-md" />
        <div class="orbit-ring orbit-ring-sm" />

        <!-- 中心太阳(带透明感与光芒) -->
        <div class="sun">
          <div class="sun-glow" />
          <div class="sun-core" />
          <div class="sun-ray sun-ray-1" />
          <div class="sun-ray sun-ray-2" />
          <div class="sun-ray sun-ray-3" />
          <div class="sun-ray sun-ray-4" />
        </div>

        <!-- 外轨道行星组(顺时针 42s) -->
        <div class="planet-track planet-track-lg">
          <div class="planet planet-lg" style="--angle: 45deg;">🪐</div>
          <div class="planet planet-lg" style="--angle: 225deg;">🌍</div>
        </div>

        <!-- 中轨道行星组(逆时针 28s) -->
        <div class="planet-track planet-track-md">
          <div class="planet planet-md" style="--angle: 0deg;">🌕</div>
          <div class="planet planet-md" style="--angle: 120deg;">🌑</div>
          <div class="planet planet-md" style="--angle: 240deg;">⭐</div>
        </div>

        <!-- 内轨道行星组(顺时针 18s) -->
        <div class="planet-track planet-track-sm">
          <div class="planet planet-sm" style="--angle: 30deg;">☄️</div>
          <div class="planet planet-sm" style="--angle: 120deg;">🌠</div>
          <div class="planet planet-sm" style="--angle: 210deg;">💫</div>
          <div class="planet planet-sm" style="--angle: 300deg;">✨</div>
        </div>
      </div>
    </div>
  </section>

  <!-- 按分类分组:每个分类最多 3 篇文章 -->
  <section class="section-gap container">
    <div class="section-headline">
      <div>
        <span class="section-chip">catalog / by category</span>
        <h2>分类阅读</h2>
      </div>
      <RouterLink class="article-link" to="/tech">查看全部分类 →</RouterLink>
    </div>

    <div v-if="pending" class="state-box glass-card">正在加载文章...</div>
    <div v-else-if="categories.length === 0" class="state-box glass-card">
      还没有分类,先去后台创建分类和文章。
    </div>

    <div v-else class="category-groups">
      <section
        v-for="category in categories"
        :key="category.id"
        class="category-group"
      >
        <header class="category-head">
          <div>
            <span class="stack-slug">/ {{ category.slug }}</span>
            <h3 class="category-title">{{ category.name }}</h3>
          </div>
          <RouterLink class="article-link" :to="`/tech`">more →</RouterLink>
        </header>

        <div
          v-if="articlesByCategory(category.id).length === 0"
          class="state-box glass-card subtle"
        >
          这个分类下暂时还没有已发布文章。
        </div>

        <div v-else class="category-cards">
          <article
            v-for="article in articlesByCategory(category.id)"
            :key="article.slug"
            class="glass-card article-card"
          >
            <div class="article-tags">
              <span v-for="tag in article.tags.slice(0, 2)" :key="tag.id" class="tech-pill">
                {{ tag.name }}
              </span>
            </div>
            <h4>{{ article.title }}</h4>
            <p>{{ article.summary || '这是一篇还没有摘要的文章。' }}</p>
            <div class="card-footline">
              <span>{{ article.view_count }} views</span>
              <RouterLink class="article-link" :to="`/blog/${article.slug}`">
                阅读全文 →
              </RouterLink>
            </div>
          </article>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsyncState } from '@vueuse/core'
import { fetchPublicArticles } from '../../api/article'
import { fetchCategories } from '../../api/meta'
import TypingTerminal from '../../components/TypingTerminal.vue'

const { state: articleState, isLoading: pending } = useAsyncState(
  () => fetchPublicArticles({ page: 1, page_size: 60 }),
  null
)
const { state: categoryState } = useAsyncState(() => fetchCategories(), null)

const items = computed(() => articleState.value?.data.items ?? [])
const categories = computed(() => categoryState.value?.data.items ?? [])

function articlesByCategory(categoryId: number) {
  return items.value.filter((a) => a.category_id === categoryId).slice(0, 3)
}
</script>

<template>
  <section class="section-gap container">
    <div class="section-headline stack-headline">
      <div>
        <span class="section-chip">taxonomy / tech</span>
        <h1 class="page-title">技术地图</h1>
        <p class="page-copy">按分类浏览 Go、MySQL、Redis、MQ 与部署实践,这会是前台最核心的内容入口。</p>
      </div>
    </div>

    <div v-if="categories.length" class="stack-board">
      <button
        v-for="category in categories"
        :key="category.id"
        class="stack-item"
        :class="{ active: currentCategoryId === category.id }"
        @click="selectCategory(category.id)"
      >
        <span class="stack-slug">{{ category.slug }}</span>
        <strong>{{ category.name }}</strong>
        <small>{{ category.description || '工程实践笔记' }}</small>
      </button>
    </div>
    <div v-else class="state-box glass-card">还没有分类数据,先去后台创建分类。</div>

    <div class="section-subline">
      <span class="section-chip subtle">filtered / articles</span>
      <span class="muted-copy">当前分类下的已发布文章</span>
    </div>

    <div v-if="pendingArticles" class="state-box glass-card">正在加载文章...</div>
    <div v-else-if="items.length === 0" class="state-box glass-card">这个分类下暂时还没有已发布文章。</div>

    <div v-else class="feature-grid">
      <article v-for="article in items" :key="article.slug" class="glass-card article-card">
        <div class="article-tags">
          <span v-for="tag in article.tags.slice(0, 2)" :key="tag.id" class="tech-pill">{{ tag.name }}</span>
        </div>
        <h3>{{ article.title }}</h3>
        <p>{{ article.summary || '这是一篇还没有摘要的文章。' }}</p>
        <div class="card-footline">
          <span>{{ article.view_count }} views</span>
          <RouterLink class="article-link" :to="`/blog/${article.slug}`">阅读全文 →</RouterLink>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsyncState } from '@vueuse/core'
import { fetchCategories } from '../../api/meta'
import { fetchPublicArticles } from '../../api/article'

const currentCategoryId = ref<number | null>(null)

const { state: categoriesState } = useAsyncState(() => fetchCategories(), null)
const categories = computed(() => categoriesState.value?.data.items ?? [])

const {
  state: articlesState,
  isLoading: pendingArticles,
  execute
} = useAsyncState(
  () => fetchPublicArticles(currentCategoryId.value ? { category_id: currentCategoryId.value } : {}),
  null,
  { immediate: false }
)

const items = computed(() => articlesState.value?.data.items ?? [])

watch(categories, (value) => {
  if (!currentCategoryId.value && value.length > 0) {
    currentCategoryId.value = value[0].id
    execute()
  }
}, { immediate: true })

function selectCategory(id: number) {
  if (currentCategoryId.value === id) {
    return
  }
  currentCategoryId.value = id
  execute()
}
</script>

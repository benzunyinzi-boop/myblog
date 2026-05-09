<template>
  <section class="section-gap container">
    <div class="section-headline stack-headline">
      <div>
        <span class="section-chip">taxonomy / tech</span>
        <h1 class="page-title">技术地图</h1>
        <p class="page-copy">按分类浏览你沉淀下来的 Go、MySQL、Redis、MQ 和部署经验。</p>
      </div>
    </div>

    <div class="stack-board">
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

    <div class="feature-grid">
      <article v-for="article in items" :key="article.slug" class="glass-card article-card">
        <div class="article-tags">
          <span v-for="tag in article.tags.slice(0, 2)" :key="tag.id" class="tech-pill">{{ tag.name }}</span>
        </div>
        <h3>{{ article.title }}</h3>
        <p>{{ article.summary }}</p>
        <RouterLink class="article-link" :to="`/blog/${article.slug}`">阅读全文 →</RouterLink>
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

const { state: articlesState, execute } = useAsyncState(
  () => fetchPublicArticles(currentCategoryId.value ? { category_id: currentCategoryId.value } : {}),
  null,
  { immediate: true }
)

const items = computed(() => articlesState.value?.data.items ?? [])

watch(categories, (value) => {
  if (!currentCategoryId.value && value.length > 0) {
    currentCategoryId.value = value[0].id
    execute()
  }
})

function selectCategory(id: number) {
  currentCategoryId.value = id
  execute()
}
</script>

<template>
  <section class="section-gap container article-page">
    <div v-if="article" class="article-layout">
      <article class="glass-card article-shell">
        <div class="article-topline">{{ article.slug }} · {{ article.view_count }} views</div>
        <h1 class="page-title article-heading">{{ article.title }}</h1>
        <p class="page-copy article-summary">{{ article.summary }}</p>

        <div class="article-tag-row">
          <span v-for="tag in article.tags" :key="tag.id" class="tech-pill">{{ tag.name }}</span>
        </div>

        <div class="markdown-shell">
          <pre>{{ article.content }}</pre>
        </div>
      </article>

      <aside class="glass-card article-aside">
        <span class="section-chip">article / metadata</span>
        <dl>
          <div>
            <dt>Published</dt>
            <dd>{{ article.published_at || '-' }}</dd>
          </div>
          <div>
            <dt>Category ID</dt>
            <dd>{{ article.category_id }}</dd>
          </div>
          <div>
            <dt>Status</dt>
            <dd>{{ article.status }}</dd>
          </div>
        </dl>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAsyncState } from '@vueuse/core'
import { fetchPublicArticleBySlug } from '../../api/article'

const route = useRoute()
const slug = computed(() => String(route.params.slug || ''))
const { state } = useAsyncState(() => fetchPublicArticleBySlug(slug.value), null)
const article = computed(() => state.value?.data ?? null)
</script>

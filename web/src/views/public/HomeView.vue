<template>
  <section class="hero section-gap">
    <div class="hero-glow" />
    <div class="container hero-inner">
      <div>
        <div class="eyebrow"><span class="status-dot" />online · backend · decade log</div>
        <h1 class="hero-title">
          记录系统的边界,<br />
          也记录<span class="gradient-text">自己</span>的成长。
        </h1>
        <p class="hero-subtitle">
          这是一个带着霓虹感的工程日志。十年后端经验,分布式系统、数据库、缓存与部署实践,都会在这里慢慢沉淀。
        </p>
        <div class="hero-actions">
          <RouterLink class="primary-button" to="/tech">探索技术栈</RouterLink>
          <RouterLink class="ghost-button" to="/about">关于我</RouterLink>
        </div>
      </div>

      <div class="hero-orbit">
        <div class="orbit orbit-lg" />
        <div class="orbit orbit-md" />
        <div class="orbit orbit-sm" />
        <div class="core" />
      </div>
    </div>
  </section>

  <section class="section-gap container">
    <div class="section-headline">
      <div>
        <span class="section-chip">latest / articles</span>
        <h2>最近更新</h2>
      </div>
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
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsyncState } from '@vueuse/core'
import { fetchPublicArticles } from '../../api/article'

const { state } = useAsyncState(() => fetchPublicArticles({ page: 1, page_size: 6 }), null)

const items = computed(() => state.value?.data.items ?? [])
</script>

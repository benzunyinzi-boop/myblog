<template>
  <section class="hero section-gap">
    <div class="hero-glow" />
    <div class="container hero-inner">
      <div>
        <div class="eyebrow"><span class="status-dot" />online · backend · decade log</div>
        <h1 class="hero-title">
          把分布式系统的锋利,<br />
          写成<span class="gradient-text">可读的经验</span>。
        </h1>
        <p class="hero-subtitle">
          这是一个偏工程化、也偏审美化的个人站点。十年后端开发经验,会在这里被整理成文章、分类和长期可复用的思考。
        </p>
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
        <!-- 外轨道 - 2颗行星 -->
        <div class="orbit orbit-lg">
          <div class="planet planet-purple" style="--angle: 45deg; --size: 24px;" />
          <div class="planet planet-blue" style="--angle: 225deg; --size: 20px;" />
        </div>
        <!-- 中轨道 - 3颗行星 -->
        <div class="orbit orbit-md">
          <div class="planet planet-cyan" style="--angle: 0deg; --size: 18px;" />
          <div class="planet planet-pink" style="--angle: 120deg; --size: 16px;" />
          <div class="planet planet-green" style="--angle: 240deg; --size: 14px;" />
        </div>
        <!-- 内轨道 - 4颗小行星 -->
        <div class="orbit orbit-sm">
          <div class="planet planet-orange" style="--angle: 30deg; --size: 12px;" />
          <div class="planet planet-teal" style="--angle: 120deg; --size: 10px;" />
          <div class="planet planet-violet" style="--angle: 210deg; --size: 11px;" />
          <div class="planet planet-amber" style="--angle: 300deg; --size: 9px;" />
        </div>
        <!-- 中心恒星 -->
        <div class="core" />
      </div>
    </div>
  </section>

  <section class="section-gap container">
    <div class="section-headline">
      <div>
        <span class="section-chip">focus / taxonomy</span>
        <h2>技术主轴</h2>
      </div>
      <RouterLink class="article-link" to="/tech">查看全部分类 →</RouterLink>
    </div>

    <div class="stack-board compact">
      <button
        v-for="category in categories.slice(0, 4)"
        :key="category.id"
        class="stack-item active"
      >
        <span class="stack-slug">{{ category.slug }}</span>
        <strong>{{ category.name }}</strong>
        <small>{{ category.description || '工程实践笔记' }}</small>
      </button>
    </div>
  </section>

  <section class="section-gap container">
    <div class="section-headline">
      <div>
        <span class="section-chip">latest / articles</span>
        <h2>最近更新</h2>
      </div>
    </div>

    <div v-if="pending" class="state-box glass-card">正在加载最近文章...</div>
    <div v-else-if="items.length === 0" class="state-box glass-card">还没有已发布文章,可以先去后台创建一篇。</div>

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
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsyncState } from '@vueuse/core'
import { fetchPublicArticles } from '../../api/article'
import { fetchCategories } from '../../api/meta'

const { state: articleState, isLoading: pending } = useAsyncState(() => fetchPublicArticles({ page: 1, page_size: 6 }), null)
const { state: categoryState } = useAsyncState(() => fetchCategories(), null)

const items = computed(() => articleState.value?.data.items ?? [])
const categories = computed(() => categoryState.value?.data.items ?? [])
</script>

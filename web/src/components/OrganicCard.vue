<template>
  <RouterLink :to="`/blog/${article.slug}`" class="organic-card">
    <div class="organic-content">
      <div class="organic-tags">
        <span v-for="t in article.tags?.slice(0, 2) ?? []" :key="t.id" class="tech-pill">{{ t.name }}</span>
      </div>
      <h3 class="organic-title">{{ article.title }}</h3>
      <p class="organic-summary">{{ article.summary || '这是一篇还没有摘要的文章。' }}</p>
      <div class="organic-foot">
        <span>{{ article.view_count ?? 0 }} views</span>
        <span class="organic-link">阅读全文 →</span>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { ArticleSummary } from '../api/article'

defineProps<{ article: ArticleSummary }>()
</script>

<style scoped>
.organic-card {
  position: relative;
  display: block;
  padding: 26px 24px 22px;
  border-radius: 16px;
  color: inherit;
  text-decoration: none;
  background: linear-gradient(160deg, rgba(22, 14, 42, 0.75), rgba(14, 8, 30, 0.92));
  border: 1px solid rgba(196, 181, 253, 0.14);
  transition: transform 320ms cubic-bezier(0.2, 0.8, 0.2, 1), border-color 280ms ease, box-shadow 280ms ease;
  overflow: hidden;
}

.organic-card::before {
  content: '';
  position: absolute;
  inset: -2px;
  border-radius: 18px;
  padding: 2px;
  background: conic-gradient(
    from 0deg,
    transparent 0%,
    rgba(34, 211, 238, 0.7) 25%,
    rgba(139, 92, 246, 0.7) 50%,
    rgba(34, 211, 238, 0.7) 75%,
    transparent 100%
  );
  -webkit-mask:
    linear-gradient(#fff 0 0) content-box,
    linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
          mask-composite: exclude;
  opacity: 0;
  transition: opacity 320ms ease;
  pointer-events: none;
}

.organic-card:hover::before {
  opacity: 1;
  animation: border-rotate 2.4s linear infinite;
}

.organic-card:hover {
  transform: translateY(-4px);
  border-color: rgba(139, 92, 246, 0.35);
  box-shadow: 0 12px 36px -12px rgba(139, 92, 246, 0.4);
}

@keyframes border-rotate {
  to {
    transform: rotate(360deg);
  }
}

.organic-content {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 200px;
}

.organic-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.organic-title {
  margin: 6px 0 4px;
  font-family: var(--font-display);
  font-size: 22px;
  line-height: 1.25;
  letter-spacing: -0.02em;
  color: var(--text-primary);
  transition: color 280ms ease;
}

.organic-card:hover .organic-title {
  background: var(--grad-primary);
  -webkit-background-clip: text;
          background-clip: text;
  color: transparent;
}

.organic-summary {
  margin: 0;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.65;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
}

.organic-foot {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 12px;
  padding-top: 10px;
}

.organic-link {
  color: var(--neon-cyan);
  letter-spacing: 0.05em;
}
</style>

<template>
  <RouterLink :to="`/blog/${article.slug}`" class="organic-card">
    <svg class="organic-border" viewBox="0 0 400 280" preserveAspectRatio="none" overflow="visible" aria-hidden="true">
      <path :d="path" class="organic-fill" />
      <path :d="path" class="organic-stroke" />
      <path :d="path" class="organic-flow" />
    </svg>

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
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import type { ArticleSummary } from '../api/article'

const props = defineProps<{ article: ArticleSummary }>()

const W = 400
const H = 280

function seedRandom(id: number, i: number) {
  const x = Math.sin(id * 9301 + i * 49297) * 10000
  return x - Math.floor(x)
}

function buildPath(id: number): string {
  const off = (i: number) => 6 + seedRandom(id, i) * 10
  const pts: Array<[number, number]> = [
    [off(1), 0],
    [W * 0.28, -off(2)],
    [W * 0.55, off(3) * 0.6],
    [W * 0.78, -off(4) * 0.8],
    [W - off(5), 0],
    [W + off(6) * 0.6, H * 0.32],
    [W - off(7) * 0.8, H * 0.58],
    [W + off(8) * 0.5, H * 0.82],
    [W - off(9), H],
    [W * 0.62, H + off(10) * 0.6],
    [W * 0.32, H - off(11) * 0.6],
    [off(12), H],
    [-off(13) * 0.6, H * 0.72],
    [off(14) * 0.8, H * 0.48],
    [-off(15) * 0.6, H * 0.22]
  ]
  const head = `M ${pts[0][0].toFixed(2)},${pts[0][1].toFixed(2)}`
  const tail = pts.slice(1).map((p) => `T ${p[0].toFixed(2)},${p[1].toFixed(2)}`).join(' ')
  return `${head} ${tail} Z`
}

const path = computed(() => buildPath(props.article.id))
</script>

<style scoped>
.organic-card {
  position: relative;
  display: block;
  padding: 28px 26px 22px;
  color: inherit;
  text-decoration: none;
  transition: transform 320ms cubic-bezier(0.2, 0.8, 0.2, 1);
}

.organic-card:hover {
  transform: translateY(-4px);
}

.organic-border {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: visible;
}

.organic-fill {
  fill: rgba(22, 14, 42, 0.72);
  stroke: none;
}

.organic-stroke {
  fill: none;
  stroke: rgba(196, 181, 253, 0.3);
  stroke-width: 1.1;
  transition: stroke 280ms ease, stroke-width 280ms ease;
}

.organic-flow {
  fill: none;
  stroke: var(--neon-cyan);
  stroke-width: 1.6;
  stroke-dasharray: 14 22;
  stroke-dashoffset: 0;
  opacity: 0;
  filter: drop-shadow(0 0 6px rgba(34, 211, 238, 0.55));
  transition: opacity 280ms ease;
}

.organic-card:hover .organic-stroke {
  stroke: rgba(167, 139, 250, 0.65);
  stroke-width: 1.3;
}

.organic-card:hover .organic-flow {
  opacity: 1;
  animation: organic-dash 2.6s linear infinite;
}

@keyframes organic-dash {
  to { stroke-dashoffset: -144; }
}

.organic-content {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 220px;
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

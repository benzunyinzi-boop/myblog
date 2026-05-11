<template>
  <div v-if="items.length" class="exp-timeline">
    <article
      v-for="(e, i) in items"
      :key="i"
      class="exp-item"
      :style="{ '--i': i }"
    >
      <span class="exp-node" />
      <div class="exp-card glass-card">
        <div class="exp-period">[ {{ e.period }} ]</div>
        <h4 class="exp-role">
          {{ e.role }}
          <span class="exp-org">@ {{ e.org }}</span>
        </h4>
        <div class="exp-stack">
          <span v-for="s in e.stack" :key="s" class="tech-pill">{{ s }}</span>
        </div>
        <ul class="exp-highlights">
          <li v-for="h in e.highlights" :key="h">{{ h }}</li>
        </ul>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import type { Experience } from '../data/experience'

defineProps<{ items: Experience[] }>()
</script>

<style scoped>
.exp-timeline {
  position: relative;
  display: grid;
  gap: 26px;
  padding-left: 32px;
}

.exp-timeline::before {
  content: '';
  position: absolute;
  left: 10px;
  top: 6px;
  bottom: 6px;
  width: 1px;
  background: linear-gradient(
    to bottom,
    rgba(139, 92, 246, 0) 0%,
    rgba(139, 92, 246, 0.55) 12%,
    rgba(34, 211, 238, 0.45) 88%,
    rgba(34, 211, 238, 0) 100%
  );
}

.exp-item {
  position: relative;
}

.exp-node {
  position: absolute;
  left: -27px;
  top: 22px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--grad-primary);
  box-shadow: 0 0 14px rgba(139, 92, 246, 0.65);
  animation: exp-pulse 3.2s ease-in-out infinite;
  animation-delay: calc(var(--i) * 320ms);
}

@keyframes exp-pulse {
  0%, 100% { transform: scale(1);   opacity: 1; }
  50%      { transform: scale(1.25); opacity: 0.75; }
}

.exp-card {
  padding: 22px 24px;
}

.exp-period {
  font-family: var(--font-mono);
  font-size: 12px;
  letter-spacing: 0.1em;
  color: var(--neon-cyan);
  margin-bottom: 6px;
}

.exp-role {
  margin: 0 0 10px;
  font-family: var(--font-display);
  font-size: 22px;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.exp-org {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 400;
  color: var(--text-secondary);
  margin-left: 4px;
}

.exp-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.exp-highlights {
  margin: 0;
  padding-left: 18px;
  color: var(--text-secondary);
  line-height: 1.85;
  font-size: 14px;
}

.exp-highlights li {
  margin-bottom: 4px;
}

.exp-highlights li::marker {
  color: var(--primary-2);
}
</style>

<template>
  <Transition name="nav-preview">
    <div v-if="items.length" class="nav-preview" :style="{ left: `${anchorX}px` }">
      <div class="nav-preview-inner">
        <span
          v-for="(item, i) in items"
          :key="item.label"
          class="nav-preview-item"
          :style="{ '--i': i }"
        >
          <span class="nav-preview-icon">{{ item.icon }}</span>
          <span class="nav-preview-label">{{ item.label }}</span>
        </span>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type PreviewItem = { icon: string; label: string }

const PREVIEWS: Record<string, PreviewItem[]> = {
  '/tech': [
    { icon: '🐹', label: 'Go' },
    { icon: '🐬', label: 'MySQL' },
    { icon: '🟥', label: 'Redis' },
    { icon: '🪶', label: 'Kafka' },
    { icon: '✨', label: 'AI' },
    { icon: '☁️', label: 'Cloud' }
  ],
  '/about': [
    { icon: '👨‍💻', label: 'Engineer' },
    { icon: '📚', label: 'Writer' },
    { icon: '📷', label: 'Camera' },
    { icon: '🎧', label: 'Music' }
  ],
  '/photo': [
    { icon: '📷', label: 'Street' },
    { icon: '🌆', label: 'City' },
    { icon: '🏞️', label: 'Landscape' },
    { icon: '🎞️', label: 'Film' }
  ]
}

const props = defineProps<{
  path: string | null
  anchorX: number
}>()

const items = computed<PreviewItem[]>(() =>
  props.path && PREVIEWS[props.path] ? PREVIEWS[props.path] : []
)
</script>

<style scoped>
.nav-preview {
  position: absolute;
  top: calc(100% + 14px);
  transform: translateX(-50%);
  z-index: 60;
  pointer-events: none;
}

.nav-preview-inner {
  display: inline-flex;
  align-items: center;
  gap: 16px;
  padding: 12px 18px;
  background: linear-gradient(
    160deg,
    rgba(30, 20, 60, 0.78) 0%,
    rgba(18, 11, 34, 0.88) 100%
  );
  backdrop-filter: blur(18px) saturate(160%);
  -webkit-backdrop-filter: blur(18px) saturate(160%);
  border: 1px solid rgba(196, 181, 253, 0.22);
  border-radius: 18px;
  box-shadow:
    0 1px 0 rgba(255, 255, 255, 0.06) inset,
    0 24px 48px -18px rgba(0, 0, 0, 0.6);
  position: relative;
}

.nav-preview-inner::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: inherit;
  padding: 1px;
  background: linear-gradient(
    135deg,
    rgba(139, 92, 246, 0),
    rgba(139, 92, 246, 0.55),
    rgba(34, 211, 238, 0.4),
    rgba(139, 92, 246, 0)
  );
  -webkit-mask:
    linear-gradient(#000, #000) content-box,
    linear-gradient(#000, #000);
  -webkit-mask-composite: xor;
          mask-composite: exclude;
  pointer-events: none;
  opacity: 0.5;
  animation: nav-preview-breathe 4s ease-in-out infinite;
}

.nav-preview-inner::after {
  content: '';
  position: absolute;
  top: -7px;
  left: 50%;
  width: 12px;
  height: 12px;
  transform: translateX(-50%) rotate(45deg);
  background: inherit;
  border-top: 1px solid rgba(196, 181, 253, 0.22);
  border-left: 1px solid rgba(196, 181, 253, 0.22);
}

.nav-preview-item {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  min-width: 44px;
  opacity: 0;
  transform: translateY(8px);
  animation: nav-item-in 360ms ease forwards, nav-item-bob 3s ease-in-out infinite;
  animation-delay:
    calc(var(--i) * 80ms),
    calc(var(--i) * 80ms + 360ms);
}

.nav-preview-icon {
  font-size: 24px;
  line-height: 1;
  filter: drop-shadow(0 0 6px rgba(139, 92, 246, 0.45));
}

.nav-preview-label {
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-secondary);
}

@keyframes nav-item-in {
  to { opacity: 1; transform: translateY(0); }
}

@keyframes nav-item-bob {
  0%, 100% { transform: translateY(0); }
  50%      { transform: translateY(-4px); }
}

@keyframes nav-preview-breathe {
  0%, 100% { opacity: 0.3; }
  50%      { opacity: 0.7; }
}

.nav-preview-enter-active,
.nav-preview-leave-active {
  transition: opacity 220ms ease, transform 220ms ease;
}
.nav-preview-enter-from,
.nav-preview-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-6px);
}
.nav-preview-enter-to,
.nav-preview-leave-from {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

@media (max-width: 960px), (hover: none), (pointer: coarse) {
  .nav-preview { display: none; }
}
</style>

<template>
  <div class="identity-typer" :aria-label="current">
    <span class="identity-prompt">&gt;</span>
    <span class="identity-current">{{ displayed }}</span>
    <span class="identity-caret" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
  items: string[]
  typeSpeed?: number
  holdMs?: number
  eraseSpeed?: number
}>()

const typeSpeed = props.typeSpeed ?? 90
const holdMs = props.holdMs ?? 1600
const eraseSpeed = props.eraseSpeed ?? 45

const displayed = ref('')
const current = ref('')
let idx = 0
let timer: number | null = null

function stop() {
  if (timer !== null) {
    window.clearTimeout(timer)
    timer = null
  }
}

function typeOne(target: string, i: number) {
  current.value = target
  displayed.value = target.slice(0, i)
  if (i < target.length) {
    timer = window.setTimeout(() => typeOne(target, i + 1), typeSpeed)
  } else {
    timer = window.setTimeout(() => eraseOne(target.length), holdMs)
  }
}

function eraseOne(i: number) {
  displayed.value = displayed.value.slice(0, i)
  if (i > 0) {
    timer = window.setTimeout(() => eraseOne(i - 1), eraseSpeed)
  } else {
    idx = (idx + 1) % props.items.length
    timer = window.setTimeout(() => typeOne(props.items[idx], 0), 320)
  }
}

onMounted(() => {
  if (props.items.length) {
    typeOne(props.items[0], 0)
  }
})

onUnmounted(stop)
</script>

<style scoped>
.identity-typer {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  font-family: var(--font-mono);
  font-size: 18px;
  letter-spacing: 0.02em;
  color: var(--text-primary);
}

.identity-prompt {
  color: var(--neon-cyan);
}

.identity-current {
  background: var(--grad-primary);
  -webkit-background-clip: text;
          background-clip: text;
  color: transparent;
  font-weight: 600;
  min-height: 1em;
}

.identity-caret {
  display: inline-block;
  width: 7px;
  height: 1em;
  background: #c4b5fd;
  vertical-align: text-bottom;
  animation: identity-caret-blink 1s steps(1) infinite;
  box-shadow: 0 0 8px rgba(196, 181, 253, 0.8);
}

@keyframes identity-caret-blink {
  0%, 50%      { opacity: 1; }
  50.01%, 100% { opacity: 0; }
}
</style>

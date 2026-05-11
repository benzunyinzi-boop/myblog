<template>
  <div class="term">
    <!-- 窗口头部:红黄绿灯 + 标题 -->
    <header class="term-head">
      <div class="term-dots">
        <span class="dot dot-red" />
        <span class="dot dot-yellow" />
        <span class="dot dot-green" />
      </div>
      <div class="term-title">myblog@core &mdash; zsh &mdash; 80×24</div>
      <div class="term-clock">{{ clock }}</div>
    </header>

    <!-- slogan:启动 / 空闲横幅,保留品牌认同 -->
    <div class="term-slogan">
      <span class="slogan-prompt">:: boot ::</span>
      把分布式系统的锋利,写成<span class="slogan-hl">可读的经验</span>。
    </div>

    <!-- 命令行输出区 -->
    <div class="term-body">
      <template v-if="state === 'loading'">
        <div class="term-line muted">[ OK ] Mounting filesystems...</div>
        <div class="term-line muted">[ OK ] Starting network service...</div>
        <div class="term-line muted">[ OK ] Reached target graphical interface.</div>
        <div class="term-line ok">welcome back, yinyin ✨</div>
      </template>

      <template v-else>
        <!-- 输入行:提示符 + 打字 + 闪烁光标 -->
        <div class="term-line">
          <span class="prompt">yinyin@blog</span>
          <span class="prompt-sep">:</span>
          <span class="path">~</span>
          <span class="prompt-sep">$</span>
          <span class="cmd">{{ typed }}</span>
          <span v-if="state === 'typing'" class="caret" />
        </div>

        <!-- 输出行:模拟 systemctl status 的可读格式 -->
        <template v-if="state === 'output' || state === 'idle'">
          <div class="term-line">
            <span class="bullet ok">●</span>
            <span class="unit">core.service</span>
            &mdash; myblog backend runtime
          </div>
          <div class="term-line muted indent">Loaded: loaded (/etc/systemd/system/core.service; enabled)</div>
          <div class="term-line indent">
            <span>Active: </span><span class="ok">active (running)</span>
            <span class="muted"> since {{ bootTime }}; {{ uptime }}</span>
          </div>
          <div class="term-line muted indent">Tasks: 42 (limit: infinity)</div>
          <div class="term-line muted indent">Memory: 128.5M&nbsp;&nbsp;CPU: 3.2%</div>
          <div class="term-line muted indent">
            CGroup: /backend.slice/core.service
          </div>
          <div class="term-line muted indent-2">
            └─ myblog-server --release --trace=on
          </div>
          <div v-if="state === 'idle'" class="term-line">
            <span class="prompt">yinyin@blog</span>
            <span class="prompt-sep">:</span>
            <span class="path">~</span>
            <span class="prompt-sep">$</span>
            <span class="caret" />
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const COMMAND = 'systemctl status core'
const PROMPT_STATES = ['loading', 'typing', 'output', 'idle'] as const
type TermState = typeof PROMPT_STATES[number]

const state = ref<TermState>('loading')
const typed = ref('')
const clock = ref('')
const uptime = ref('0s')
const bootTime = ref('')

let typingTimer: number | null = null
let sequenceTimer: number | null = null
let clockTimer: number | null = null
let bootedAt = 0

function tickClock() {
  const d = new Date()
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  const ss = String(d.getSeconds()).padStart(2, '0')
  clock.value = `${hh}:${mm}:${ss}`

  if (bootedAt > 0) {
    const delta = Math.floor((Date.now() - bootedAt) / 1000)
    if (delta < 60) uptime.value = `${delta}s ago`
    else if (delta < 3600) uptime.value = `${Math.floor(delta / 60)}m ${delta % 60}s ago`
    else uptime.value = `${Math.floor(delta / 3600)}h ${Math.floor((delta % 3600) / 60)}m ago`
  }
}

function typeCommand() {
  typed.value = ''
  state.value = 'typing'
  let i = 0
  typingTimer = window.setInterval(() => {
    if (i >= COMMAND.length) {
      if (typingTimer) window.clearInterval(typingTimer)
      typingTimer = null
      sequenceTimer = window.setTimeout(() => {
        state.value = 'output'
        sequenceTimer = window.setTimeout(() => {
          state.value = 'idle'
          sequenceTimer = window.setTimeout(runLoop, 6000)
        }, 2500)
      }, 400)
      return
    }
    typed.value += COMMAND[i++]
  }, 90)
}

function runLoop() {
  typeCommand()
}

function startBootSequence() {
  state.value = 'loading'
  sequenceTimer = window.setTimeout(() => {
    bootedAt = Date.now()
    const d = new Date(bootedAt)
    bootTime.value = d.toLocaleString('zh-CN', { hour12: false })
    runLoop()
  }, 1600)
}

onMounted(() => {
  tickClock()
  clockTimer = window.setInterval(tickClock, 1000)
  startBootSequence()
})

onUnmounted(() => {
  if (typingTimer) window.clearInterval(typingTimer)
  if (sequenceTimer) window.clearTimeout(sequenceTimer)
  if (clockTimer) window.clearInterval(clockTimer)
})
</script>

<style scoped>
.term {
  position: relative;
  width: 100%;
  max-width: 560px;
  border-radius: 16px;
  padding: 0;
  overflow: hidden;
  /* 毛玻璃主体 */
  background: linear-gradient(
    160deg,
    rgba(30, 20, 60, 0.55) 0%,
    rgba(18, 12, 38, 0.65) 100%
  );
  backdrop-filter: blur(18px) saturate(160%);
  -webkit-backdrop-filter: blur(18px) saturate(160%);
  border: 1px solid rgba(196, 181, 253, 0.18);
  box-shadow:
    0 1px 0 rgba(255, 255, 255, 0.05) inset,
    0 24px 60px -16px rgba(0, 0, 0, 0.55);
}

/* 紫色呼吸灯:在卡片外圈脉动 */
.term::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: inherit;
  padding: 1px;
  background: linear-gradient(
    135deg,
    rgba(139, 92, 246, 0.0),
    rgba(139, 92, 246, 0.6),
    rgba(34, 211, 238, 0.4),
    rgba(139, 92, 246, 0.0)
  );
  -webkit-mask:
    linear-gradient(#000, #000) content-box,
    linear-gradient(#000, #000);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  pointer-events: none;
  opacity: 0.45;
  animation: term-breathe 4.5s ease-in-out infinite;
}

.term-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(14, 8, 28, 0.55);
  border-bottom: 1px solid rgba(196, 181, 253, 0.12);
}

.term-dots { display: flex; gap: 6px; }
.dot { width: 11px; height: 11px; border-radius: 50%; display: inline-block; }
.dot-red    { background: #ff5f57; }
.dot-yellow { background: #febc2e; }
.dot-green  { background: #28c840; }

.term-title {
  flex: 1;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 12px;
  color: rgba(196, 181, 253, 0.55);
  letter-spacing: 0.02em;
}

.term-clock {
  font-family: var(--font-mono);
  font-size: 12px;
  color: rgba(196, 181, 253, 0.7);
  letter-spacing: 0.04em;
  font-variant-numeric: tabular-nums;
}

.term-slogan {
  padding: 12px 18px 10px;
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: rgba(220, 210, 250, 0.75);
  letter-spacing: 0.01em;
  border-bottom: 1px dashed rgba(196, 181, 253, 0.14);
}

.slogan-prompt {
  display: inline-block;
  margin-right: 8px;
  color: rgba(34, 211, 238, 0.75);
}

.slogan-hl {
  background: linear-gradient(90deg, #c4b5fd, #22d3ee);
  -webkit-background-clip: text;
          background-clip: text;
  color: transparent;
  font-weight: 600;
}

.term-body {
  padding: 16px 18px 20px;
  font-family: var(--font-mono);
  font-size: 13.5px;
  line-height: 1.75;
  color: rgba(234, 227, 255, 0.9);
  min-height: 240px;
}

.term-line { white-space: pre-wrap; word-break: break-all; }
.term-line.muted { color: rgba(180, 170, 210, 0.55); }
.term-line.ok { color: #6ee7b7; }
.term-line.indent { padding-left: 16px; }
.term-line.indent-2 { padding-left: 32px; }

.prompt { color: #8b5cf6; font-weight: 600; }
.prompt-sep { color: rgba(180, 170, 210, 0.5); margin: 0 2px; }
.path { color: #22d3ee; }
.cmd { margin-left: 8px; color: #ede9fe; }
.unit { color: #ede9fe; font-weight: 600; margin-right: 4px; }
.bullet.ok { color: #34d399; margin-right: 8px; }

.caret {
  display: inline-block;
  width: 7px;
  height: 1em;
  margin-left: 4px;
  background: #c4b5fd;
  vertical-align: text-bottom;
  animation: caret-blink 1s steps(1) infinite;
  box-shadow: 0 0 8px rgba(196, 181, 253, 0.8);
}

@keyframes caret-blink {
  0%, 50%   { opacity: 1; }
  50.01%, 100% { opacity: 0; }
}

@keyframes term-breathe {
  0%, 100% { opacity: 0.25; }
  50%      { opacity: 0.7; }
}

@media (max-width: 720px) {
  .term { max-width: 100%; }
  .term-body { font-size: 12.5px; }
}
</style>

<template>
  <div class="public-shell">
    <header class="site-header">
      <div class="container nav-inner">
        <RouterLink class="brand" to="/">
          <div class="brand-logo">
            <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <defs>
                <linearGradient id="logo-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#c4b5fd;stop-opacity:1" />
                  <stop offset="50%" style="stop-color:#8b5cf6;stop-opacity:1" />
                  <stop offset="100%" style="stop-color:#22d3ee;stop-opacity:1" />
                </linearGradient>
              </defs>
              <polygon points="20,2 38,12 38,28 20,38 2,28 2,12"
                       fill="url(#logo-gradient)"
                       stroke="url(#logo-gradient)"
                       stroke-width="1.5" />
              <text x="20" y="26"
                    font-family="'Courier New', monospace"
                    font-size="16"
                    font-weight="bold"
                    fill="#0a0118"
                    text-anchor="middle">Y</text>
            </svg>
          </div>
        </RouterLink>

        <nav ref="navEl" class="nav-links" @mouseleave="clearHover">
          <RouterLink
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            :class="linkClass(link.path)"
            :data-path="link.path"
            @mouseenter="onEnter($event, link.path)"
          >{{ link.label }}</RouterLink>

          <NavHoverPreview :path="hoveredPath" :anchor-x="anchorX" />
        </nav>

        <RouterLink class="admin-entry" to="/admin/login">Admin</RouterLink>
      </div>
    </header>

    <main>
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import NavHoverPreview from '../components/NavHoverPreview.vue'

const route = useRoute()

const navLinks = [
  { path: '/', label: '首页' },
  { path: '/tech', label: '技术' },
  { path: '/about', label: '关于' },
  { path: '/photo', label: '摄影' }
]

const navEl = ref<HTMLElement | null>(null)
const hoveredPath = ref<string | null>(null)
const anchorX = ref(0)

function linkClass(path: string) {
  return route.path === path ? 'nav-link active' : 'nav-link'
}

function onEnter(e: MouseEvent, path: string) {
  if (path === '/') {
    hoveredPath.value = null
    return
  }
  const target = e.currentTarget as HTMLElement
  const nav = navEl.value
  if (!target || !nav) return
  const tRect = target.getBoundingClientRect()
  const nRect = nav.getBoundingClientRect()
  anchorX.value = tRect.left - nRect.left + tRect.width / 2
  hoveredPath.value = path
}

function clearHover() {
  hoveredPath.value = null
}
</script>

<style scoped>
.nav-links {
  position: relative;
}
</style>

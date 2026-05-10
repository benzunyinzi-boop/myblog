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

        <nav class="nav-links">
          <RouterLink to="/" :class="linkClass('/')">首页</RouterLink>
          <RouterLink to="/tech" :class="linkClass('/tech')">技术</RouterLink>
          <RouterLink to="/about" :class="linkClass('/about')">关于</RouterLink>
          <a class="ghost-link" href="#photography">摄影</a>
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
import { RouterLink, RouterView, useRoute } from 'vue-router'

const route = useRoute()

function linkClass(path: string) {
  return route.path === path ? 'nav-link active' : 'nav-link'
}
</script>

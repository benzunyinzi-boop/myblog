<template>
  <div class="public-shell">
    <header class="site-header">
      <div class="container nav-inner">
        <RouterLink class="brand" to="/">
          <span class="brand-logo" />
          <span class="brand-text">yinyin.dev</span>
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

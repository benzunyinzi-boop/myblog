<template>
  <section class="admin-shell">
    <aside class="admin-sidebar glass-card">
      <div>
        <span class="section-chip">admin / panel</span>
        <h1 class="sidebar-title">内容控制台</h1>
        <p class="page-copy">后台 CRUD 面板。内容围绕 M3 后端 API 展开。</p>
      </div>

      <nav class="admin-nav">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="admin-nav-item"
          :class="{ active: isActive(item.path) }"
        >
          <span class="nav-label">{{ item.label }}</span>
          <span class="nav-meta">{{ item.meta }}</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="admin-main">
      <header class="admin-topbar glass-card">
        <div>
          <span class="section-chip subtle">session / active</span>
          <h2>欢迎回来, {{ auth.profile?.username || 'admin' }}</h2>
        </div>
        <n-button quaternary @click="handleLogout">退出登录</n-button>
      </header>

      <RouterView />
    </main>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { NButton } from 'naive-ui'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const navItems = [
  { path: '/admin', label: 'Dashboard', meta: '概览' },
  { path: '/admin/articles', label: 'Articles', meta: '文章' },
  { path: '/admin/categories', label: 'Categories', meta: '分类' },
  { path: '/admin/tags', label: 'Tags', meta: '标签' },
  { path: '/admin/profile', label: 'Profile', meta: '关于我' }
]

const currentPath = computed(() => route.path)

function isActive(path: string) {
  if (path === '/admin') return currentPath.value === '/admin'
  return currentPath.value.startsWith(path)
}

function handleLogout() {
  auth.signOut()
  router.push('/admin/login')
}
</script>

<template>
  <section class="admin-shell">
    <aside class="admin-sidebar glass-card">
      <div>
        <span class="section-chip">admin / panel</span>
        <h1 class="sidebar-title">内容控制台</h1>
        <p class="page-copy">这是后台第一版入口,下一步会在这里接文章/分类/标签/关于我/上传管理。</p>
      </div>

      <nav class="admin-nav">
        <span class="admin-nav-item active">Dashboard</span>
        <span class="admin-nav-item">Articles</span>
        <span class="admin-nav-item">Categories</span>
        <span class="admin-nav-item">Tags</span>
        <span class="admin-nav-item">Profile</span>
        <span class="admin-nav-item">Uploads</span>
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

      <section class="admin-cards">
        <article class="glass-card admin-card">
          <span class="metric-label">Mode</span>
          <strong>Editor Ready</strong>
          <p>认证已经打通,下一步可以在这里接真正的文章管理列表。</p>
        </article>
        <article class="glass-card admin-card">
          <span class="metric-label">API</span>
          <strong>M3 connected</strong>
          <p>后端 Category / Tag / Article / Profile / Upload 接口都已经就绪。</p>
        </article>
        <article class="glass-card admin-card">
          <span class="metric-label">Theme</span>
          <strong>Purple tech</strong>
          <p>后台也沿用了同一套霓虹科技感视觉变量,后续只需扩展组件。</p>
        </article>
      </section>
    </main>
  </section>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NButton } from 'naive-ui'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const router = useRouter()

function handleLogout() {
  auth.signOut()
  router.push('/admin/login')
}
</script>

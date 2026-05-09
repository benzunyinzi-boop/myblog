<template>
  <section class="section-gap container about-page">
    <div class="profile-orbit" />
    <div class="about-grid">
      <div class="glass-card profile-card">
        <div class="avatar-shell">
          <div class="avatar-core">{{ initials }}</div>
        </div>
        <div class="profile-meta">
          <span class="section-chip">profile / public</span>
          <h1 class="page-title">{{ profile.name || 'Yinyin' }}</h1>
          <p class="page-copy">{{ profile.bio || '十年 Go 后端工程经验,正在把工程实践变成可读的内容。' }}</p>
        </div>
      </div>

      <div class="glass-card contact-card">
        <h2>数字名片</h2>
        <ul class="contact-list">
          <li><strong>Email</strong><span>{{ profile.email || 'me@example.com' }}</span></li>
          <li><strong>GitHub</strong><span>{{ profile.github || 'https://github.com/' }}</span></li>
          <li><strong>LinkedIn</strong><span>{{ profile.linkedin || 'https://linkedin.com/' }}</span></li>
          <li><strong>Website</strong><span>{{ profile.website || 'https://yinyin.dev' }}</span></li>
        </ul>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAsyncState } from '@vueuse/core'
import { fetchProfile } from '../../api/profile'

const { state } = useAsyncState(() => fetchProfile(), null)
const profile = computed(() => state.value?.data ?? {
  name: '',
  bio: '',
  avatar: '',
  email: '',
  github: '',
  twitter: '',
  linkedin: '',
  website: ''
})

const initials = computed(() => {
  const name = profile.value.name || 'YY'
  return name.slice(0, 2).toUpperCase()
})
</script>

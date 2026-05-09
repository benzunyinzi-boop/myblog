<template>
  <section class="admin-page">
    <header class="admin-page-head">
      <div>
        <span class="section-chip">profile / edit</span>
        <h2>关于我</h2>
        <p class="page-copy">公开站点的个人资料,会展示在 /about 和 /public/profile 接口里。</p>
      </div>
      <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
    </header>

    <section class="glass-card admin-form-card">
      <n-form :model="form" label-placement="top">
        <div class="form-grid">
          <n-form-item label="姓名 / 昵称" required>
            <n-input v-model:value="form.name" placeholder="Yinyin" />
          </n-form-item>
          <n-form-item label="Email" required>
            <n-input v-model:value="form.email" placeholder="me@example.com" />
          </n-form-item>
        </div>

        <n-form-item label="简介" required>
          <n-input
            v-model:value="form.bio"
            type="textarea"
            :rows="4"
            placeholder="用一段话介绍自己"
          />
        </n-form-item>

        <div class="form-grid">
          <n-form-item label="头像 URL">
            <n-input v-model:value="form.avatar" placeholder="/uploads/... 或外链" />
          </n-form-item>
          <n-form-item label="个人网站">
            <n-input v-model:value="form.website" placeholder="https://yinyin.dev" />
          </n-form-item>
        </div>

        <div class="form-grid">
          <n-form-item label="GitHub">
            <n-input v-model:value="form.github" placeholder="https://github.com/yourname" />
          </n-form-item>
          <n-form-item label="Twitter">
            <n-input v-model:value="form.twitter" placeholder="https://twitter.com/yourname" />
          </n-form-item>
        </div>

        <n-form-item label="LinkedIn">
          <n-input v-model:value="form.linkedin" placeholder="https://linkedin.com/in/yourname" />
        </n-form-item>
      </n-form>
    </section>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { NButton, NForm, NFormItem, NInput } from 'naive-ui'
import { fetchProfile } from '../../api/profile'
import { adminUpdateProfile } from '../../api/admin'
import { discrete } from '../../main'

const saving = ref(false)

const form = reactive({
  name: '',
  bio: '',
  avatar: '',
  email: '',
  github: '',
  twitter: '',
  linkedin: '',
  website: ''
})

async function load() {
  try {
    const resp = await fetchProfile()
    Object.assign(form, resp.data)
  } catch (err) {
    discrete.message.error('加载个人资料失败')
    console.error(err)
  }
}

async function handleSave() {
  if (!form.name.trim() || !form.email.trim() || !form.bio.trim()) {
    discrete.message.warning('姓名、Email、简介不能为空')
    return
  }
  saving.value = true
  try {
    const resp = await adminUpdateProfile({
      name: form.name.trim(),
      bio: form.bio.trim(),
      avatar: form.avatar.trim(),
      email: form.email.trim(),
      github: form.github.trim(),
      twitter: form.twitter.trim(),
      linkedin: form.linkedin.trim(),
      website: form.website.trim()
    })
    if (resp.code !== 0) {
      discrete.message.error(resp.message || '保存失败')
      return
    }
    discrete.message.success('已保存')
  } catch (err: any) {
    const msg = err?.response?.data?.message || '保存失败'
    discrete.message.error(msg)
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

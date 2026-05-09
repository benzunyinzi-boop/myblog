<template>
  <section class="admin-login-shell">
    <div class="admin-grid">
      <div class="login-brand glass-card">
        <span class="section-chip">admin / gate</span>
        <h1 class="page-title">控制台登录</h1>
        <p class="page-copy">
          这是内容后台的入口。登录后将进入文章、分类、标签、关于我与图片上传的管理界面。
        </p>
      </div>

      <n-card class="login-card" :bordered="false">
        <n-form :model="form" @submit.prevent="handleSubmit">
          <n-form-item label="用户名">
            <n-input v-model:value="form.username" placeholder="admin" />
          </n-form-item>
          <n-form-item label="密码">
            <n-input v-model:value="form.password" type="password" show-password-on="click" placeholder="Admin@123" />
          </n-form-item>
          <n-button type="primary" block :loading="auth.pending" attr-type="submit">
            进入控制台
          </n-button>
        </n-form>
      </n-card>
    </div>
  </section>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NCard, NForm, NFormItem, NInput } from 'naive-ui'
import { useAuthStore } from '../../stores/auth'
import { discrete } from '../../main'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const form = reactive({
  username: 'admin',
  password: 'Admin@123'
})

async function handleSubmit() {
  try {
    const envelope = await auth.signIn(form)
    if (envelope.code === 0) {
      discrete.message.success('登录成功')
      const redirect = (route.query.redirect as string) || '/admin'
      router.push(redirect)
      return
    }
    discrete.message.error(envelope.message)
  } catch (error: any) {
    const apiMessage = error?.response?.data?.message
    discrete.message.error(apiMessage || '登录失败,请检查后端服务是否已启动')
    console.error(error)
  }
}
</script>

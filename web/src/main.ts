import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createDiscreteApi } from 'naive-ui'
import App from './App.vue'
import { router } from './router'
import { setUnauthorizedHandler } from './api/http'
import { useAuthStore } from './stores/auth'
import './styles/theme.css'
import './styles/global.css'
import './styles/markdown.css'
import 'highlight.js/styles/atom-one-dark.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 现在 pinia 已经 installed,可以安全使用 useAuthStore
const auth = useAuthStore()
export const discrete = createDiscreteApi(['message', 'notification', 'dialog'])

setUnauthorizedHandler(() => {
  // 已经在登录页就不需要再跳了,避免无限循环
  if (!auth.isLoggedIn && router.currentRoute.value.path === '/admin/login') {
    return
  }
  auth.signOut()
  discrete.message.warning('登录已过期,请重新登录')
  router.push({ path: '/admin/login', query: { redirect: router.currentRoute.value.fullPath } })
})

app.mount('#app')

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createDiscreteApi } from 'naive-ui'
import App from './App.vue'
import { router } from './router'
import './styles/theme.css'
import './styles/global.css'
import './styles/markdown.css'
import 'highlight.js/styles/atom-one-dark.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.mount('#app')

export const discrete = createDiscreteApi(['message', 'notification', 'dialog'])

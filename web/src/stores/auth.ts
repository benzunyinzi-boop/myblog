import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { login, type LoginReq } from '../api/auth'

const TOKEN_KEY = 'myblog-admin-token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const profile = ref<{ username: string; role: string } | null>(null)
  const pending = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function signIn(payload: LoginReq) {
    pending.value = true
    try {
      const envelope = await login(payload)
      token.value = envelope.data.access_token
      profile.value = {
        username: envelope.data.user.username,
        role: envelope.data.user.role
      }
      localStorage.setItem(TOKEN_KEY, token.value)
      return envelope
    } finally {
      pending.value = false
    }
  }

  function signOut() {
    token.value = ''
    profile.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  return {
    token,
    profile,
    pending,
    isLoggedIn,
    signIn,
    signOut
  }
})

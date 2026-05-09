import axios, { AxiosError } from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:18080/api/v1'

export const TOKEN_KEY = 'myblog-admin-token'

export const http = axios.create({
  baseURL,
  timeout: 10000
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

/**
 * 业务层需要的 token 失效回调。
 * 在 main.ts 里注入 store/router 的清理和跳转逻辑,
 * 避免这里直接依赖 pinia/router 造成循环引用。
 */
let unauthorizedHandler: (() => void) | null = null

export function setUnauthorizedHandler(handler: () => void) {
  unauthorizedHandler = handler
}

http.interceptors.response.use(
  (response) => {
    const envelope = response.data
    if (envelope && typeof envelope === 'object' && 'code' in envelope) {
      const code = envelope.code
      // 10002 / 30004 / 30005 都视作登录态问题
      if (code === 10002 || code === 30004 || code === 30005) {
        unauthorizedHandler?.()
      }
    }
    return response
  },
  (error: AxiosError) => {
    const status = error.response?.status
    if (status === 401) {
      unauthorizedHandler?.()
    }
    return Promise.reject(error)
  }
)

export type ApiEnvelope<T> = {
  code: number
  message: string
  data: T
}

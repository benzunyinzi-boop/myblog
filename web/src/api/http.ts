import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:18080/api/v1'

export const http = axios.create({
  baseURL,
  timeout: 10000
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('myblog-admin-token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error)
)

export type ApiEnvelope<T> = {
  code: number
  message: string
  data: T
}

import { http, type ApiEnvelope } from './http'

export type LoginReq = {
  username: string
  password: string
}

export type LoginResp = {
  access_token: string
  refresh_token: string
  token_type: string
  expires_at: number
  user: {
    id: number
    username: string
    nickname: string
    avatar: string
    role: string
  }
}

export async function login(payload: LoginReq) {
  const { data } = await http.post<ApiEnvelope<LoginResp>>('/admin/auth/login', payload)
  return data
}

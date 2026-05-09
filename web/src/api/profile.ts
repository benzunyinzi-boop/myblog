import { http, type ApiEnvelope } from './http'

export type Profile = {
  name: string
  bio: string
  avatar: string
  email: string
  github: string
  twitter: string
  linkedin: string
  website: string
}

export async function fetchProfile() {
  const { data } = await http.get<ApiEnvelope<Profile>>('/public/profile')
  return data
}

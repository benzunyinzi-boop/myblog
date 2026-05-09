import { http, type ApiEnvelope } from './http'

export type Category = {
  id: number
  name: string
  slug: string
  description: string
  sort_order: number
}

export type Tag = {
  id: number
  name: string
  slug: string
}

export async function fetchCategories() {
  const { data } = await http.get<ApiEnvelope<{ items: Category[]; total: number }>>('/public/categories')
  return data
}

export async function fetchTags() {
  const { data } = await http.get<ApiEnvelope<{ items: Tag[]; total: number }>>('/public/tags')
  return data
}

import { http, type ApiEnvelope } from './http'
import type { Category, Tag } from './meta'
import type { ArticleDetail, ArticleListResp, ArticleSummary } from './article'

// ---------- Category ----------

export type CategoryPayload = {
  name: string
  slug: string
  description?: string
  sort_order?: number
}

export async function adminCreateCategory(payload: CategoryPayload) {
  const { data } = await http.post<ApiEnvelope<Category>>('/admin/categories', payload)
  return data
}

export async function adminUpdateCategory(id: number, payload: CategoryPayload) {
  const { data } = await http.put<ApiEnvelope<Category>>(`/admin/categories/${id}`, payload)
  return data
}

export async function adminDeleteCategory(id: number) {
  const { data } = await http.delete<ApiEnvelope<null>>(`/admin/categories/${id}`)
  return data
}

// ---------- Tag ----------

export type TagPayload = {
  name: string
  slug: string
}

export async function adminCreateTag(payload: TagPayload) {
  const { data } = await http.post<ApiEnvelope<Tag>>('/admin/tags', payload)
  return data
}

export async function adminDeleteTag(id: number) {
  const { data } = await http.delete<ApiEnvelope<null>>(`/admin/tags/${id}`)
  return data
}

// ---------- Article ----------

export type ArticlePayload = {
  title: string
  slug: string
  summary?: string
  content: string
  cover_image?: string
  category_id?: number
  tag_ids?: number[]
  status?: 'draft' | 'published'
}

export type AdminArticleListParams = {
  page?: number
  page_size?: number
  status?: 'draft' | 'published'
  category_id?: number
  tag_id?: number
  keyword?: string
}

export async function adminListArticles(params: AdminArticleListParams = {}) {
  const { data } = await http.get<ApiEnvelope<ArticleListResp>>('/admin/articles', { params })
  return data
}

export async function adminGetArticle(id: number) {
  const { data } = await http.get<ApiEnvelope<ArticleDetail>>(`/admin/articles/${id}`)
  return data
}

export async function adminCreateArticle(payload: ArticlePayload) {
  const { data } = await http.post<ApiEnvelope<ArticleDetail>>('/admin/articles', payload)
  return data
}

export async function adminUpdateArticle(id: number, payload: ArticlePayload) {
  const { data } = await http.put<ApiEnvelope<ArticleDetail>>(`/admin/articles/${id}`, payload)
  return data
}

export async function adminDeleteArticle(id: number) {
  const { data } = await http.delete<ApiEnvelope<null>>(`/admin/articles/${id}`)
  return data
}

export async function adminPublishArticle(id: number) {
  const { data } = await http.post<ApiEnvelope<ArticleDetail>>(`/admin/articles/${id}/publish`)
  return data
}

export async function adminUnpublishArticle(id: number) {
  const { data } = await http.post<ApiEnvelope<ArticleDetail>>(`/admin/articles/${id}/unpublish`)
  return data
}

// 类型透传,避免外层重复 import
export type { Category, Tag, ArticleSummary, ArticleDetail, ArticleListResp }

// ---------- Profile ----------

export type ProfilePayload = {
  name: string
  bio: string
  avatar?: string
  email: string
  github?: string
  twitter?: string
  linkedin?: string
  website?: string
}

export async function adminUpdateProfile(payload: ProfilePayload) {
  const { data } = await http.put<ApiEnvelope<ProfilePayload>>('/admin/profile', payload)
  return data
}

// ---------- Upload ----------

export async function adminUploadFile(file: File) {
  const form = new FormData()
  form.append('file', file)
  const { data } = await http.post<ApiEnvelope<{ url: string }>>('/admin/uploads', form, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return data
}

import { http, type ApiEnvelope } from './http'

export type ArticleSummary = {
  id: number
  title: string
  slug: string
  summary: string
  cover_image: string
  category_id: number
  author_id: number
  status: string
  view_count: number
  published_at?: number
  created_at: number
  tags: Array<{ id: number; name: string; slug: string }>
}

export type ArticleDetail = ArticleSummary & {
  content: string
}

export type ArticleListResp = {
  items: ArticleSummary[]
  total: number
  page: number
  page_size: number
}

export async function fetchPublicArticles(params?: Record<string, unknown>) {
  const { data } = await http.get<ApiEnvelope<ArticleListResp>>('/public/articles', { params })
  return data
}

export async function fetchPublicArticleBySlug(slug: string) {
  const { data } = await http.get<ApiEnvelope<ArticleDetail>>(`/public/articles/${slug}`)
  return data
}

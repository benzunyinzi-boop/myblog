import { createRouter, createWebHistory } from 'vue-router'

const PublicLayout = () => import('../layouts/PublicLayout.vue')
const AdminLoginView = () => import('../views/admin/AdminLoginView.vue')
const HomeView = () => import('../views/public/HomeView.vue')
const TechView = () => import('../views/public/TechView.vue')
const AboutView = () => import('../views/public/AboutView.vue')
const ArticleDetailView = () => import('../views/public/ArticleDetailView.vue')

export const router = createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    return { top: 0, behavior: 'smooth' }
  },
  routes: [
    {
      path: '/',
      component: PublicLayout,
      children: [
        { path: '', name: 'home', component: HomeView },
        { path: 'tech', name: 'tech', component: TechView },
        { path: 'about', name: 'about', component: AboutView },
        { path: 'blog/:slug', name: 'article-detail', component: ArticleDetailView, props: true }
      ]
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: AdminLoginView
    }
  ]
})

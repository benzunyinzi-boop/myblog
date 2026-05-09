import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const PublicLayout = () => import('../layouts/PublicLayout.vue')
const AdminLoginView = () => import('../views/admin/AdminLoginView.vue')
const AdminDashboardView = () => import('../views/admin/AdminDashboardView.vue')
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
      component: AdminLoginView,
      meta: { guestOnly: true }
    },
    {
      path: '/admin',
      name: 'admin-dashboard',
      component: AdminDashboardView,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return '/admin/login'
  }

  if (to.meta.guestOnly && auth.isLoggedIn) {
    return '/admin'
  }

  return true
})

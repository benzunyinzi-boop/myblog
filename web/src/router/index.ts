import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const PublicLayout = () => import('../layouts/PublicLayout.vue')
const AdminLayout = () => import('../layouts/AdminLayout.vue')
const AdminLoginView = () => import('../views/admin/AdminLoginView.vue')
const AdminDashboardView = () => import('../views/admin/AdminDashboardView.vue')
const AdminArticlesView = () => import('../views/admin/AdminArticlesView.vue')
const AdminArticleEditView = () => import('../views/admin/AdminArticleEditView.vue')
const AdminCategoriesView = () => import('../views/admin/AdminCategoriesView.vue')
const AdminTagsView = () => import('../views/admin/AdminTagsView.vue')
const AdminProfileView = () => import('../views/admin/AdminProfileView.vue')
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
      component: AdminLayout,
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'admin-dashboard', component: AdminDashboardView },
        { path: 'articles', name: 'admin-articles', component: AdminArticlesView },
        { path: 'articles/new', name: 'admin-article-new', component: AdminArticleEditView },
        { path: 'articles/:id/edit', name: 'admin-article-edit', component: AdminArticleEditView, props: true },
        { path: 'categories', name: 'admin-categories', component: AdminCategoriesView },
        { path: 'tags', name: 'admin-tags', component: AdminTagsView },
        { path: 'profile', name: 'admin-profile', component: AdminProfileView }
      ]
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

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import DashboardView from '@/views/DashboardView.vue'
import LinksView from '@/views/LinksView.vue'
import AnalyticsView from '@/views/AnalyticsView.vue'
import PasswordPromptView from '@/views/PasswordPromptView.vue'
import StatusView from '@/views/StatusView.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { guestOnly: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
    meta: { guestOnly: true },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true },
  },
  {
    path: '/links',
    name: 'Links',
    component: LinksView,
    meta: { requiresAuth: true },
  },
  {
    path: '/analytics/:id',
    name: 'Analytics',
    component: AnalyticsView,
    meta: { requiresAuth: true },
  },
  {
    path: '/p/:slug',
    name: 'PasswordPrompt',
    component: PasswordPromptView,
  },
  {
    path: '/link-not-found',
    name: 'LinkNotFound',
    component: StatusView,
    props: { statusType: 'not-found' },
  },
  {
    path: '/link-inactive',
    name: 'LinkInactive',
    component: StatusView,
    props: { statusType: 'inactive' },
  },
  {
    path: '/link-expired',
    name: 'LinkExpired',
    component: StatusView,
    props: { statusType: 'expired' },
  },
  {
    path: '/link-limit-reached',
    name: 'LinkLimitReached',
    component: StatusView,
    props: { statusType: 'limit' },
  },
  {
    path: '/link-error',
    name: 'LinkError',
    component: StatusView,
    props: { statusType: 'error' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: StatusView,
    props: { statusType: 'not-found' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // If page requires auth and user is not authenticated
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // If page is for guests only (login/register) and user is already logged in
  if (to.meta.guestOnly && authStore.isAuthenticated) {
    next({ name: 'Dashboard' })
    return
  }

  next()
})

export default router

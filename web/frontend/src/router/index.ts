// Vue Router Configuration
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Lazy load views
const ProcessList = () => import('../views/ProcessListView.vue')
const ProcessDetail = () => import('../views/ProcessDetailView.vue')
const Search = () => import('../views/SearchView.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'list',
    component: ProcessList,
    meta: { title: 'Processes' },
  },
  {
    path: '/process/:id',
    name: 'detail',
    component: ProcessDetail,
    meta: { title: 'Process Detail' },
    props: true,
  },
  {
    path: '/search',
    name: 'search',
    component: Search,
    meta: { title: 'Search' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Update page title
router.beforeEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - Tasker` : 'Tasker'
})

export default router

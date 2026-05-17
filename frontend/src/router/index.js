import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from '../pages/DashboardPage.vue'
import LoginPage from '../pages/LoginPage.vue'
import TaskPage from '../pages/TaskPage.vue'
import AdminPage from '../pages/AdminPage.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'login',
    component: LoginPage
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: DashboardPage
  },
  {
    path: '/task/:id',
    name: 'task',
    component: TaskPage
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminPage
  }
]

export default createRouter({
  history: createWebHistory(),
  routes
})


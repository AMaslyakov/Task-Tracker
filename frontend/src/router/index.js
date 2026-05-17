import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from '../pages/DashboardPage.vue'
import LoginPage from '../pages/LoginPage.vue'
import TaskPage from '../pages/TaskPage.vue'

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
  }
]

export default createRouter({
  history: createWebHistory(),
  routes
})


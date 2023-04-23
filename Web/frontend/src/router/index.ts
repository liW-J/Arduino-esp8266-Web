import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        redirect: '/index'
    },
    {
        path: '/index',
        component: () => import('@/pages/client/index.vue')
    },
    {
        path: '/home',
        component: () => import('@/pages/client/home.vue')
    },
    {
        path: '/manage',
        component: () => import('@/pages/client/manage.vue')
    },
    {
        path: '/admin',
        redirect: '/admin/login'
    },
    {
        path: '/admin/login',
        component: () => import('@/pages/admin/login.vue')
    },
    {
        path: '/admin/home',
        component: () => import('@/pages/admin/home.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router
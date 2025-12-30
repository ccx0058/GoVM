import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: '/version'
        },
        {
            path: '/version',
            name: 'version',
            component: () => import('../views/VersionManager.vue')
        },
        {
            path: '/package',
            name: 'package',
            component: () => import('../views/PackageManager.vue')
        },
        {
            path: '/env',
            name: 'env',
            component: () => import('../views/EnvConfig.vue')
        },
        {
            path: '/cache',
            name: 'cache',
            component: () => import('../views/CacheManager.vue')
        },
        {
            path: '/settings',
            name: 'settings',
            component: () => import('../views/Settings.vue')
        }
    ]
})

export default router

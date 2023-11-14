import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import Hosts from "./pages/Hosts.vue";
import Settings from "./pages/Settings.vue";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/hosts',
    },
    {
        path: '/hosts',
        name: 'Hosts',
        component: Hosts,
    },
    {
        path: '/settings',
        name: 'Settings',
        component: Settings,
    },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router

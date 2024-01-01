import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";
import Hosts from "./pages/Hosts.vue";
import Settings from "./pages/Settings.vue";
import About from "./pages/About.vue";

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
    {
        path: '/about',
        name: 'About',
        component: About,
    },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router

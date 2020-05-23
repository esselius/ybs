import Vue from 'vue';
import VueRouter from 'vue-router'

// @ts-ignore
import Login from './components/Login';
// @ts-ignore
import About from './components/About';

Vue.use(VueRouter)

const routes = [
    { path: '/login', component: Login },
    { path: '/about', component: About },
]

const router = new VueRouter({
    routes,
})

const app = new Vue({
    router
}).$mount('#app')

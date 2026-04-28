import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RoomView from '../views/RoomView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',          component: HomeView },
    { path: '/room/:id',  component: RoomView }
  ]
})

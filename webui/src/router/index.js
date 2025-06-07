import {createRouter, createWebHashHistory} from 'vue-router'

import HomeView from '../views/HomeView.vue'

import LoginView from '../views/LoginView.vue'



const router = createRouter({

	history: createWebHashHistory(import.meta.env.BASE_URL),

	routes: [

		{ path: '/', name: 'Login', component: LoginView },

		{ path: '/home', name: 'Home', component: HomeView }, // already there

	  ]

})



export default router


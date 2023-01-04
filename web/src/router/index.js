import { createRouter, createWebHistory } from 'vue-router'
import ListView from '../views/ListView.vue'
import CreateView from '../views/CreateView.vue'
import AboutView from '../views/AboutView.vue'
import LoginView from '../views/LoginView.vue'
import { useAuthStore } from '../stores/auth';



const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/create',
      name: 'Create',
      component: CreateView
    },
    {
      path: '/list',
      name: 'List',
      component: ListView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: AboutView
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginView
    }
  ]
})

router.beforeEach(async (to) => {
  // redirect to login page if not logged in and trying to access a restricted page
  const publicPages = ['/login'];
  const authRequired = !publicPages.includes(to.path);
  const auth = useAuthStore();

  if (authRequired && !auth.token) {
    // auth.returnUrl = to.fullPath;
    return '/login';
  }
});

export default router

import Vue from 'vue'
import Router from 'vue-router'
import Layout from './components/Layout.vue'
import Login from './views/Login.vue'
import routes from './routes'
import store from './store'
import api from './api'

Vue.use(Router)

const buildRoutes = routes => {
  let rs = []
  for (let i = 0; i < routes.length; i++) {
    if (routes[i].children instanceof Array) {
      let _rs = buildRoutes(routes[i].children)
      for (let j = 0; j < _rs.length; j++) {
        rs.push(_rs[j])
      }
    } else {
      rs.push({
        path: routes[i].path.replace(/^\/+/, '').replace(/\/+$/, ''),
        component: routes[i].component,
        meta: {
          role: routes[i].role
        }
      })
    }
  }
  return rs
}

const router = new Router({
  routes: [
    {
      path: '/Login',
      component: Login
    },
    {
      path: '/',
      component: Layout,
      children: buildRoutes(routes)
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  let r = await api.getloginuser()
  if (r.code === 1) {
    store.commit('storeuser', r.data.user)
  }

  if (to.path === '/Login') {
    if (store.state.user) {
      next('/')
    } else {
      next()
    }
    return
  }
  if (store.state.user) {
    // pass
  } else {
    next('/Login')
    return
  }
  for (let i = 0; i < router.options.routes.length; i++) {
    let p1 = router.options.routes[i].path
    let c1 = router.options.routes[i].children
    if (p1 === '/') {
      if (c1 instanceof Array && c1.length > 0) {
        for (let j = 0; j < c1.length; j++) {
          let p2 = c1[j].path
          if (p1 + p2 === to.path) {
            next()
            return
          }
        }
      }
      next(p1 + c1[0].path)
      return
    }
  }
})

export default router

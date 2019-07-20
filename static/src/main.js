import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import 'animate.css/animate.css'

import ace from 'brace'
import 'brace/ext/language_tools'
import 'brace/mode/lua'
import 'brace/snippets/lua'
import 'brace/theme/eclipse'

import { VueAceEditor } from 'vue2x-ace-editor'

Vue.component(VueAceEditor.name, VueAceEditor)

Vue.use(ElementUI)
Vue.config.productionTip = false

new Vue({
  ace,
  router,
  store,
  render: h => h(App)
}).$mount('#app')

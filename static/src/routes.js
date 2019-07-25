import LuaScript from './views/LuaScript.vue'
import Task from './views/Task.vue'

export default [
  { icon: 'el-icon-folder', name: '任务脚本管理', path: '/LuaScript', component: LuaScript },
  { icon: 'el-icon-monitor', name: '任务命令管理', path: '/Task', component: Task }
]

import LuaScript from './views/LuaScript.vue'
import Task from './views/Task.vue'
import Server from './views/Server.vue'

export default [
  { icon: 'el-icon-document', name: '脚本任务管理', path: '/LuaScript', component: LuaScript },
  { icon: 'el-icon-finished', name: '执行命令任务管理', path: '/Task', component: Task },
  { icon: 'el-icon-office-building', name: '服务器管理', path: '/Server', component: Server }
]

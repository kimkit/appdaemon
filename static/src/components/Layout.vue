<template>
  <el-container>
    <el-header>
      <router-link class="logo" to="/">鲁小柒</router-link>
      <el-dropdown @command="userCommand" trigger="click" style="float: right;">
        <div class="user">
          <span class="nickname">{{ $store.state.user }}</span>
          <el-avatar icon="el-icon-user-solid" :size="40"></el-avatar>
        </div>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item icon="el-icon-switch-button" command="/Logout">退出</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </el-header>
    <el-container>
      <el-aside width="200px">
        <el-menu
          background-color="#1a1a1a"
          text-color="#cccccc"
          active-text-color="#eeeeee"
          router
          :default-active="$router.currentRoute.path"
        >
          <template v-for="r in routes">
            <template v-if="r.children">
              <el-submenu index="/" v-bind:key="r.name">
                <template slot="title">
                  <i :class="r.icon"></i>
                  <span>{{ r.name }}</span>
                </template>
                <template v-for="_r in r.children">
                  <el-menu-item
                    v-if="$store.state.role >= _r.role"
                    :index="_r.path"
                    v-bind:key="_r.name"
                  >
                    <i :class="_r.icon"></i>
                    <span>{{ _r.name }}</span>
                  </el-menu-item>
                </template>
              </el-submenu>
            </template>
            <template v-else>
              <el-menu-item :index="r.path" v-bind:key="r.name">
                <i :class="r.icon"></i>
                <span>{{ r.name }}</span>
              </el-menu-item>
            </template>
          </template>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import routes from '../routes'
import api from '../api'
export default {
  name: 'Layout',
  data () {
    return {
      routes
    }
  },
  methods: {
    userCommand (path) {
      if (path === '/Logout') {
        api.removetoken()
        this.$store.commit('storeuser', '')
        this.$router.push('/Login')
      } else {
        this.$router.push(path)
      }
    }
  }
}
</script>

<style>
.el-container {
  height: 100%;
  background-color: #eceef2;
}
.el-header {
  background-color: #ffffff;
  box-shadow: rgba(0, 21, 41, 0.08) 0px 1px 4px;
}
.el-aside {
  background-color: #1a1a1a;
  box-shadow: rgba(0, 21, 41, 0.38) 0px 1px 4px;
  width: 200px;
  overflow-y: scroll;
  height: 100%;
}
.el-main {
  padding: 10px !important;
}
.el-menu {
  border-right: none !important;
}
.el-menu-item {
  height: 40px !important;
  line-height: 40px !important;
  margin-top: 10px;
}
.el-menu-item:hover i {
  color: inherit;
}
.el-menu-item:hover {
  color: #eeeeee !important;
}
.el-menu-item.is-active {
  background-color: #409eff !important;
  color: #eeeeee !important;
}
.el-submenu__title {
  height: 40px !important;
  line-height: 40px !important;
  margin-top: 10px;
}
.logo {
  float: left;
  height: 60px;
  padding-left: 70px;
  cursor: pointer;
  line-height: 60px;
  font-size: 20px;
  text-decoration: none;
  color: #4086d2;
  background-image: url(../assets/logo.png);
  background-size: 50px;
  background-repeat: no-repeat;
  background-position: 5px center;
}
.user {
  height: 48px;
  padding-top: 10px;
  cursor: pointer;
}
.user .nickname {
  float: left;
  line-height: 40px;
  font-size: 14px;
  margin-right: 10px;
}
</style>

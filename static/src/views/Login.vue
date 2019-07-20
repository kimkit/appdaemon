<template>
  <div class="login">
    <div class="login-box">
      <el-form @submit.native.prevent>
        <el-form-item>
          <el-input
            size="medium"
            v-model="password"
            type="password"
            placeholder="请输入(账号:密码)"
            prefix-icon="el-icon-user-solid"
            suffix-icon="el-icon-right"
            @keyup.enter.native="login"
          ></el-input>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import api from '../api'
export default {
  name: 'Login',
  data () {
    return {
      password: ''
    }
  },
  methods: {
    animateCSS (element, animationName, callback) {
      const node = document.querySelector(element)
      node.classList.add('animated', animationName)

      function handleAnimationEnd () {
        node.classList.remove('animated', animationName)
        node.removeEventListener('animationend', handleAnimationEnd)

        if (typeof callback === 'function') callback()
      }

      node.addEventListener('animationend', handleAnimationEnd)
    },
    shake () {
      this.animateCSS('.login-box', 'shake', () => {
        this.$data.password = ''
      })
    },
    async login () {
      let pos = this.$data.password.indexOf(':')
      if (pos < 0) {
        this.shake()
        return
      }
      let username = this.$data.password.substr(0, pos)
      let password = this.$data.password.substr(pos + 1)
      let r = await api.login(username, password)
      if (r.code === 1) {
        this.$router.push('/')
      } else {
        this.shake()
      }
    }
  }
}
</script>

<style>
.login {
  height: 100%;
  background-color: #eceef2;
  padding-top: 160px;
  background-image: url(../assets/logo.png);
  background-repeat: no-repeat;
  background-size: 80px;
  background-position: center 60px;
}
.login-box {
  width: 320px;
  padding: 20px;
  margin: 0 auto;
  background-color: #ffffff;
  box-shadow: rgba(0, 21, 41, 0.08) 0px 1px 4px;
  border-radius: 2px;
}
.login-box .el-form-item {
  margin-bottom: 0;
}
</style>

<template>
  <el-dialog title="日志" :visible.sync="outputshow" width="800px" @close="outputhide">
    <el-form :inline="true" @submit.native.prevent size="small">
      <el-form-item>
        <el-select
          class="addrlist"
          v-model="subaddr"
          filterable
          clearable
          placeholder="留空查找所有服务器"
          @change="changeoutput"
        >
          <el-option
            v-for="item in addrlist"
            :key="item.addr"
            :label="item.addr"
            :value="item.addr"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-input
          v-model="jobname"
          type="text"
          placeholder="请输入原始任务名称"
          prefix-icon="el-icon-search"
          suffix-icon="el-icon-right"
          class="scriptname-input-wide"
          @keyup.enter.native="changeoutput"
        ></el-input>
      </el-form-item>
    </el-form>
    <div class="output">
      <p>
        <i class="el-icon-loading"></i>
      </p>
    </div>
    <div class="output" v-for="line in outputlines" v-bind:key="line.id">
      <p class="outputtitle" v-if="line.title" @click="line.hide = line.hide ? false : true">
        <i :class="line.hide ? 'el-icon-arrow-up' : 'el-icon-arrow-right'"></i>
        <strong>{{ line.title }}</strong>
      </p>
      <p v-if="line.data" :class="line.hide ? 'outputfold' : ''">{{ line.data }}</p>
    </div>
  </el-dialog>
</template>

<script>
import api from '../api'
export default {
  name: 'Output',
  data () {
    return {
      addrlist: [],
      jobname: '',
      subaddr: '',
      outputshow: false,
      outputlines: []
    }
  },
  async mounted () {
    let r = await api.getserverlist()
    if (r.code === 1) {
      this.$data.addrlist = r.data.list
    }
  },
  methods: {
    getoutput (jobname, subaddr) {
      this.$data.outputshow = true
      this.$data.jobname = jobname
      this.$data.subaddr = subaddr
      if (this.ws) {
        this.ws.close()
      }
      let url = api
        .geturl('GetOutput', { token: api.gettoken(), jobname, subaddr })
        .replace(/^http/i, 'ws')
      if (!url.match(/^ws/i)) {
        url =
          document.location.protocol.replace(/^http/i, 'ws') +
          '//' +
          document.location.host +
          url
      }
      this.$data.outputshow = true
      this.$data.outputlines = []
      this.ws = new WebSocket(url)
      this.ws.onopen = e => {
        console.log(new Date(), e)
        this.$data.outputlines.unshift({
          id: Math.random()
            .toString(36)
            .slice(-8),
          data: 'Connected'
        })
        this.$data.outputlines = this.$data.outputlines.slice(0, 100)
      }
      this.ws.onclose = e => {
        console.log(new Date(), e)
        this.$data.outputlines.unshift({
          id: Math.random()
            .toString(36)
            .slice(-8),
          data: 'Connection closed'
        })
        this.$data.outputlines = this.$data.outputlines.slice(0, 100)
      }
      this.ws.onerror = e => {
        console.log(new Date(), e)
        this.$data.outputlines.unshift({
          id: Math.random()
            .toString(36)
            .slice(-8),
          data: 'Error occur'
        })
        this.$data.outputlines = this.$data.outputlines.slice(0, 100)
      }
      this.ws.onmessage = e => {
        console.log(new Date(), e)
        let r = JSON.parse(e.data)
        if (r.code === 1) {
          this.$data.outputlines.unshift({
            id: Math.random()
              .toString(36)
              .slice(-8),
            data: r.data
          })
        } else if (r.code === -2) {
          // ignore
        } else {
          this.$data.outputlines.unshift({
            id: Math.random()
              .toString(36)
              .slice(-8),
            data: r.message + ' (' + r.code + ')'
          })
        }
        this.$data.outputlines = this.$data.outputlines.slice(0, 100)
      }
    },
    outputhide () {
      if (this.ws) {
        this.ws.close()
        this.ws = null
      }
    },
    outpushshow () {
      this.$data.outpushshow = true
      this.$data.jobname = ''
      this.$data.subaddr = ''
      this.outputhide()
    },
    changeoutput () {
      this.getoutput(this.$data.jobname, this.$data.subaddr)
    }
  }
}
</script>

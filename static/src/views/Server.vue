<template>
  <el-card class="box-card">
    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="name" label="服务器" width="220">
        <template slot-scope="scope">
          <span class="addr" v-if="scope.row.addr">
            <i class="el-icon-arrow-right"></i>
            {{scope.row.addr}}
          </span>
          <span class="addr" v-else>
            <i class="el-icon-check"></i>
            all
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="更新时间" width="220">
        <template slot-scope="scope">
          <span class="addr">{{scope.row.updatetime}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="任务数" width="100">
        <template slot-scope="scope">
          <span class="addr">{{scope.row.jobcount}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template slot-scope="scope">
          <el-tag type="success" size="mini" v-if="scope.row.status == 1" style="cursor: auto;">启用</el-tag>
          <el-tag type="info" size="mini" v-if="scope.row.status == 0" style="cursor: auto;">停用</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" min-width="150">
        <template slot-scope="scope">
          <el-button
            class="el-button-action"
            icon="el-icon-delete"
            size="mini"
            circle
            @click="deleteserver(scope.row.id, scope.row.addr)"
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
import api from '../api'
export default {
  name: 'Server',
  data () {
    return {
      list: [],
      loading: false
    }
  },
  async mounted () {
    await this.getserverlist()
  },
  methods: {
    async getserverlist () {
      let r = await api.getserverlist(1)
      if (r.code === 1) {
        this.$data.list = r.data.list
      }
    },
    async updateserverstatus (id, status) {
      let r = await api.updateserverstatus(id, status)
      if (r.code === 1) {
        await this.getserverlist(1)
      }
    },
    async deleteserver (id, addr) {
      let r = await this.$confirm(
        '确认要删除服务器：<code>' + addr + '</code> ？',
        '提示',
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).catch(() => {
        // pass
      })
      if (r === 'confirm') {
        r = await api.deleteserver(id)
        if (r.code === 1) {
          this.getserverlist(1)
        }
      }
    }
  }
}
</script>

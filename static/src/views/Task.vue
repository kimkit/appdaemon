<template>
  <el-card class="box-card">
    <el-form :inline="true" @submit.native.prevent size="small">
      <el-form-item>
        <el-select
          class="addrlist"
          v-model="addr"
          filterable
          clearable
          placeholder="留空查找所有服务器"
          @change="searchtasklist"
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
          v-model="keyword"
          type="text"
          placeholder="请输入查找任务名称"
          prefix-icon="el-icon-search"
          suffix-icon="el-icon-right"
          class="scriptname-input"
          @keyup.enter.native="searchtasklist"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-document-add" @click="addtaskshow">添加</el-button>
      </el-form-item>
    </el-form>
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
      <el-table-column prop="name" label="任务名称 / 规则 / 命令" width="320">
        <template slot-scope="scope">
          <div class="task-info">{{scope.row.name}}</div>
          <div class="task-info">{{scope.row.rule}}</div>
          <div class="task-info task-info-command">{{scope.row.command}}</div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="执行 / 上报时间" width="220">
        <template slot-scope="scope">
          <span class="addr">{{scope.row.updatetime}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template slot-scope="scope">
          <el-tag type="success" size="mini" v-if="scope.row.status == 1">启用</el-tag>
          <el-tag type="info" size="mini" v-if="scope.row.status == 0">停用</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" min-width="150">
        <template slot-scope="scope">
          <el-button
            class="el-button-action"
            icon="el-icon-edit"
            size="mini"
            circle
            @click="updatetaskshow(scope.row.id)"
          ></el-button>
          <el-button
            class="el-button-action"
            icon="el-icon-delete"
            size="mini"
            circle
            @click="deletetask(scope.row.id, scope.row.name)"
          ></el-button>
          <el-button
            class="el-button-action"
            icon="el-icon-toilet-paper"
            size="mini"
            circle
            @click="getoutputlist(scope.row.name)"
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      background
      layout="prev, pager, next, jumper"
      :total="total"
      :page-size="pagesize"
      :current-page="page"
      @current-change="gettasklist"
      style="margin-top: 20px;"
    />
    <el-dialog title="添加" :visible.sync="addshow" width="800px">
      <el-form :inline="true" size="small">
        <el-form-item>
          <el-select
            class="addrlist"
            v-model="addaddr"
            filterable
            clearable
            placeholder="留空所有服务器执行"
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
            v-model="addname"
            prefix-icon="el-icon-document"
            placeholder="任务名称"
            class="scriptname-input"
          ></el-input>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="addstatus"></el-switch>
        </el-form-item>
      </el-form>
      <el-form size="small">
        <el-form-item>
          <el-input class="edit-input" v-model="addrule" placeholder="定时任务规则或者执行进程数量"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input
            type="textarea"
            :autosize="{ minRows: 4}"
            placeholder="任务命令(不支持管道,重定向,环境变量,内嵌命令)"
            v-model="addcommand"
            class="edit-input"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="medium" @click="addshow = false">取消</el-button>
        <el-button size="medium" type="primary" @click="addtask">确定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="编辑" :visible.sync="updateshow" width="800px">
      <el-form :inline="true" size="small">
        <el-form-item>
          <el-select
            class="addrlist"
            v-model="updateaddr"
            filterable
            clearable
            placeholder="留空所有服务器执行"
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
            v-model="updatename"
            prefix-icon="el-icon-document"
            placeholder="任务名称"
            class="scriptname-input"
          ></el-input>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="updatestatus"></el-switch>
        </el-form-item>
      </el-form>
      <el-form size="small">
        <el-form-item>
          <el-input class="edit-input" v-model="updaterule" placeholder="定时任务规则或者执行进程数量"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input
            type="textarea"
            :autosize="{ minRows: 4}"
            placeholder="任务命令(不支持管道,重定向,环境变量,内嵌命令)"
            v-model="updatecommand"
            class="edit-input"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="medium" @click="updateshow = false">取消</el-button>
        <el-button size="medium" type="primary" @click="updatetask">确定</el-button>
      </div>
    </el-dialog>
  </el-card>
</template>

<script>
import api from '../api'
export default {
  name: 'Task',
  data () {
    return {
      list: [],
      total: 0,
      pagesize: 8,
      page: 1,
      keyword: '',
      loading: false,
      addrlist: [],
      addr: '',
      addshow: false,
      addaddr: '',
      addname: '',
      addrule: '',
      addcommand: '',
      addstatus: false,
      updateshow: false,
      updateaddr: '',
      updateid: 0,
      updatename: '',
      updatestatus: false,
      updaterule: '',
      updatecommand: ''
    }
  },
  async mounted () {
    await this.gettasklist()
    let r = await api.getserverlist()
    if (r.code === 1) {
      this.$data.addrlist = r.data.list
    }
  },
  methods: {
    async gettasklist (page) {
      if (page) {
        this.$data.page = page
      }
      let r = await api.gettasklist(
        this.$data.page,
        this.$data.pagesize,
        this.$data.keyword,
        this.$data.addr
      )
      if (r.code === 1) {
        this.$data.list = r.data.list
        this.$data.total = r.data.total
        this.$data.page = r.data.page
        this.$data.pagesize = r.data.pagesize
      }
    },
    async searchtasklist () {
      await this.gettasklist(1)
    },
    addtaskshow () {
      this.$data.addname =
        'tmp_' +
        Math.random()
          .toString(36)
          .slice(-8)
      this.$data.addrule = '*/10 * * * * * *'
      this.$data.addcommand = 'echo hello'
      this.$data.addshow = true
    },
    async addtask () {
      let r = await api.addtask(
        this.$data.addname,
        this.$data.addrule,
        this.$data.addcommand,
        this.$data.addstatus ? 1 : 0,
        this.$data.addaddr
      )
      if (r.code === 1) {
        this.$message({
          type: 'success',
          message: '添加任务命令成功（*＾-＾*）',
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
        setTimeout(() => {
          this.$data.addname = ''
          this.$data.addstatus = false
          this.$data.addrule = ''
          this.$data.addcommand = ''
          this.$data.addshow = false
          this.$data.keyword = ''
          this.$data.addr = ''
          this.gettasklist(1)
        }, 500)
      } else {
        this.$message({
          type: 'error',
          message: '(' + r.code + ') ' + r.message,
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
      }
    },
    async updatetaskshow (id) {
      let r = await api.gettask(id)
      if (r.code === 1) {
        this.$data.updateid = r.data.id
        this.$data.updatename = r.data.name
        this.$data.updatestatus = r.data.status === '1'
        this.$data.updaterule = r.data.rule
        this.$data.updatecommand = r.data.command
        this.$data.updateaddr = r.data.addr
        this.$data.updateshow = true
      }
    },
    async updatetask () {
      let r = await api.updatetask(
        this.$data.updateid,
        this.$data.updatename,
        this.$data.updaterule,
        this.$data.updatecommand,
        this.$data.updatestatus ? 1 : 0,
        this.$data.updateaddr
      )
      if (r.code === 1) {
        this.$message({
          type: 'success',
          message: '更新任务命令成功（*＾-＾*）',
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
        setTimeout(() => {
          this.$data.updateshow = false
          this.gettasklist()
        }, 500)
      } else {
        this.$message({
          type: 'error',
          message: '(' + r.code + ') ' + r.message,
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
      }
    },
    async deletetask (id, name) {
      let r = await this.$confirm(
        '确认要删除任务：<code>' + name + '</code> ？',
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
        r = await api.deletetask(id)
        if (r.code === 1) {
          this.gettasklist()
        }
      }
    },
    getoutputlist (name) {}
  }
}
</script>

<style>
.task-info {
  font-family: "Monaco", "Menlo", "Ubuntu Mono", "Consolas", "source-code-pro",
    monospace;
  border-bottom: 1px solid #eeeeee;
  font-size: 13px;
  word-wrap: break-word;
  word-break: break-all;
  white-space: pre-wrap;
}
.task-info-command {
  border-bottom: none;
}
</style>

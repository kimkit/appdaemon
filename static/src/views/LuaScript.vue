<template>
  <el-card class="box-card">
    <el-form :inline="true" @submit.native.prevent size="small">
      <el-form-item style="padding-top: 1px;">
        <el-select
          class="addrlist"
          v-model="addr"
          filterable
          clearable
          placeholder="留空查找所有服务器"
          @change="searchluascriptlist"
        >
          <el-option
            v-for="item in addrlist"
            :key="item.addr"
            :label="item.addr"
            :value="item.addr"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item style="padding-top: 1px;">
        <el-input
          v-model="keyword"
          type="text"
          placeholder="请输入查找脚本名称"
          prefix-icon="el-icon-search"
          suffix-icon="el-icon-right"
          class="scriptname-input"
          @keyup.enter.native="searchluascriptlist"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-document-add" @click="addluascriptshow">添加</el-button>
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
      <el-table-column prop="name" label="脚本名称" width="320">
        <template slot-scope="scope">
          <code>{{scope.row.name}}</code>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template slot-scope="scope">
          <el-tag
            type="success"
            size="mini"
            v-if="scope.row.status == 1"
            @click="updateluascriptstatus(scope.row.id, 0)"
          >启用</el-tag>
          <el-tag
            type="info"
            size="mini"
            v-if="scope.row.status == 0"
            @click="updateluascriptstatus(scope.row.id, 1)"
          >停用</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" min-width="150">
        <template slot-scope="scope">
          <el-button
            class="el-button-action"
            icon="el-icon-edit"
            size="mini"
            circle
            @click="updateluascriptshow(scope.row.id)"
          ></el-button>
          <el-button
            class="el-button-action"
            icon="el-icon-delete"
            size="mini"
            circle
            @click="deleteluascript(scope.row.id, scope.row.name)"
          ></el-button>
          <el-button
            class="el-button-action"
            icon="el-icon-toilet-paper"
            size="mini"
            circle
            @click="getoutput(scope.row.jobname, scope.row.subaddr)"
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
      @current-change="getluascriptlist"
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
            placeholder="脚本名称"
            class="scriptname-input"
          ></el-input>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="addstatus"></el-switch>
        </el-form-item>
      </el-form>
      <vue-ace-editor
        ref="addeditor"
        :content="addscript"
        :fontSize="12"
        height="257px"
        lang="lua"
        theme="eclipse"
        :options="editoroptions"
        class="scripteditor"
      ></vue-ace-editor>
      <div slot="footer" class="dialog-footer">
        <el-button size="medium" @click="addshow = false">取消</el-button>
        <el-button size="medium" type="primary" @click="addluascript">确定</el-button>
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
            placeholder="脚本名称"
            class="scriptname-input"
          ></el-input>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="updatestatus"></el-switch>
        </el-form-item>
      </el-form>
      <vue-ace-editor
        ref="updateeditor"
        :content="updatescript"
        :fontSize="12"
        height="257px"
        lang="lua"
        theme="eclipse"
        :options="editoroptions"
        class="scripteditor"
      ></vue-ace-editor>
      <div slot="footer" class="dialog-footer">
        <el-button size="medium" @click="updateshow = false">取消</el-button>
        <el-button size="medium" type="primary" @click="updateluascript">确定</el-button>
      </div>
    </el-dialog>
    <Output ref="output"></Output>
  </el-card>
</template>

<script>
import api from '../api'
import Output from '../components/Output.vue'
export default {
  name: 'LuaScript',
  components: {
    Output
  },
  data () {
    return {
      list: [],
      total: 0,
      pagesize: 8,
      page: 1,
      keyword: '',
      loading: false,
      editoroptions: { wrap: 'free' },
      addrlist: [],
      addr: '',
      addshow: false,
      addaddr: '',
      addname: '',
      addstatus: false,
      addscript: '',
      updateshow: false,
      updateaddr: '',
      updateid: 0,
      updatename: '',
      updatestatus: false,
      updatescript: ''
    }
  },
  async mounted () {
    await this.getluascriptlist()
    let r = await api.getserverlist()
    if (r.code === 1) {
      this.$data.addrlist = r.data.list
    }
  },
  methods: {
    async getluascriptlist (page) {
      if (page) {
        this.$data.page = page
      }
      let r = await api.getluascriptlist(
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
    async searchluascriptlist () {
      await this.getluascriptlist(1)
    },
    addluascriptshow () {
      this.$data.addname =
        'tmp_' +
        Math.random()
          .toString(36)
          .slice(-8)
      this.$data.addscript = `if cron == nil then cron = newcron("*/10 * * * * * *") end
if nexttime == nil then nexttime = cron:next() end

now = os.time()
if now >= nexttime then
    log.debug("%v: ...", jobname)
    nexttime = cron:next()
else
    sleep(200)
end`
      this.$data.addshow = true
    },
    async addluascript () {
      let r = await api.addluascript(
        this.$data.addname,
        this.$refs.addeditor.editor.getValue(),
        this.$data.addstatus ? 1 : 0,
        this.$data.addaddr
      )
      if (r.code === 1) {
        this.$message({
          type: 'success',
          message: '添加任务脚本成功（*＾-＾*）',
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
        setTimeout(() => {
          this.$data.addname = ''
          this.$data.addstatus = false
          this.$data.addscript = ''
          this.$refs.addeditor.editor.setValue('')
          this.$data.addshow = false
          this.$data.keyword = ''
          this.$data.addr = ''
          this.getluascriptlist(1)
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
    async updateluascriptstatus (id, status) {
      let r = await api.updateluascriptstatus(id, status)
      if (r.code === 1) {
        await this.getluascriptlist()
      }
    },
    async updateluascriptshow (id) {
      let r = await api.getluascript(id)
      if (r.code === 1) {
        this.$data.updateid = r.data.id
        this.$data.updatename = r.data.name
        this.$data.updatestatus = r.data.status === '1'
        this.$data.updatescript = r.data.script
        this.$data.updateaddr = r.data.addr
        this.$data.updateshow = true
      }
    },
    async updateluascript () {
      let r = await api.updateluascript(
        this.$data.updateid,
        this.$data.updatename,
        this.$refs.updateeditor.editor.getValue(),
        this.$data.updatestatus ? 1 : 0,
        this.$data.updateaddr
      )
      if (r.code === 1) {
        this.$message({
          type: 'success',
          message: '更新任务脚本成功（*＾-＾*）',
          offset: 12,
          duration: 1000,
          customClass: 'message'
        })
        setTimeout(() => {
          this.$data.updateshow = false
          this.getluascriptlist()
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
    async deleteluascript (id, name) {
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
        r = await api.deleteluascript(id)
        if (r.code === 1) {
          this.getluascriptlist()
        }
      }
    },
    getoutput (jobname, subaddr) {
      this.$refs.output.getoutput(jobname, subaddr)
    }
  }
}
</script>

<style>
.el-tag {
  cursor: pointer;
}
code {
  font-family: "Monaco", "Menlo", "Ubuntu Mono", "Consolas", "source-code-pro",
    monospace;
  color: #5e6d82;
  background-color: #e6effb;
  margin: 0;
  display: inline-block;
  padding: 1px 5px;
  font-size: 12px;
  border-radius: 3px;
  height: 18px;
  line-height: 18px;
}
.scriptname-input {
  width: 310px !important;
}
.scriptname-input-wide {
  width: 410px !important;
}
.scriptname-input .el-input__inner,
.scriptname-input-wide .el-input__inner,
.addrlist input,
.el-select-dropdown__item,
.addr,
.edit-input .el-input__inner,
.edit-input .el-textarea__inner {
  font-family: "Monaco", "Menlo", "Ubuntu Mono", "Consolas", "source-code-pro",
    monospace !important;
  font-size: 13px !important;
}
.addrlist {
  width: 220px;
}
.scripteditor {
  border: 1px solid #dcdfe6;
}
.output {
  word-wrap: break-word;
  word-break: break-all;
  white-space: pre-wrap;
}
.output p {
  font-family: "Monaco", "Menlo", "Ubuntu Mono", "Consolas", "source-code-pro",
    monospace;
  border-bottom: 1px solid #eeeeee;
  font-size: 13px;
}
.outputtitle {
  cursor: pointer;
}
.outputfold {
  display: none;
}
</style>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  GetEnvInfo, 
  DiagnoseEnv, 
  SetEnvVar,
  FixGoRoot,
  FixGoProxy,
  GetInstalledVersions,
  SelectDirectory
} from '../../wailsjs/go/main/App'

interface EnvInfo {
  goroot: string
  gopath: string
  goproxy: string
  gobin: string
  path: string
}

interface DiagnoseResult {
  item: string
  status: string
  message: string
}

interface InstalledVersion {
  version: string
  path: string
  isCurrent: boolean
}

const envInfo = ref<EnvInfo>({
  goroot: '',
  gopath: '',
  goproxy: '',
  gobin: '',
  path: ''
})

const diagnoseResults = ref<DiagnoseResult[]>([])
const installedVersions = ref<InstalledVersion[]>([])
const loading = ref(false)
const diagnosing = ref(false)
const fixing = ref(false)

// 编辑相关
const editDialog = ref(false)
const editForm = ref({
  name: '',
  value: ''
})

const loadEnvInfo = async () => {
  loading.value = true
  try {
    const [env, versions] = await Promise.all([
      GetEnvInfo(),
      GetInstalledVersions()
    ])
    envInfo.value = env
    installedVersions.value = versions || []
  } catch (e: any) {
    ElMessage.error('获取环境信息失败: ' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

const runDiagnose = async () => {
  diagnosing.value = true
  try {
    diagnoseResults.value = await DiagnoseEnv()
  } catch (e: any) {
    ElMessage.error('诊断失败: ' + (e?.message || e))
  } finally {
    diagnosing.value = false
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制到剪贴板')
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'ok': return 'success'
    case 'warning': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

// 打开编辑对话框
const openEditDialog = (name: string, currentValue: string) => {
  editForm.value = {
    name,
    value: currentValue || ''
  }
  editDialog.value = true
}

// 保存环境变量
const saveEnvVar = async () => {
  try {
    await SetEnvVar(editForm.value.name, editForm.value.value)
    ElMessage.success(`${editForm.value.name} 已更新，重启终端后生效`)
    editDialog.value = false
    await loadEnvInfo()
    await runDiagnose()
  } catch (e: any) {
    ElMessage.error('设置失败: ' + (e?.message || e))
  }
}

// 选择目录
const selectDir = async () => {
  try {
    const dir = await SelectDirectory(`选择 ${editForm.value.name} 目录`)
    if (dir) {
      editForm.value.value = dir
    }
  } catch (e) {
    // 用户取消
  }
}

// 一键修复 GOROOT
const handleFixGoRoot = async () => {
  if (installedVersions.value.length === 0) {
    ElMessage.warning('未检测到已安装的 Go 版本')
    return
  }
  
  const version = installedVersions.value.find(v => v.isCurrent) || installedVersions.value[0]
  
  try {
    await ElMessageBox.confirm(
      `将 GOROOT 设置为: ${version.path}`,
      '确认修复',
      { type: 'info' }
    )
    fixing.value = true
    await FixGoRoot(version.path)
    ElMessage.success('GOROOT 已设置，重启终端后生效')
    await loadEnvInfo()
    await runDiagnose()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('修复失败: ' + (e?.message || e))
    }
  } finally {
    fixing.value = false
  }
}

// 一键修复 GOPROXY
const handleFixGoProxy = async () => {
  try {
    await ElMessageBox.confirm(
      '将 GOPROXY 设置为: https://goproxy.cn,direct',
      '确认修复',
      { type: 'info' }
    )
    fixing.value = true
    await FixGoProxy()
    ElMessage.success('GOPROXY 已设置，重启终端后生效')
    await loadEnvInfo()
    await runDiagnose()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('修复失败: ' + (e?.message || e))
    }
  } finally {
    fixing.value = false
  }
}

onMounted(() => {
  loadEnvInfo()
  // 监听版本切换事件，自动刷新环境信息
  EventsOn('version-changed', () => {
    loadEnvInfo()
    runDiagnose()
  })
})

onUnmounted(() => {
  EventsOff('version-changed')
})
</script>

<template>
  <div class="env-config">
    <div class="section">
      <div class="section-header">
        <h3>环境变量</h3>
        <el-button :icon="'Refresh'" circle size="small" @click="loadEnvInfo" :loading="loading" />
      </div>

      <el-table :data="[
        { name: 'GOROOT', value: envInfo.goroot, desc: 'Go 安装目录' },
        { name: 'GOPATH', value: envInfo.gopath, desc: '工作区目录' },
        { name: 'GOPROXY', value: envInfo.goproxy, desc: '模块代理' },
        { name: 'GOBIN', value: envInfo.gobin, desc: '可执行文件目录' }
      ]" v-loading="loading">
        <el-table-column prop="name" label="变量名" width="120" />
        <el-table-column prop="value" label="当前值">
          <template #default="{ row }">
            <span :class="{ 'not-set': !row.value }">{{ row.value || '未设置' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button 
              v-if="row.value" 
              :icon="'CopyDocument'" 
              size="small" 
              text 
              @click="copyToClipboard(row.value)" 
            />
            <el-button 
              type="primary" 
              size="small" 
              text 
              @click="openEditDialog(row.name, row.value)"
            >
              配置
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="section">
      <div class="section-header">
        <h3>环境诊断</h3>
        <div class="header-actions">
          <el-button type="primary" @click="runDiagnose" :loading="diagnosing">
            开始诊断
          </el-button>
        </div>
      </div>

      <el-table :data="diagnoseResults" v-if="diagnoseResults.length > 0">
        <el-table-column prop="item" label="检查项" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ row.status === 'ok' ? '正常' : row.status === 'warning' ? '警告' : '错误' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="详情" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button 
              v-if="row.status !== 'ok' && row.item === 'GOROOT'"
              type="primary"
              size="small"
              text
              @click="handleFixGoRoot"
              :loading="fixing"
            >
              修复
            </el-button>
            <el-button 
              v-if="row.status !== 'ok' && row.item === 'GOPROXY'"
              type="primary"
              size="small"
              text
              @click="handleFixGoProxy"
              :loading="fixing"
            >
              修复
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-else description="点击「开始诊断」检查环境配置" />
    </div>

    <div class="section">
      <div class="section-header">
        <h3>快捷操作</h3>
      </div>
      
      <div class="quick-actions">
        <el-button type="primary" @click="handleFixGoRoot" :loading="fixing" :disabled="installedVersions.length === 0">
          <el-icon><Setting /></el-icon>
          自动配置 GOROOT
        </el-button>
        <el-button type="success" @click="handleFixGoProxy" :loading="fixing">
          <el-icon><Connection /></el-icon>
          配置国内代理 (goproxy.cn)
        </el-button>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialog" :title="`配置 ${editForm.name}`" width="500px">
      <el-form label-width="80px">
        <el-form-item label="变量名">
          <el-input v-model="editForm.name" disabled />
        </el-form-item>
        <el-form-item label="值">
          <div class="value-input">
            <el-input 
              v-model="editForm.value" 
              :placeholder="`请输入 ${editForm.name} 的值`"
              clearable
            />
            <el-button 
              v-if="editForm.name === 'GOROOT' || editForm.name === 'GOPATH' || editForm.name === 'GOBIN'"
              @click="selectDir"
            >
              浏览
            </el-button>
          </div>
        </el-form-item>
        <el-form-item v-if="editForm.name === 'GOPROXY'">
          <div class="preset-values">
            <span>常用值：</span>
            <el-tag 
              class="preset-tag" 
              @click="editForm.value = 'https://goproxy.cn,direct'"
              style="cursor: pointer"
            >
              goproxy.cn
            </el-tag>
            <el-tag 
              class="preset-tag" 
              @click="editForm.value = 'https://proxy.golang.org,direct'"
              style="cursor: pointer"
            >
              官方代理
            </el-tag>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog = false">取消</el-button>
        <el-button type="primary" @click="saveEnvVar">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.env-config {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  border: 1px solid #e8e8e8;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.not-set {
  color: #909399;
}

.quick-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.preset-values {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preset-values span {
  color: #909399;
  font-size: 12px;
}

.preset-tag:hover {
  background: #409eff;
  color: #fff;
  cursor: pointer;
}

.value-input {
  display: flex;
  gap: 8px;
  width: 100%;
}

.value-input .el-input {
  flex: 1;
}
</style>

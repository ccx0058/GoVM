<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import {
  GetRemoteVersions,
  GetInstalledVersions,
  InstallVersion,
  UninstallVersion,
  UseVersion,
  GetCurrentVersion,
  GetInstallSettings,
  SetInstallSettings,
  SelectDirectory
} from '../../wailsjs/go/main/App'

interface GoVersion {
  version: string
  stable: boolean
}

interface InstalledVersion {
  version: string
  path: string
  isCurrent: boolean
}

interface DownloadProgress {
  total: number
  downloaded: number
  percent: number
}

interface InstallSettings {
  installDir: string
  cacheDir: string
}

const remoteVersions = ref<GoVersion[]>([])
const installedVersions = ref<InstalledVersion[]>([])
const currentVersion = ref('')
const loading = ref(false)
const installing = ref(false)
const installingVersion = ref('')
const downloadProgress = ref(0)
const showStableOnly = ref(true)

// 安装设置
const installSettings = ref<InstallSettings>({
  installDir: '',
  cacheDir: ''
})
const settingsDialog = ref(false)
const tempSettings = ref<InstallSettings>({
  installDir: '',
  cacheDir: ''
})

// 安装确认对话框
const installDialog = ref(false)
const selectedVersion = ref('')
const customInstallDir = ref('')

const filteredRemoteVersions = computed(() => {
  if (showStableOnly.value) {
    return remoteVersions.value.filter(v => v.stable)
  }
  return remoteVersions.value
})

const isInstalled = (version: string) => {
  return installedVersions.value.some(v => v.version === version)
}

const loadData = async () => {
  loading.value = true
  try {
    const [remote, installed, current, settings] = await Promise.all([
      GetRemoteVersions(!showStableOnly.value),
      GetInstalledVersions(),
      GetCurrentVersion(),
      GetInstallSettings()
    ])
    remoteVersions.value = remote || []
    installedVersions.value = installed || []
    currentVersion.value = current
    installSettings.value = settings
  } catch (e: any) {
    const msg = e?.message || e?.toString() || '未知错误'
    ElMessage.error('加载数据失败: ' + msg)
  } finally {
    loading.value = false
  }
}

// 打开安装确认对话框
const openInstallDialog = (version: string) => {
  selectedVersion.value = version
  customInstallDir.value = installSettings.value.installDir
  installDialog.value = true
}

// 选择安装目录
const selectInstallDir = async () => {
  try {
    const dir = await SelectDirectory('选择安装目录')
    if (dir) {
      customInstallDir.value = dir
    }
  } catch (e) {
    // 用户取消选择
  }
}

// 确认安装
const confirmInstall = async () => {
  installDialog.value = false
  await handleInstall(selectedVersion.value, customInstallDir.value)
}

const handleInstall = async (version: string, installDir?: string) => {
  installing.value = true
  installingVersion.value = version
  downloadProgress.value = 0

  try {
    await InstallVersion(version, installDir || '')
    ElMessage.success(`${version} 安装成功`)
    await loadData()
  } catch (e: any) {
    ElMessage.error('安装失败: ' + (e?.message || e))
  } finally {
    installing.value = false
    installingVersion.value = ''
  }
}

const handleUninstall = async (version: string) => {
  try {
    await ElMessageBox.confirm(`确定要卸载 ${version} 吗？`, '确认卸载', {
      type: 'warning'
    })
    await UninstallVersion(version)
    ElMessage.success(`${version} 已卸载`)
    await loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('卸载失败: ' + (e?.message || e))
    }
  }
}

const handleUse = async (version: string) => {
  try {
    await UseVersion(version)
    ElMessage.success(`已切换到 ${version}`)
    await loadData()
  } catch (e: any) {
    ElMessage.error('切换失败: ' + (e?.message || e))
  }
}

// 打开设置对话框
const openSettingsDialog = () => {
  tempSettings.value = { ...installSettings.value }
  settingsDialog.value = true
}

// 选择目录
const selectDir = async (type: 'install' | 'cache') => {
  try {
    const title = type === 'install' ? '选择安装目录' : '选择下载缓存目录'
    const dir = await SelectDirectory(title)
    if (dir) {
      if (type === 'install') {
        tempSettings.value.installDir = dir
      } else {
        tempSettings.value.cacheDir = dir
      }
    }
  } catch (e) {
    // 用户取消选择
  }
}

// 保存设置
const saveSettings = async () => {
  try {
    await SetInstallSettings(tempSettings.value.installDir, tempSettings.value.cacheDir)
    installSettings.value = { ...tempSettings.value }
    settingsDialog.value = false
    ElMessage.success('设置已保存')
  } catch (e: any) {
    ElMessage.error('保存失败: ' + (e?.message || e))
  }
}

onMounted(() => {
  loadData()

  // 监听下载进度
  EventsOn('download-progress', (progress: DownloadProgress) => {
    downloadProgress.value = Math.round(progress.percent)
  })
})
</script>

<template>
  <div class="version-manager">
    <div class="section">
      <div class="section-header">
        <h3>已安装版本</h3>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="openSettingsDialog">
            <el-icon><Setting /></el-icon>
            路径设置
          </el-button>
          <el-button :icon="'Refresh'" circle size="small" @click="loadData" :loading="loading" />
        </div>
      </div>

      <!-- 当前路径信息 -->
      <div class="path-info" v-if="installSettings.installDir || installSettings.cacheDir">
        <el-tag type="info" v-if="installSettings.installDir">
          安装目录: {{ installSettings.installDir }}
        </el-tag>
        <el-tag type="info" v-if="installSettings.cacheDir">
          缓存目录: {{ installSettings.cacheDir }}
        </el-tag>
      </div>

      <el-table :data="installedVersions" v-loading="loading" empty-text="暂无已安装版本">
        <el-table-column prop="version" label="版本" width="180">
          <template #default="{ row }">
            <span>{{ row.version }}</span>
            <el-tag v-if="row.isCurrent" size="small" type="success" style="margin-left: 8px">当前</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button
              v-if="!row.isCurrent"
              type="primary"
              size="small"
              @click="handleUse(row.version)"
            >
              使用
            </el-button>
            <el-button
              v-if="!row.isCurrent"
              type="danger"
              size="small"
              @click="handleUninstall(row.version)"
            >
              卸载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="section">
      <div class="section-header">
        <h3>可用版本</h3>
        <el-checkbox v-model="showStableOnly" @change="loadData">只显示稳定版</el-checkbox>
      </div>

      <el-table :data="filteredRemoteVersions.slice(0, 20)" v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="version" label="版本" width="180">
          <template #default="{ row }">
            <span>{{ row.version }}</span>
            <el-tag v-if="row.stable" size="small" style="margin-left: 8px">稳定版</el-tag>
            <el-tag v-else size="small" type="warning" style="margin-left: 8px">预览版</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="isInstalled(row.version)" type="success">已安装</el-tag>
            <el-tag v-else type="info">未安装</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <template v-if="installing && installingVersion === row.version">
              <el-progress :percentage="downloadProgress" :width="100" />
            </template>
            <el-button
              v-else-if="!isInstalled(row.version)"
              type="primary"
              size="small"
              @click="openInstallDialog(row.version)"
              :disabled="installing"
            >
              安装
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 路径设置对话框 -->
    <el-dialog v-model="settingsDialog" title="路径设置" width="550px">
      <el-form label-width="120px">
        <el-form-item label="安装目录">
          <div class="dir-input">
            <el-input v-model="tempSettings.installDir" placeholder="默认: ~/.govm/versions" />
            <el-button @click="selectDir('install')">浏览</el-button>
          </div>
          <div class="form-tip">Go 版本将安装到此目录</div>
        </el-form-item>
        <el-form-item label="下载缓存目录">
          <div class="dir-input">
            <el-input v-model="tempSettings.cacheDir" placeholder="默认: ~/.govm/cache" />
            <el-button @click="selectDir('cache')">浏览</el-button>
          </div>
          <div class="form-tip">下载的安装包将缓存到此目录</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settingsDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSettings">保存</el-button>
      </template>
    </el-dialog>

    <!-- 安装确认对话框 -->
    <el-dialog v-model="installDialog" title="安装确认" width="500px">
      <el-form label-width="100px">
        <el-form-item label="版本">
          <el-tag>{{ selectedVersion }}</el-tag>
        </el-form-item>
        <el-form-item label="安装目录">
          <div class="dir-input">
            <el-input v-model="customInstallDir" placeholder="使用默认目录" />
            <el-button @click="selectInstallDir">浏览</el-button>
          </div>
          <div class="form-tip">留空则使用默认安装目录</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="installDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmInstall">开始安装</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.version-manager {
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
  align-items: center;
}

.path-info {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.dir-input {
  display: flex;
  gap: 8px;
  width: 100%;
}

.dir-input .el-input {
  flex: 1;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>

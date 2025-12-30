<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  GetGoPathMode,
  SetGoPathMode,
  GetSharedGoPath,
  SetSharedGoPath,
  GetModuleCacheStats,
  GetModules,
  CleanModuleCache,
  CleanModule,
  VerifyModules,
  GetCurrentGoPath,
  GetModuleCachePath,
  SelectDirectory,
  SearchPackages,
  GetPackage
} from '../../wailsjs/go/main/App'

interface ModuleInfo {
  path: string
  version: string
  size: number
  dir: string
  name: string
  description: string
  category: string
}

interface ModuleCacheStats {
  totalModules: number
  totalSize: number
  totalSizeStr: string
  cachePath: string
}

interface VerifyResult {
  module: string
  version: string
  status: string
  message: string
}

interface SearchResult {
  path: string
  version: string
  description: string
}

const loading = ref(false)
const goPathMode = ref('shared')
const sharedGoPath = ref('')
const currentGoPath = ref('')
const moduleCachePath = ref('')
const cacheStats = ref<ModuleCacheStats | null>(null)
const modules = ref<ModuleInfo[]>([])
const searchQuery = ref('')
const verifyResults = ref<VerifyResult[]>([])
const verifying = ref(false)

// 分页
const currentPage = ref(1)
const pageSize = ref(15)
const categoryFilter = ref('thirdparty') // 默认只显示第三方包

// 设置对话框
const settingsDialog = ref(false)
const tempMode = ref('shared')
const tempSharedPath = ref('')

// 搜索安装
const installDialog = ref(false)
const installSearchQuery = ref('')
const searchResults = ref<SearchResult[]>([])
const searching = ref(false)
const installing = ref('')

// 过滤后的模块列表
const filteredModules = computed(() => {
  let result = modules.value

  // 按分类筛选
  if (categoryFilter.value !== 'all') {
    result = result.filter(m => m.category === categoryFilter.value)
  }

  // 按搜索词筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(m => 
      m.path.toLowerCase().includes(query) || 
      m.version.toLowerCase().includes(query) ||
      m.name.toLowerCase().includes(query) ||
      (m.description && m.description.toLowerCase().includes(query))
    )
  }

  return result
})

// 统计数量
const thirdpartyCount = computed(() => modules.value.filter(m => m.category === 'thirdparty').length)
const officialCount = computed(() => modules.value.filter(m => m.category === 'official').length)

// 分页后的模块列表
const pagedModules = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredModules.value.slice(start, end)
})

// 总数
const totalModules = computed(() => filteredModules.value.length)

// 搜索时重置页码
const handleSearch = () => {
  currentPage.value = 1
}

// 切换分类时重置页码
const handleCategoryChange = () => {
  currentPage.value = 1
}

// 格式化大小
const formatSize = (size: number) => {
  const KB = 1024
  const MB = KB * 1024
  const GB = MB * 1024
  if (size >= GB) return (size / GB).toFixed(2) + ' GB'
  if (size >= MB) return (size / MB).toFixed(2) + ' MB'
  if (size >= KB) return (size / KB).toFixed(2) + ' KB'
  return size + ' B'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const [mode, sharedPath, gopath, cachePath, stats, moduleList] = await Promise.all([
      GetGoPathMode(),
      GetSharedGoPath(),
      GetCurrentGoPath(),
      GetModuleCachePath(),
      GetModuleCacheStats(),
      GetModules()
    ])
    goPathMode.value = mode || 'shared'
    sharedGoPath.value = sharedPath || ''
    currentGoPath.value = gopath
    moduleCachePath.value = cachePath
    cacheStats.value = stats
    modules.value = moduleList || []
  } catch (e: any) {
    ElMessage.error('加载数据失败: ' + (e?.message || e))
  } finally {
    loading.value = false
  }
}

// 打开设置对话框
const openSettings = () => {
  tempMode.value = goPathMode.value
  tempSharedPath.value = sharedGoPath.value
  settingsDialog.value = true
}

// 选择目录
const selectSharedPath = async () => {
  try {
    const dir = await SelectDirectory('选择共享 GOPATH 目录')
    if (dir) {
      tempSharedPath.value = dir
    }
  } catch (e) {}
}

// 保存设置
const saveSettings = async () => {
  try {
    await SetGoPathMode(tempMode.value)
    if (tempMode.value === 'shared' && tempSharedPath.value) {
      await SetSharedGoPath(tempSharedPath.value)
    }
    goPathMode.value = tempMode.value
    sharedGoPath.value = tempSharedPath.value
    settingsDialog.value = false
    ElMessage.success('设置已保存')
    await loadData()
  } catch (e: any) {
    ElMessage.error('保存失败: ' + (e?.message || e))
  }
}

// 清理所有模块缓存
const handleCleanAll = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理所有模块缓存吗？这将删除所有已下载的第三方包。',
      '确认清理',
      { type: 'warning' }
    )
    await CleanModuleCache()
    ElMessage.success('模块缓存已清理')
    await loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('清理失败: ' + (e?.message || e))
    }
  }
}

// 删除单个模块
const handleDeleteModule = async (mod: ModuleInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除模块 ${mod.path}@${mod.version} 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await CleanModule(mod.path, mod.version)
    ElMessage.success('模块已删除')
    await loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败: ' + (e?.message || e))
    }
  }
}

// 校验模块
const handleVerify = async () => {
  verifying.value = true
  try {
    const results = await VerifyModules()
    verifyResults.value = results
    if (results.length === 1 && results[0].status === 'ok') {
      ElMessage.success('所有模块校验通过')
    } else {
      ElMessage.warning('部分模块校验失败，请查看详情')
    }
  } catch (e: any) {
    ElMessage.error('校验失败: ' + (e?.message || e))
  } finally {
    verifying.value = false
  }
}

// 打开安装对话框
const openInstallDialog = () => {
  installSearchQuery.value = ''
  searchResults.value = []
  installDialog.value = true
}

// 搜索包
const handleSearchPackages = async () => {
  if (!installSearchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  searching.value = true
  try {
    const results = await SearchPackages(installSearchQuery.value.trim())
    searchResults.value = results || []
  } catch (e: any) {
    ElMessage.error('搜索失败: ' + (e?.message || e))
  } finally {
    searching.value = false
  }
}

// 安装包
const handleInstallPackage = async (pkg: SearchResult) => {
  installing.value = pkg.path
  try {
    await GetPackage(pkg.path, pkg.version || 'latest')
    ElMessage.success(`${pkg.path} 下载成功`)
    await loadData()
  } catch (e: any) {
    ElMessage.error('下载失败: ' + (e?.message || e))
  } finally {
    installing.value = ''
  }
}

// 直接安装输入的包
const handleDirectInstall = async () => {
  if (!installSearchQuery.value.trim()) return
  const pkgPath = installSearchQuery.value.trim()
  installing.value = pkgPath
  try {
    await GetPackage(pkgPath, 'latest')
    ElMessage.success(`${pkgPath} 下载成功`)
    installDialog.value = false
    await loadData()
  } catch (e: any) {
    ElMessage.error('下载失败: ' + (e?.message || e))
  } finally {
    installing.value = ''
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="package-manager">
    <!-- GOPATH 模式设置 -->
    <div class="section">
      <div class="section-header">
        <h3>GOPATH 模式</h3>
        <el-button type="primary" size="small" @click="openSettings">
          <el-icon><Setting /></el-icon>
          设置
        </el-button>
      </div>
      
      <el-descriptions :column="1" border>
        <el-descriptions-item label="当前模式">
          <el-tag :type="goPathMode === 'shared' ? 'success' : 'warning'">
            {{ goPathMode === 'shared' ? '共享模式' : '隔离模式' }}
          </el-tag>
          <span class="mode-desc">
            {{ goPathMode === 'shared' 
              ? '所有 Go 版本共用一个 GOPATH，节省磁盘空间' 
              : '每个 Go 版本独立 GOPATH，完全隔离' }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="GOPATH">
          <span class="path-text">{{ currentGoPath || '未设置' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="模块缓存">
          <span class="path-text">{{ moduleCachePath || '未设置' }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 缓存统计 -->
    <div class="section">
      <div class="section-header">
        <h3>模块缓存统计</h3>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="openInstallDialog">
            <el-icon><Plus /></el-icon>
            安装包
          </el-button>
          <el-button size="small" @click="handleVerify" :loading="verifying">
            <el-icon><Check /></el-icon>
            校验模块
          </el-button>
          <el-button type="danger" size="small" @click="handleCleanAll">
            <el-icon><Delete /></el-icon>
            清理全部
          </el-button>
          <el-button :icon="'Refresh'" circle size="small" @click="loadData" :loading="loading" />
        </div>
      </div>

      <div class="stats-cards" v-if="cacheStats">
        <div class="stat-card">
          <div class="stat-value">{{ cacheStats.totalModules }}</div>
          <div class="stat-label">模块数量</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ cacheStats.totalSizeStr }}</div>
          <div class="stat-label">占用空间</div>
        </div>
      </div>
    </div>

    <!-- 模块列表 -->
    <div class="section">
      <div class="section-header">
        <h3>已下载模块</h3>
        <div class="filter-bar">
          <el-radio-group v-model="categoryFilter" size="small" @change="handleCategoryChange">
            <el-radio-button value="thirdparty">第三方包 ({{ thirdpartyCount }})</el-radio-button>
            <el-radio-button value="official">官方包 ({{ officialCount }})</el-radio-button>
            <el-radio-button value="all">全部 ({{ modules.length }})</el-radio-button>
          </el-radio-group>
          <el-input
            v-model="searchQuery"
            placeholder="搜索..."
            style="width: 200px"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>

      <el-table 
        :data="pagedModules" 
        v-loading="loading" 
        empty-text="暂无模块"
      >
        <el-table-column label="包信息" min-width="350">
          <template #default="{ row }">
            <div class="module-info">
              <div class="module-name">{{ row.name || row.path.split('/').pop() }}</div>
              <div class="module-path">{{ row.path }}</div>
              <div class="module-desc" v-if="row.description">{{ row.description }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="140">
          <template #default="{ row }">
            <el-tag size="small">{{ row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button
              type="danger"
              size="small"
              link
              @click="handleDeleteModule(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 15, 20, 50]"
          :total="totalModules"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </div>

    <!-- 设置对话框 -->
    <el-dialog v-model="settingsDialog" title="GOPATH 模式设置" width="500px">
      <div class="mode-options">
        <div 
          class="mode-option" 
          :class="{ active: tempMode === 'shared' }"
          @click="tempMode = 'shared'"
        >
          <div class="mode-option-header">
            <el-radio v-model="tempMode" value="shared">共享模式</el-radio>
          </div>
          <div class="mode-option-desc">所有 Go 版本共用一个 GOPATH，节省磁盘空间</div>
        </div>
        
        <div 
          class="mode-option" 
          :class="{ active: tempMode === 'isolated' }"
          @click="tempMode = 'isolated'"
        >
          <div class="mode-option-header">
            <el-radio v-model="tempMode" value="isolated">隔离模式</el-radio>
          </div>
          <div class="mode-option-desc">每个 Go 版本独立 GOPATH，完全隔离避免兼容问题</div>
        </div>
      </div>

      <div class="shared-path-setting" v-if="tempMode === 'shared'">
        <div class="setting-label">共享路径（可选）</div>
        <div class="dir-input">
          <el-input v-model="tempSharedPath" placeholder="留空使用系统默认 GOPATH" />
          <el-button @click="selectSharedPath">浏览</el-button>
        </div>
        <div class="form-tip">自定义共享 GOPATH 路径，留空则使用系统默认</div>
      </div>

      <template #footer>
        <el-button @click="settingsDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSettings">保存</el-button>
      </template>
    </el-dialog>

    <!-- 安装包对话框 -->
    <el-dialog v-model="installDialog" title="搜索安装包" width="600px">
      <div class="install-search">
        <el-input
          v-model="installSearchQuery"
          placeholder="输入包名或关键词，如 gin、redis、gorm..."
          size="large"
          @keyup.enter="handleSearchPackages"
        >
          <template #append>
            <el-button @click="handleSearchPackages" :loading="searching">搜索</el-button>
          </template>
        </el-input>
        <div class="search-tip">
          输入完整包路径可直接下载，如 github.com/gin-gonic/gin
        </div>
      </div>

      <div class="search-results" v-if="searchResults.length > 0">
        <div 
          class="search-result-item" 
          v-for="pkg in searchResults" 
          :key="pkg.path"
        >
          <div class="pkg-info">
            <div class="pkg-path">{{ pkg.path }}</div>
            <div class="pkg-desc" v-if="pkg.description">{{ pkg.description }}</div>
          </div>
          <el-button 
            type="primary" 
            size="small"
            :loading="installing === pkg.path"
            @click="handleInstallPackage(pkg)"
          >
            下载
          </el-button>
        </div>
      </div>

      <div class="no-results" v-else-if="installSearchQuery && !searching">
        <p>没有找到匹配的包建议</p>
        <p class="hint">你可以直接输入完整的包路径进行下载</p>
      </div>

      <template #footer>
        <el-button @click="installDialog = false">关闭</el-button>
        <el-button 
          type="primary" 
          @click="handleDirectInstall"
          :loading="installing === installSearchQuery"
          :disabled="!installSearchQuery.trim()"
        >
          直接下载
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.package-manager {
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

.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
}

.mode-desc {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}

.path-text {
  font-family: monospace;
  font-size: 13px;
  color: #606266;
}

.stats-cards {
  display: flex;
  gap: 20px;
}

.stat-card {
  flex: 1;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #1890ff;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.module-info {
  padding: 4px 0;
}

.module-name {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
}

.module-path {
  font-family: monospace;
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.module-desc {
  font-size: 12px;
  color: #606266;
  margin-top: 4px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
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

.mode-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.mode-option {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.mode-option:hover {
  border-color: #1890ff;
}

.mode-option.active {
  border-color: #1890ff;
  background: rgba(24, 144, 255, 0.05);
}

.mode-option-header {
  margin-bottom: 6px;
}

.mode-option-desc {
  font-size: 13px;
  color: #909399;
  margin-left: 24px;
}

.shared-path-setting {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e8e8e8;
}

.setting-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

/* 安装包对话框 */
.install-search {
  margin-bottom: 20px;
}

.search-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.search-results {
  max-height: 350px;
  overflow-y: auto;
}

.search-result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  margin-bottom: 8px;
}

.search-result-item:hover {
  border-color: #1890ff;
  background: rgba(24, 144, 255, 0.02);
}

.pkg-info {
  flex: 1;
  min-width: 0;
}

.pkg-path {
  font-family: monospace;
  font-size: 14px;
  color: #303133;
  word-break: break-all;
}

.pkg-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.no-results {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.no-results p {
  margin: 0;
}

.no-results .hint {
  font-size: 12px;
  margin-top: 8px;
}
</style>

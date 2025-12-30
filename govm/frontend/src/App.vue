<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { GetCurrentVersion, GetAppVersion, GetConfig } from '../wailsjs/go/main/App'

const router = useRouter()
const currentVersion = ref('')
const appVersion = ref('')
const activeMenu = ref('/version')
const theme = ref('light')

const isDark = computed(() => theme.value === 'dark')

const menuItems = [
  { path: '/version', icon: 'Box', title: '版本管理' },
  { path: '/package', icon: 'Files', title: '包管理' },
  { path: '/env', icon: 'Monitor', title: '环境配置' },
  { path: '/cache', icon: 'Delete', title: '缓存清理' },
  { path: '/settings', icon: 'Setting', title: '设置' }
]

const handleMenuSelect = (path: string) => {
  router.push(path)
}

const loadCurrentVersion = async () => {
  try {
    currentVersion.value = await GetCurrentVersion()
  } catch (e) {
    currentVersion.value = '未设置'
  }
}

onMounted(async () => {
  await loadCurrentVersion()
  appVersion.value = await GetAppVersion()
  activeMenu.value = router.currentRoute.value.path

  // 加载主题配置
  try {
    const config = await GetConfig()
    theme.value = config.theme || 'light'
  } catch (e) {}

  // 监听版本切换事件
  EventsOn('version-changed', (version: string) => {
    currentVersion.value = version
  })

  // 监听主题切换事件
  EventsOn('theme-changed', (newTheme: string) => {
    theme.value = newTheme
  })
})
</script>

<template>
  <el-container class="app-container" :class="{ 'theme-dark': isDark }">
    <el-aside width="180px" class="app-aside">
      <div class="logo">
        <el-icon size="24"><Box /></el-icon>
        <span>GoVM</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="app-menu"
        @select="handleMenuSelect"
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-main class="app-main">
        <router-view :theme="theme" />
      </el-main>

      <el-footer height="32px" class="app-footer">
        <span>当前版本: {{ currentVersion || '未设置' }}</span>
        <span>GoVM v{{ appVersion }}</span>
      </el-footer>
    </el-container>
  </el-container>
</template>

<style scoped>
/* ========== 浅色主题（默认） ========== */
.app-container {
  height: 100vh;
  background: #f0f2f5;
}

.app-aside {
  background: #fff;
  border-right: 1px solid #e8e8e8;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  font-size: 18px;
  font-weight: bold;
  color: #1890ff;
}

.app-menu {
  border-right: none;
  background: transparent;
}

.app-main {
  background: #f0f2f5;
  padding: 20px;
}

.app-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: #fff;
  border-top: 1px solid #e8e8e8;
  font-size: 12px;
  color: #666;
}

/* ========== 深色主题 ========== */
.app-container.theme-dark {
  background: #1f1f1f;
}

.app-container.theme-dark .app-aside {
  background: #141414;
  border-color: #303030;
}

.app-container.theme-dark .logo {
  color: #1890ff;
}

.app-container.theme-dark .app-main {
  background: #1f1f1f;
}

.app-container.theme-dark .app-footer {
  background: #141414;
  border-color: #303030;
  color: #8c8c8c;
}

/* 深色菜单 */
.app-container.theme-dark :deep(.el-menu) {
  background: transparent;
}

.app-container.theme-dark :deep(.el-menu-item) {
  color: #a6a6a6;
}

.app-container.theme-dark :deep(.el-menu-item:hover) {
  background: #262626;
  color: #fff;
}

.app-container.theme-dark :deep(.el-menu-item.is-active) {
  color: #1890ff;
  background: #262626;
}
</style>

<style>
/* ========== 深色主题全局样式 ========== */
.theme-dark .section {
  background: #262626 !important;
  border-color: #303030 !important;
}

.theme-dark .section h3,
.theme-dark .section-header h3 {
  color: #e8e8e8 !important;
}

/* 表格 */
.theme-dark .el-table {
  --el-table-bg-color: #262626;
  --el-table-tr-bg-color: #262626;
  --el-table-header-bg-color: #1f1f1f;
  --el-table-row-hover-bg-color: #303030;
  --el-table-border-color: #303030;
  --el-table-text-color: #e8e8e8;
  --el-table-header-text-color: #a6a6a6;
}

.theme-dark .el-table th.el-table__cell {
  background: #1f1f1f !important;
}

.theme-dark .el-table td.el-table__cell {
  background: #262626 !important;
}

.theme-dark .el-table__row:hover > td.el-table__cell {
  background: #303030 !important;
}

.theme-dark .el-table__empty-text {
  color: #8c8c8c;
}

/* 输入框 */
.theme-dark .el-input__wrapper {
  background: #1f1f1f !important;
  box-shadow: 0 0 0 1px #303030 inset !important;
}

.theme-dark .el-input__inner {
  color: #e8e8e8 !important;
}

.theme-dark .el-input__inner::placeholder {
  color: #6b6b6b !important;
}

/* 选择器 */
.theme-dark .el-select__wrapper {
  background: #1f1f1f !important;
  box-shadow: 0 0 0 1px #303030 inset !important;
}

.theme-dark .el-select__placeholder {
  color: #e8e8e8 !important;
}

.theme-dark .el-select-dropdown {
  background: #262626 !important;
  border-color: #303030 !important;
}

.theme-dark .el-select-dropdown__item {
  color: #e8e8e8 !important;
}

.theme-dark .el-select-dropdown__item:hover {
  background: #303030 !important;
}

/* 描述列表 */
.theme-dark .el-descriptions__label {
  background: #1f1f1f !important;
  color: #a6a6a6 !important;
}

.theme-dark .el-descriptions__content {
  background: #262626 !important;
  color: #e8e8e8 !important;
}

.theme-dark .el-descriptions--bordered .el-descriptions__cell {
  border-color: #303030 !important;
}

/* 卡片 */
.theme-dark .el-card {
  background: #262626 !important;
  border-color: #303030 !important;
}

.theme-dark .el-card__header {
  background: #1f1f1f !important;
  border-color: #303030 !important;
  color: #e8e8e8 !important;
}

/* 标签 */
.theme-dark .el-tag {
  background: #303030 !important;
  border-color: #404040 !important;
  color: #e8e8e8 !important;
}

.theme-dark .el-tag--info {
  background: #303030 !important;
  color: #a6a6a6 !important;
}

.theme-dark .el-tag--success {
  background: #274916 !important;
  border-color: #306317 !important;
  color: #95de64 !important;
}

/* 表单 */
.theme-dark .el-form-item__label {
  color: #a6a6a6 !important;
}

/* 单选/复选 */
.theme-dark .el-radio__label,
.theme-dark .el-checkbox__label {
  color: #e8e8e8 !important;
}

.theme-dark .el-radio__inner,
.theme-dark .el-checkbox__inner {
  background: #1f1f1f !important;
  border-color: #404040 !important;
}

/* 按钮 */
.theme-dark .el-button--default {
  background: #303030 !important;
  border-color: #404040 !important;
  color: #e8e8e8 !important;
}

.theme-dark .el-button--default:hover {
  background: #404040 !important;
  color: #fff !important;
}

/* 对话框 */
.theme-dark .el-dialog {
  background: #262626 !important;
}

.theme-dark .el-dialog__title {
  color: #e8e8e8 !important;
}

/* 空状态 */
.theme-dark .el-empty__description {
  color: #8c8c8c !important;
}

/* 其他 */
.theme-dark .form-tip,
.theme-dark .not-set {
  color: #8c8c8c !important;
}

.theme-dark .cache-size {
  color: #1890ff !important;
}

.theme-dark .clean-card p,
.theme-dark .card-header {
  color: #a6a6a6 !important;
}
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { GetCacheInfo, CleanDownloadCache, CleanAllCache } from '../../wailsjs/go/main/App'

interface CacheInfo {
  downloadCacheSize: number
  downloadCachePath: string
  totalSize: number
  totalSizeHuman: string
}

const cacheInfo = ref<CacheInfo>({
  downloadCacheSize: 0,
  downloadCachePath: '',
  totalSize: 0,
  totalSizeHuman: '0 B'
})

const loading = ref(false)
const cleaning = ref(false)

const loadCacheInfo = async () => {
  loading.value = true
  try {
    cacheInfo.value = await GetCacheInfo()
  } catch (e: any) {
    ElMessage.error('获取缓存信息失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

const handleCleanDownload = async () => {
  try {
    await ElMessageBox.confirm('确定要清理下载缓存吗？', '确认清理', {
      type: 'warning'
    })
    cleaning.value = true
    await CleanDownloadCache()
    ElMessage.success('下载缓存已清理')
    await loadCacheInfo()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('清理失败: ' + e.message)
    }
  } finally {
    cleaning.value = false
  }
}

const handleCleanAll = async () => {
  try {
    await ElMessageBox.confirm('确定要清理所有缓存吗？', '确认清理', {
      type: 'warning'
    })
    cleaning.value = true
    await CleanAllCache()
    ElMessage.success('所有缓存已清理')
    await loadCacheInfo()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('清理失败: ' + e.message)
    }
  } finally {
    cleaning.value = false
  }
}

onMounted(() => {
  loadCacheInfo()
})
</script>

<template>
  <div class="cache-manager">
    <div class="section">
      <div class="section-header">
        <h3>缓存信息</h3>
        <el-button :icon="'Refresh'" circle size="small" @click="loadCacheInfo" :loading="loading" />
      </div>

      <el-descriptions :column="1" border v-loading="loading">
        <el-descriptions-item label="缓存总大小">
          <span class="cache-size">{{ cacheInfo.totalSizeHuman }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="下载缓存路径">
          {{ cacheInfo.downloadCachePath }}
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="section">
      <div class="section-header">
        <h3>清理操作</h3>
      </div>

      <div class="clean-actions">
        <el-card shadow="hover" class="clean-card">
          <template #header>
            <div class="card-header">
              <el-icon size="24"><Delete /></el-icon>
              <span>清理下载缓存</span>
            </div>
          </template>
          <p>清理已下载的 Go 安装包，不影响已安装的版本</p>
          <el-button type="primary" @click="handleCleanDownload" :loading="cleaning">
            清理下载缓存
          </el-button>
        </el-card>

        <el-card shadow="hover" class="clean-card">
          <template #header>
            <div class="card-header">
              <el-icon size="24"><DeleteFilled /></el-icon>
              <span>清理所有缓存</span>
            </div>
          </template>
          <p>清理所有缓存文件，释放磁盘空间</p>
          <el-button type="danger" @click="handleCleanAll" :loading="cleaning">
            清理全部
          </el-button>
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cache-manager {
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

.cache-size {
  font-size: 18px;
  font-weight: bold;
  color: #409eff;
}

.clean-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.clean-card {
  text-align: center;
}

.clean-card p {
  color: #909399;
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
</style>

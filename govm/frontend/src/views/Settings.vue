<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  GetConfig,
  GetMirrorOptions,
  SetMirror,
  SetProxy,
  SetGoProxy,
  SetTheme,
  ResetConfig,
  GetSystemInfo,
  GetAppVersion
} from '../../wailsjs/go/main/App'

interface Config {
  mirror: string
  proxy: string
  goproxy: string
  theme: string
  language: string
}

interface MirrorOption {
  name: string
  url: string
}

const config = ref<Config>({
  mirror: '',
  proxy: '',
  goproxy: '',
  theme: 'light',
  language: 'zh-CN'
})

const mirrorOptions = ref<MirrorOption[]>([])
const systemInfo = ref<Record<string, string>>({})
const appVersion = ref('')
const loading = ref(false)
const saving = ref(false)

const loadData = async () => {
  loading.value = true
  try {
    const [cfg, mirrors, sysInfo, version] = await Promise.all([
      GetConfig(),
      GetMirrorOptions(),
      GetSystemInfo(),
      GetAppVersion()
    ])
    config.value = cfg
    mirrorOptions.value = mirrors
    systemInfo.value = sysInfo
    appVersion.value = version
  } catch (e: any) {
    ElMessage.error('加载配置失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

const handleMirrorChange = async (value: string) => {
  saving.value = true
  try {
    await SetMirror(value)
    ElMessage.success('镜像源已更新')
  } catch (e: any) {
    ElMessage.error('保存失败: ' + e.message)
  } finally {
    saving.value = false
  }
}

const handleProxyChange = async () => {
  saving.value = true
  try {
    await SetProxy(config.value.proxy)
    ElMessage.success('代理设置已更新')
  } catch (e: any) {
    ElMessage.error('保存失败: ' + e.message)
  } finally {
    saving.value = false
  }
}

const handleGoProxyChange = async () => {
  saving.value = true
  try {
    await SetGoProxy(config.value.goproxy)
    ElMessage.success('GOPROXY 已更新')
  } catch (e: any) {
    ElMessage.error('保存失败: ' + e.message)
  } finally {
    saving.value = false
  }
}

const handleThemeChange = async (value: string) => {
  saving.value = true
  try {
    await SetTheme(value)
    ElMessage.success('主题已更新')
  } catch (e: any) {
    ElMessage.error('保存失败: ' + e.message)
  } finally {
    saving.value = false
  }
}

const handleReset = async () => {
  try {
    await ElMessageBox.confirm('确定要重置所有配置吗？', '确认重置', {
      type: 'warning'
    })
    await ResetConfig()
    ElMessage.success('配置已重置')
    await loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('重置失败: ' + e.message)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="settings" v-loading="loading">
    <div class="section">
      <h3>下载设置</h3>

      <el-form label-width="120px">
        <el-form-item label="镜像源">
          <el-select v-model="config.mirror" @change="handleMirrorChange" style="width: 100%">
            <el-option
              v-for="option in mirrorOptions"
              :key="option.url"
              :label="option.name"
              :value="option.url"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="HTTP 代理">
          <el-input
            v-model="config.proxy"
            placeholder="例如: http://127.0.0.1:7890"
            @blur="handleProxyChange"
          />
        </el-form-item>

        <el-form-item label="GOPROXY">
          <el-input
            v-model="config.goproxy"
            placeholder="例如: https://goproxy.cn,direct"
            @blur="handleGoProxyChange"
          />
        </el-form-item>
      </el-form>
    </div>

    <div class="section">
      <h3>外观设置</h3>

      <el-form label-width="120px">
        <el-form-item label="主题">
          <el-radio-group v-model="config.theme" @change="handleThemeChange">
            <el-radio value="light">浅色</el-radio>
            <el-radio value="dark">深色</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </div>

    <div class="section">
      <h3>系统信息</h3>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="操作系统">{{ systemInfo.os }}</el-descriptions-item>
        <el-descriptions-item label="架构">{{ systemInfo.arch }}</el-descriptions-item>
        <el-descriptions-item label="数据目录">{{ systemInfo.baseDir }}</el-descriptions-item>
        <el-descriptions-item label="GoVM 版本">v{{ appVersion }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="section">
      <h3>其他操作</h3>

      <el-button type="warning" @click="handleReset">
        恢复默认配置
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.settings {
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

.section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #303133;
}
</style>

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"govm/internal/cache"
	"govm/internal/config"
	"govm/internal/env"
	"govm/internal/module"
	"govm/internal/version"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	baseDir        string
	versionManager *version.Manager
	configManager  *config.Manager
	envManager     *env.Manager
	cacheManager   *cache.Manager
	moduleManager  *module.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, _ := os.UserHomeDir()
	baseDir := filepath.Join(homeDir, ".govm")

	return &App{
		baseDir: baseDir,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化管理器
	a.configManager = config.NewManager(a.baseDir)
	a.configManager.Load()

	cfg := a.configManager.Get()
	a.versionManager = version.NewManager(a.baseDir, cfg.Mirror)
	// 设置自定义路径
	a.versionManager.SetPaths(cfg.InstallDir, cfg.CacheDir)

	a.envManager = env.NewManager(a.baseDir)
	a.cacheManager = cache.NewManager(a.baseDir)
	// 设置缓存路径
	a.cacheManager.SetCacheDir(cfg.CacheDir)

	// 初始化模块管理器
	a.moduleManager = module.NewManager(a.baseDir)
	a.moduleManager.SetGoPathMode(cfg.GoPathMode, cfg.SharedGoPath)

	// 确保目录存在
	os.MkdirAll(filepath.Join(a.baseDir, "versions"), 0755)
	os.MkdirAll(filepath.Join(a.baseDir, "cache"), 0755)
}

// ==================== 版本管理 ====================

// InstallSettings 安装设置
type InstallSettings struct {
	InstallDir string `json:"installDir"`
	CacheDir   string `json:"cacheDir"`
}

// GetRemoteVersions 获取远程版本列表
func (a *App) GetRemoteVersions(includeAll bool) ([]version.GoVersion, error) {
	return a.versionManager.GetRemoteVersions(includeAll)
}

// GetInstalledVersions 获取已安装版本列表
func (a *App) GetInstalledVersions() ([]version.InstalledVersion, error) {
	return a.versionManager.GetInstalledVersions()
}

// GetCurrentVersion 获取当前版本
func (a *App) GetCurrentVersion() string {
	return a.versionManager.GetCurrentVersion()
}

// InstallVersion 安装指定版本
func (a *App) InstallVersion(ver string, installDir string) error {
	return a.versionManager.Install(ver, installDir, func(progress version.DownloadProgress) {
		wailsRuntime.EventsEmit(a.ctx, "download-progress", progress)
	})
}

// UninstallVersion 卸载指定版本
func (a *App) UninstallVersion(ver string) error {
	return a.versionManager.Uninstall(ver)
}

// UseVersion 切换到指定版本
func (a *App) UseVersion(ver string) error {
	// 获取版本路径（可能是 GoVM 管理的，也可能是系统已安装的）
	goroot := ""
	
	// 先检查是否是已安装版本列表中的
	installed, _ := a.versionManager.GetInstalledVersions()
	for _, v := range installed {
		// 支持多种匹配方式：精确匹配、路径匹配
		if v.Version == ver || v.Path == ver {
			goroot = v.Path
			break
		}
	}
	
	// 如果没找到，尝试直接用传入的值作为路径
	if goroot == "" {
		for _, v := range installed {
			// 检查路径是否存在
			if strings.HasSuffix(v.Path, ver) || strings.Contains(v.Path, ver) {
				goroot = v.Path
				break
			}
		}
	}
	
	if goroot == "" {
		// 返回更详细的错误信息
		var versions []string
		for _, v := range installed {
			versions = append(versions, fmt.Sprintf("%s(%s)", v.Version, v.Path))
		}
		return fmt.Errorf("版本 %s 未找到，已安装: %v", ver, versions)
	}

	// 记录当前版本
	if err := a.versionManager.Use(ver); err != nil {
		return err
	}

	// 更新 GOROOT
	if err := a.envManager.SetGOROOT(goroot); err != nil {
		return err
	}

	// 更新 PATH（添加 Go bin 目录）
	if err := a.envManager.UpdatePATH(goroot); err != nil {
		return err
	}

	// 根据 GOPATH 模式更新 GOPATH
	cfg := a.configManager.Get()
	if cfg.GoPathMode == "isolated" {
		// 隔离模式：设置版本专属的 GOPATH
		gopath := a.moduleManager.GetGoPath(ver)
		os.MkdirAll(gopath, 0755)
		if err := a.envManager.SetEnvVar("GOPATH", gopath); err != nil {
			return err
		}
	} else if cfg.SharedGoPath != "" {
		// 共享模式且设置了自定义路径
		if err := a.envManager.SetEnvVar("GOPATH", cfg.SharedGoPath); err != nil {
			return err
		}
	}

	// 通知前端刷新
	wailsRuntime.EventsEmit(a.ctx, "version-changed", ver)

	return nil
}

// GetLatestStableVersion 获取最新稳定版
func (a *App) GetLatestStableVersion() (string, error) {
	return a.versionManager.GetLatestStableVersion()
}

// GetInstallSettings 获取安装设置
func (a *App) GetInstallSettings() InstallSettings {
	cfg := a.configManager.Get()
	return InstallSettings{
		InstallDir: cfg.InstallDir,
		CacheDir:   cfg.CacheDir,
	}
}

// SetInstallSettings 设置安装路径
func (a *App) SetInstallSettings(installDir, cacheDir string) error {
	if err := a.configManager.SetInstallDir(installDir); err != nil {
		return err
	}
	if err := a.configManager.SetCacheDir(cacheDir); err != nil {
		return err
	}
	// 更新版本管理器的路径
	a.versionManager.SetPaths(installDir, cacheDir)
	// 更新缓存管理器的路径
	a.cacheManager.SetCacheDir(cacheDir)
	return nil
}

// SelectDirectory 打开目录选择对话框
func (a *App) SelectDirectory(title string) (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
}

// ==================== 配置管理 ====================

// GetConfig 获取配置
func (a *App) GetConfig() *config.Config {
	return a.configManager.Get()
}

// SetMirror 设置镜像源
func (a *App) SetMirror(mirror string) error {
	if err := a.configManager.SetMirror(mirror); err != nil {
		return err
	}
	a.versionManager.SetMirror(mirror)
	return nil
}

// SetProxy 设置代理
func (a *App) SetProxy(proxy string) error {
	return a.configManager.SetProxy(proxy)
}

// SetGoProxy 设置 GOPROXY
func (a *App) SetGoProxy(goproxy string) error {
	return a.configManager.SetGoProxy(goproxy)
}

// SetTheme 设置主题
func (a *App) SetTheme(theme string) error {
	if err := a.configManager.SetTheme(theme); err != nil {
		return err
	}
	// 通知前端主题变更
	wailsRuntime.EventsEmit(a.ctx, "theme-changed", theme)
	return nil
}

// ResetConfig 重置配置
func (a *App) ResetConfig() error {
	return a.configManager.Reset()
}

// GetMirrorOptions 获取镜像源选项
func (a *App) GetMirrorOptions() []config.MirrorOption {
	return config.MirrorOptions
}

// ==================== 环境管理 ====================

// GetEnvInfo 获取环境信息
func (a *App) GetEnvInfo() *env.EnvInfo {
	return a.envManager.GetEnvInfo()
}

// DiagnoseEnv 诊断环境
func (a *App) DiagnoseEnv() []env.DiagnoseResult {
	return a.envManager.Diagnose()
}

// SetEnvVar 设置环境变量
func (a *App) SetEnvVar(name, value string) error {
	return a.envManager.SetEnvVar(name, value)
}

// FixGoRoot 修复 GOROOT
func (a *App) FixGoRoot(goroot string) error {
	if err := a.envManager.SetEnvVar("GOROOT", goroot); err != nil {
		return err
	}
	return a.envManager.UpdatePATH(goroot)
}

// FixGoProxy 修复 GOPROXY
func (a *App) FixGoProxy() error {
	return a.envManager.SetEnvVar("GOPROXY", "https://goproxy.cn,direct")
}

// ==================== 缓存管理 ====================

// GetCacheInfo 获取缓存信息
func (a *App) GetCacheInfo() (*cache.CacheInfo, error) {
	return a.cacheManager.GetCacheInfo()
}

// CleanDownloadCache 清理下载缓存
func (a *App) CleanDownloadCache() error {
	return a.cacheManager.CleanDownloadCache()
}

// CleanAllCache 清理所有缓存
func (a *App) CleanAllCache() error {
	return a.cacheManager.CleanAll()
}

// ==================== 系统信息 ====================

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() map[string]string {
	return map[string]string{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"baseDir": a.baseDir,
	}
}

// GetAppVersion 获取应用版本
func (a *App) GetAppVersion() string {
	return "1.0.0"
}

// ==================== 包管理 ====================

// GetGoPathMode 获取 GOPATH 模式
func (a *App) GetGoPathMode() string {
	return a.configManager.Get().GoPathMode
}

// SetGoPathMode 设置 GOPATH 模式
func (a *App) SetGoPathMode(mode string) error {
	if err := a.configManager.SetGoPathMode(mode); err != nil {
		return err
	}
	cfg := a.configManager.Get()
	a.moduleManager.SetGoPathMode(mode, cfg.SharedGoPath)
	return nil
}

// GetSharedGoPath 获取共享 GOPATH
func (a *App) GetSharedGoPath() string {
	return a.configManager.Get().SharedGoPath
}

// SetSharedGoPath 设置共享 GOPATH
func (a *App) SetSharedGoPath(path string) error {
	if err := a.configManager.SetSharedGoPath(path); err != nil {
		return err
	}
	cfg := a.configManager.Get()
	a.moduleManager.SetGoPathMode(cfg.GoPathMode, path)
	return nil
}

// GetModuleCacheStats 获取模块缓存统计
func (a *App) GetModuleCacheStats() (*module.ModuleCacheStats, error) {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.GetModuleCacheStats(currentVersion)
}

// GetModules 获取模块列表
func (a *App) GetModules() ([]module.ModuleInfo, error) {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.GetModules(currentVersion)
}

// CleanModuleCache 清理模块缓存
func (a *App) CleanModuleCache() error {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.CleanModuleCache(currentVersion)
}

// CleanModule 清理指定模块
func (a *App) CleanModule(modulePath, moduleVersion string) error {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.CleanModule(modulePath, moduleVersion, currentVersion)
}

// VerifyModules 校验模块
func (a *App) VerifyModules() ([]module.VerifyResult, error) {
	return a.moduleManager.VerifyModules()
}

// GetCurrentGoPath 获取当前 GOPATH
func (a *App) GetCurrentGoPath() string {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.GetGoPath(currentVersion)
}

// GetModuleCachePath 获取模块缓存路径
func (a *App) GetModuleCachePath() string {
	currentVersion := a.versionManager.GetCurrentVersion()
	return a.moduleManager.GetModuleCachePath(currentVersion)
}

// SearchPackages 搜索包
func (a *App) SearchPackages(query string) ([]module.SearchResult, error) {
	return a.moduleManager.SearchPackages(query)
}

// InstallPackage 安装包
func (a *App) InstallPackage(packagePath string) error {
	return a.moduleManager.InstallPackage(packagePath)
}

// GetPackage 下载包到缓存
func (a *App) GetPackage(packagePath, version string) error {
	return a.moduleManager.GetPackage(packagePath, version)
}

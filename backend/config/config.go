package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config 应用配置
type Config struct {
	Port           string // HTTP 服务端口
	KubeconfigPath string // 默认集群 kubeconfig 路径
	JWTSecret      string // JWT 签名密钥
	DBPath         string // SQLite 数据库文件路径
	EncryptKey     string // 集群凭据加密密钥（任意长度，内部 SHA-256 派生）
	TLSSkipVerify  bool   // 是否跳过集群 TLS 证书校验（仅开发环境）
	GinMode        string // gin 运行模式: debug/release/test
}

// App 全局配置单例，供不便通过依赖注入获取配置的包使用
var App *Config

// LoadConfig 加载配置（优先环境变量，缺失项使用安全默认值并告警）
func LoadConfig() *Config {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		KubeconfigPath: getEnv("KUBECONFIG", defaultKubeconfigPath()),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		DBPath:         getEnv("DB_PATH", filepath.Join(dataDir(), "kubeadm.db")),
		EncryptKey:     getEnv("ENCRYPT_KEY", ""),
		TLSSkipVerify:  getEnv("TLS_SKIP_VERIFY", "false") == "true",
		GinMode:        getEnv("GIN_MODE", "debug"),
	}

	// 安全告警：生产关键配置缺失时给出明确提示
	if cfg.JWTSecret == "" {
		log.Println("[WARN] JWT_SECRET 未设置，使用开发默认值，生产环境必须配置")
		cfg.JWTSecret = "dev-only-jwt-secret-change-me"
	}
	if cfg.EncryptKey == "" {
		log.Println("[WARN] ENCRYPT_KEY 未设置，使用开发默认值，集群凭据将用弱密钥加密")
		cfg.EncryptKey = "dev-only-encrypt-key-change-me"
	}
	if cfg.TLSSkipVerify {
		log.Println("[WARN] TLS_SKIP_VERIFY=true，集群 TLS 证书校验已关闭，仅限开发环境")
	}

	App = cfg
	return cfg
}

// getEnv 读取环境变量，缺失时返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// dataDir 返回数据目录路径（优先项目内 ./data，保证可写）
func dataDir() string {
	dir := "data"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		// 退化到用户主目录，避免启动失败
		if home := homeDir(); home != "" {
			dir = filepath.Join(home, ".kube-admin", "data")
			_ = os.MkdirAll(dir, 0o755)
		}
	}
	return dir
}

// defaultKubeconfigPath 返回默认 kubeconfig 路径
func defaultKubeconfigPath() string {
	if home := homeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return ""
}

// homeDir 返回用户主目录（兼容 Unix/Windows）
func homeDir() string {
	if strings.TrimSpace(os.Getenv("HOME")) != "" {
		return os.Getenv("HOME")
	}
	return os.Getenv("USERPROFILE") // windows
}

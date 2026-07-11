package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config 应用配置
type Config struct {
	Port           string // HTTP 服务端口
	KubeconfigPath string // 默认集群 kubeconfig 路径
	JWTSecret      string // JWT 签名密钥
	DBPath         string // SQLite 数据库文件路径（仅 DB_DRIVER=sqlite 生效）
	DBDriver       string // 数据库驱动：sqlite | mysql | postgres
	DBDSN          string // mysql/postgres 连接串（DB_DRIVER 非 sqlite 时必填）
	EncryptKey      string        // 集群凭据加密密钥（任意长度，内部 SHA-256 派生）
	TLSSkipVerify   bool          // 是否跳过集群 TLS 证书校验（仅开发环境）
	K8sTimeout      time.Duration // k8s API 单次请求超时（K8S_REQUEST_TIMEOUT 秒，默认 10s，避免集群不可达时挂 30s）
	GinMode         string        // gin 运行模式: debug/release/test
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
		DBDriver:       strings.ToLower(getEnv("DB_DRIVER", "sqlite")),
		DBDSN:          getEnv("DB_DSN", ""),
		EncryptKey:     getEnv("ENCRYPT_KEY", ""),
		TLSSkipVerify:  getEnv("TLS_SKIP_VERIFY", "false") == "true",
		K8sTimeout:     k8sTimeoutFromEnv(),
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

	// 校验数据库驱动
	switch cfg.DBDriver {
	case "sqlite", "mysql", "postgres":
	default:
		log.Fatalf("[FATAL] 不支持的 DB_DRIVER=%s，仅支持 sqlite/mysql/postgres", cfg.DBDriver)
	}
	if cfg.DBDriver != "sqlite" && cfg.DBDSN == "" {
		log.Printf("[WARN] DB_DRIVER=%s 但 DB_DSN 未设置，数据库连接将失败", cfg.DBDriver)
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

// k8sTimeoutFromEnv 解析 K8S_REQUEST_TIMEOUT（秒），非法或未设置回退 10s。
// 集群不可达时让请求快速失败，而非 client-go 默认挂起 30s。
func k8sTimeoutFromEnv() time.Duration {
	if v := os.Getenv("K8S_REQUEST_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
		log.Printf("[WARN] K8S_REQUEST_TIMEOUT=%s 非法，回退默认 10s", v)
	}
	return 10 * time.Second
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

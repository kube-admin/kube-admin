package config

import (
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	Port           string
	KubeconfigPath string
	JWTSecret      string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	var defaultKubeConfigDir string
	if home := homeDir(); home != "" {
		defaultKubeConfigDir = filepath.Join(home, ".kube", "config")
	} else {
		defaultKubeConfigDir = ""
	}
	return &Config{
		Port:           getEnv("PORT", "8080"),
		KubeconfigPath: getEnv("KUBECONFIG", defaultKubeConfigDir),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

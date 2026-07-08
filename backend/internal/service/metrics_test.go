package service

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

// TestCalcPercent 正常使用率计算
func TestCalcPercent(t *testing.T) {
	used := resource.MustParse("500m")
	total := resource.MustParse("1000m")
	if got := calcPercent(used, total); got != 50 {
		t.Fatalf("got %v want 50", got)
	}
}

// TestCalcPercentZeroTotal 总量为 0 应返回 0（避免除零）
func TestCalcPercentZeroTotal(t *testing.T) {
	used := resource.MustParse("500m")
	total := resource.MustParse("0")
	if got := calcPercent(used, total); got != 0 {
		t.Fatalf("zero total expect 0, got %v", got)
	}
}

// TestCalcPercentMemory 内存字节量级计算应正确
func TestCalcPercentMemory(t *testing.T) {
	used := resource.MustParse("512Mi")
	total := resource.MustParse("1Gi")
	got := calcPercent(used, total)
	if got <= 49 || got > 51 {
		t.Fatalf("512Mi/1Gi expect ~50, got %v", got)
	}
}

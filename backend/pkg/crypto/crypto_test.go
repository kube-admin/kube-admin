package crypto

import "testing"

// TestInitEmptyKey 空密钥应拒绝初始化
func TestInitEmptyKey(t *testing.T) {
	if err := Init(""); err == nil {
		t.Fatal("空密钥应返回错误")
	}
}

// TestEncryptDecryptRoundTrip 加解密往返应一致
func TestEncryptDecryptRoundTrip(t *testing.T) {
	if err := Init("unit-test-key"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	cases := []string{"hello", "kubeconfig-content", "token-abc-123", "中文凭据", ""}
	for _, plain := range cases {
		enc, err := Encrypt(plain)
		if err != nil {
			t.Fatalf("Encrypt(%q) error: %v", plain, err)
		}
		// 空串应原样返回，不加密
		if plain == "" && enc != "" {
			t.Fatalf("空串应返回空，实际 %q", enc)
		}
		if plain == "" {
			continue
		}
		dec, err := Decrypt(enc)
		if err != nil {
			t.Fatalf("Decrypt error: %v", err)
		}
		if dec != plain {
			t.Fatalf("往返不一致: got %q want %q", dec, plain)
		}
	}
}

// TestEncryptNotEqualPlain 密文不应与明文相同
func TestEncryptNotEqualPlain(t *testing.T) {
	Init("another-key")
	enc, err := Encrypt("very-secret-token")
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if enc == "very-secret-token" {
		t.Fatal("密文不应等于明文")
	}
}

// TestDecryptInvalid 密文被篡改应解密失败
func TestDecryptInvalid(t *testing.T) {
	Init("some-key")
	if _, err := Decrypt("not-a-valid-base64!!!"); err == nil {
		t.Fatal("非法密文应返回错误")
	}
}

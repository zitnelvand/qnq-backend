package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHTTPConfigDefaults(t *testing.T) {
	t.Setenv("HTTP_HOST", "")
	t.Setenv("HTTP_PORT", "")

	cfg, err := LoadHTTPConfig()
	if err != nil {
		t.Fatalf("LoadHTTPConfig() error = %v", err)
	}
	if cfg.Host != "" {
		t.Errorf("Host = %q, want empty host", cfg.Host)
	}
	if cfg.Port != defaultHTTPPort {
		t.Errorf("Port = %q, want %q", cfg.Port, defaultHTTPPort)
	}
}

func TestLoadHTTPConfigUsesEnvironment(t *testing.T) {
	t.Setenv("HTTP_HOST", "127.0.0.1")
	t.Setenv("HTTP_PORT", "8080")

	cfg, err := LoadHTTPConfig()
	if err != nil {
		t.Fatalf("LoadHTTPConfig() error = %v", err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != "8080" {
		t.Fatalf("LoadHTTPConfig() = %+v, want host 127.0.0.1 and port 8080", cfg)
	}
}

func TestHTTPConfigAddressSupportsIPv6(t *testing.T) {
	cfg := HTTPConfig{Host: "2001:db8::1", Port: "8080"}

	if got, want := cfg.Address(), "[2001:db8::1]:8080"; got != want {
		t.Errorf("Address() = %q, want %q", got, want)
	}
}

func TestLoadHTTPConfigRejectsInvalidPort(t *testing.T) {
	for _, port := range []string{"abc", "0", "65536"} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("HTTP_HOST", "")
			t.Setenv("HTTP_PORT", port)

			if _, err := LoadHTTPConfig(); err == nil {
				t.Fatal("LoadHTTPConfig() error = nil, want invalid port error")
			}
		})
	}
}

func TestLoadDotEnv(t *testing.T) {
	const key = "QNQ_DOTENV_TEST_VALUE"
	unsetEnv(t, key)

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(key+"=from-dotenv\n"), 0o600); err != nil {
		t.Fatalf("write dotenv file: %v", err)
	}

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv() error = %v", err)
	}
	if got := os.Getenv(key); got != "from-dotenv" {
		t.Errorf("%s = %q, want %q", key, got, "from-dotenv")
	}
}

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	const key = "QNQ_DOTENV_TEST_VALUE"
	t.Setenv(key, "from-environment")

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(key+"=from-dotenv\n"), 0o600); err != nil {
		t.Fatalf("write dotenv file: %v", err)
	}

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv() error = %v", err)
	}
	if got := os.Getenv(key); got != "from-environment" {
		t.Errorf("%s = %q, want %q", key, got, "from-environment")
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	old, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if existed {
			if err := os.Setenv(key, old); err != nil {
				t.Errorf("restore %s: %v", key, err)
			}
			return
		}
		if err := os.Unsetenv(key); err != nil {
			t.Errorf("unset %s: %v", key, err)
		}
	})
}

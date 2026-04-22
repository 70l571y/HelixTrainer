package cfg

import "testing"

func TestConfigDirFromBasePath(t *testing.T) {
	got := configDirFromBase("/tmp/config-root")
	want := "/tmp/config-root/hxtrainer"
	if got != want {
		t.Fatalf("configDirFromBase() = %q, want %q", got, want)
	}
}

func TestFallbackConfigBaseUsesHomeDirectory(t *testing.T) {
	got := fallbackConfigBase("linux", "/home/tester")
	want := "/home/tester/.config"
	if got != want {
		t.Fatalf("fallbackConfigBase() = %q, want %q", got, want)
	}
}

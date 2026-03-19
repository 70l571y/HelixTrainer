package buildinfo

import (
	"runtime/debug"
	"strings"
)

var Version string

func CurrentVersion() string {
	if version := normalizeVersion(Version); version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if ok {
		if version := normalizeVersion(info.Main.Version); version != "" {
			return version
		}
	}

	return "devel"
}

func normalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	switch version {
	case "", "(devel)", "devel":
		return ""
	default:
		version = strings.TrimPrefix(version, "v")
		if idx := strings.Index(version, "+"); idx >= 0 {
			version = version[:idx]
		}
		return version
	}
}

package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current AGY++ version.
	Version = "0.1.0-dev"
	// GitCommit is the commit hash populated at build time.
	GitCommit = "none"
	// BuildDate is the build date populated at build time.
	BuildDate = "unknown"
	// BaselineAgy is the target upstream AGY compatibility baseline.
	BaselineAgy = "1.2.2"
)

// String returns formatted version details.
func String() string {
	return fmt.Sprintf("AGY++ v%s (%s) built %s [%s/%s]\nUpstream AGY Baseline: %s",
		Version, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH, BaselineAgy)
}

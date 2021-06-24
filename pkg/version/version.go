package version

import (
	"fmt"
)

// These variables are substituted with real values at build time
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = "dev"
)

// String returns the version as a formatted string
func String() string {
	return fmt.Sprintf("%s (%s@%s by %s)", version, commit, date, builtBy)
}

func Version() string {
	return version
}

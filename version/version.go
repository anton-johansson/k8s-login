package version

var (
	version   string
	goVersion string
	commit    string
)

// VersionInfo holds information about the current version
type VersionInfo struct {
	Version   string
	GoVersion string
	Commit    string
}

// Info returns the version information
func Info() VersionInfo {
	return VersionInfo{
		Version:   version,
		GoVersion: goVersion,
		Commit:    commit,
	}
}
